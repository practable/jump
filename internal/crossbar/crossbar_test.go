package crossbar

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/reconws"
	"github.com/practable/jump/internal/ttlcode"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var debug bool

func init() {
	debug = false
	if debug {
		log.SetReportCaller(true)
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: false, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		log.SetLevel(log.WarnLevel)
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

}

func MakeTestToken(audience, connectionType, topic string, scopes []string, lifetime int64) permission.Token {
	begin := time.Now().Unix() - 1 //ensure it's in the past
	end := begin + lifetime
	return permission.NewToken(audience, connectionType, topic, scopes, begin, begin, end)
}

func TestRun(t *testing.T) {

	// Renew the mux to avoid multiple registrations error
	http.DefaultServeMux = new(http.ServeMux)

	// setup shellbar on local (free) port

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	audience := "ws://127.0.0.1:" + strconv.Itoa(port)
	secret := "somesecret"
	cs := ttlcode.NewDefaultCodeStore()
	config := Config{
		Listen:         port,
		Audience:       audience,
		CodeStore:      cs,
		ConnectionType: "shell",
		Secret:         secret,
		StatsEvery:     time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go Run(ctx, config)
	// safety margin to get shellbar running
	time.Sleep(time.Second)

	var timeout = 100 * time.Millisecond

	// Start tests

	// *** TestConnectUniquely ***

	// construct host token & connect
	ct := "shell"
	session := "abc"
	scopes := []string{"read", "write"} //host, client scopes are known only to access

	tokenHost := MakeTestToken(audience, ct, session, scopes, 30)
	codeHost := cs.SubmitToken(tokenHost)

	h := reconws.New()
	go func() {
		err := h.Dial(ctx, audience+"/"+ct+"/"+session+"?code="+codeHost)
		assert.NoError(t, err)
	}()

	// ensure we connect first by pausing until a dummy message sends
	//  not needed in production - shellbar would be alive long before a client connects

	h.Out <- reconws.WsMessage{Type: websocket.BinaryMessage}

	// construct client token & connect
	connectionID := "def"
	clientTopic := session + "/" + connectionID
	topicSalt := "ghi"
	topicInHub := clientTopic + topicSalt
	tokenClient := MakeTestToken(audience, ct, clientTopic, scopes, 30)
	permission.SetTopicSalt(&tokenClient, topicSalt)
	permission.SetAlertHost(&tokenClient, true)

	codeClient0 := cs.SubmitToken(tokenClient)
	c0 := reconws.New()
	client0UniqueURI := audience + "/" + ct + "/" + clientTopic

	ctx0, cancel0 := context.WithCancel(context.Background())

	go func() {
		err := c0.Dial(ctx0, client0UniqueURI+"?code="+codeClient0)
		assert.NoError(t, err)
	}()

	var ca ConnectionAction

	var c0UUID string

	select {

	case <-time.After(time.Second):
		t.Error("TestHostAdminGetsConnectAction...FAIL\n")

	case msg, ok := <-h.In:

		assert.True(t, ok)

		err = json.Unmarshal(msg.Data, &ca)
		assert.NoError(t, err)
		assert.Equal(t, "connect", ca.Action)

		base := strings.Split(ca.URI, "?")[0]
		c0UUID = ca.UUID
		assert.Equal(t, client0UniqueURI, base)
		if client0UniqueURI == base {
			t.Logf("TestHostAdminGetsConnectAction...PASS\n")
		} else {
			t.Fatal("TestHostAdminGetsConnectAction...FAIL\n")
		}
	}

	// Host now dials the unqiue connection

	h1 := reconws.New()
	go func() {
		err := h1.Dial(ctx, ca.URI)
		assert.NoError(t, err)
	}()

	time.Sleep(timeout)

	data := []byte("ping")

	h1.Out <- reconws.WsMessage{Data: data, Type: websocket.TextMessage}

	select {
	case msg := <-c0.In:
		assert.Equal(t, data, msg.Data)
		if reflect.DeepEqual(data, msg.Data) {
			t.Logf("TestHostConnectsToUniqueSession...PASS\n")
		} else {
			t.Fatal("TestHostConnectsToUniqueSession...FAIL")
		}
	case <-time.After(timeout):
		t.Fatal("TestHostConnectsToUniqueSession...FAIL")
	}

	data = []byte("pong")

	c0.Out <- reconws.WsMessage{Data: data, Type: websocket.TextMessage}
	select {
	case msg := <-h1.In:
		assert.Equal(t, data, msg.Data)
		if reflect.DeepEqual(data, msg.Data) {
			t.Logf("TestHostReceivesDataFromUniqueSession...PASS\n")
		} else {
			t.Fatal("TestHostReceivesDataFromUniqueSession...FAIL (wrong message)")
		}
	case <-time.After(timeout):
		t.Fatal("TestHostReceivesDataFromUniqueSession...FAIL")
	}

	// while connected, get stats
	scopes = []string{"read", "write"}
	statsToken := MakeTestToken(audience, ct, "stats", scopes, 30)
	statsCode := cs.SubmitToken(statsToken)
	stats := reconws.New()

	go func() {
		err := stats.Dial(ctx, audience+"/"+ct+"/stats?code="+statsCode)
		assert.NoError(t, err)
	}()

	cmd, err := json.Marshal(StatsCommand{Command: "update"})

	assert.NoError(t, err)

	stats.Out <- reconws.WsMessage{Data: cmd, Type: websocket.TextMessage}

	select {
	case msg := <-stats.In:

		t.Log("TestGetStats...PROVISIONAL-PASS")

		var reports []*ClientReport

		err := json.Unmarshal(msg.Data, &reports)

		assert.NoError(t, err)

		agents := make(map[string]int)

		for _, report := range reports {
			count, ok := agents[report.Topic]
			if !ok {
				agents[report.Topic] = 1
				continue
			}

			agents[report.Topic] = count + 1
		}

		if agents[topicInHub] == 2 {
			t.Log("TestGetStats...PASS")
		} else {
			t.Fatalf("TestGetStats...FAIL")
			pretty, err := json.MarshalIndent(reports, "", "\t")
			assert.NoError(t, err)
			fmt.Println(string(pretty))
		}

	case <-time.After(timeout):
		t.Fatalf("TestGetStats...FAIL")
	}

	time.Sleep(timeout)

	cancel0()

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	select {
	case <-time.After(time.Second):
		t.Fatal("No disconnect message")
	case msg, ok := <-h.In:
		assert.True(t, ok)

		err = json.Unmarshal(msg.Data, &ca)
		assert.NoError(t, err)
		assert.Equal(t, "disconnect", ca.Action)
		assert.Equal(t, c0UUID, ca.UUID)

		if c0UUID == ca.UUID {
			t.Logf("TestHostAdminGetsDisconnectAction...PASS\n")
		} else {
			t.Fatal("TestHostAdminGetsDisconnectAction...FAIL\n")
		}

	}

	time.Sleep(timeout)

	// construct host token & connect
	ct = "shell"
	session = "rst"
	scopes = []string{"read", "write"} //host, client scopes are known only to access

	tokenHost = MakeTestToken(audience, ct, session, scopes, 30)
	codeHost = cs.SubmitToken(tokenHost)

	hh := reconws.New()
	go func() {
		err := hh.Dial(ctx, audience+"/"+ct+"/"+session+"?code="+codeHost)
		assert.NoError(t, err)
	}()

	// ensure host connects first by pausing until a dummy message sends
	//  not needed in production - shellbar would be alive long before a client connects

	hh.Out <- reconws.WsMessage{Type: websocket.BinaryMessage}

	// construct client token & connect
	connectionID = "uvw"
	clientTopic = session + "/" + connectionID
	topicSalt = "xyz"
	tokenClient = MakeTestToken(audience, ct, clientTopic, scopes, 30)
	permission.SetTopicSalt(&tokenClient, topicSalt)
	permission.SetAlertHost(&tokenClient, true)

	codeClient2 := cs.SubmitToken(tokenClient)
	c2 := reconws.New()
	client2UniqueURI := audience + "/" + ct + "/" + clientTopic
	c2uri := client2UniqueURI + "?code=" + codeClient2
	go func() {
		err := c2.Dial(ctx, c2uri)
		assert.NoError(t, err)
	}()

	// construct second client token & connect
	connectionID = "Bf6380c7-c444-4e99-aec7-11272a690bc5"
	clientTopic = session + "/" + connectionID
	topicSalt = "B9638f36-9c20-4d8d-84e9-65d4e0410126"
	tokenClient = MakeTestToken(audience, ct, clientTopic, scopes, 30)
	permission.SetTopicSalt(&tokenClient, topicSalt)
	permission.SetAlertHost(&tokenClient, true)

	codeClient3 := cs.SubmitToken(tokenClient)
	c3 := reconws.New()
	client3UniqueURI := audience + "/" + ct + "/" + clientTopic

	c3uri := client3UniqueURI + "?code=" + codeClient3
	go func() {
		err := c3.Dial(ctx, c3uri)
		assert.NoError(t, err)
	}()

	log.Debug(c2uri)
	log.Debug(c3uri)

	// make a list of connectionActions we receive, so that we don't have to
	// rely on them coming in order - a sleep between dials does not
	// guarantee order.

	time.Sleep(1 * time.Second)

	var cas []ConnectionAction

	timeout = 50 * time.Millisecond

	for n := 0; n < 100; n++ {

		// test intermittently fails depending on the timing
		// employed in this loop
		// this is considered a test artefact
		// since many goros running in this thread
		// does not fail in >10 attempts with -race

		select {

		case <-time.After(timeout):

		case msg, ok := <-hh.In:

			assert.True(t, ok)

			err = json.Unmarshal(msg.Data, &ca)
			assert.NoError(t, err)

			cas = append(cas, ca)
		}

		if len(cas) >= 2 {
			log.Debugf("Got connections after %d loops", n)
			break
		}

	}

	var cac2, cac3 int

	for _, ca := range cas {

		assert.Equal(t, "connect", ca.Action)

		base := strings.Split(ca.URI, "?")[0]

		if client2UniqueURI == base {
			cac2 = cac2 + 1
		}
		if client3UniqueURI == base {
			cac3 = cac3 + 1
		}

	}

	assert.Equal(t, 1, cac2)
	assert.Equal(t, 1, cac3)

	if cac2 == 1 && cac3 == 1 {
		t.Logf("TestHostAdminGetsMultipleConnectActions...PASS\n")
	} else {
		t.Errorf("TestHostAdminGetsMultipleConnectActions...FAIL\n")
		fmt.Println(pretty(cas))
	}

	// let tests finish before concelling the clients
	time.Sleep(timeout)
	cancel()
}

func pretty(t interface{}) string {

	json, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		return ""
	}

	return string(json)
}

func TestSlashify(t *testing.T) {

	if "/foo" != slashify("foo") {
		t.Errorf("Slashify not prefixing slash ")
	}
	if "//foo" == slashify("/foo") {
		t.Errorf("Slashify prefixing additional slash")
	}
	if "/foo" != slashify("/foo/") {
		t.Errorf("Slashify not removing trailing slash")
	}
	if "/foo" != slashify("foo/") {
		t.Errorf("Slashify not both removing trailing slash AND prefixing slash")
	}

	b := "foo/bar/rab/oof/"
	if "/foo/bar/rab/oof" != slashify(b) {
		t.Errorf("Slashify not coping with internal slashes %s -> %s", b, slashify(b))
	}

}

func TestGetConnectionTypeFromPath(t *testing.T) {

	assert.Equal(t, "connectionType", getConnectionTypeFromPath("/connectionType/shellID"))
	assert.Equal(t, "", getConnectionTypeFromPath("NoLeadingSlash/A/B/C"))
	assert.Equal(t, "foo%20bar", getConnectionTypeFromPath("/foo%20bar/glum"))
	assert.Equal(t, "", getConnectionTypeFromPath("ooops/foo%20bar/glum"))
}

func TestGetHostTopicFromUniqueTopic(t *testing.T) {

	assert.Equal(t, "shellID", getHostTopicFromUniqueTopic("shellID"))
	assert.Equal(t, "NoLeadingSlash", getHostTopicFromUniqueTopic("NoLeadingSlash/A/B/C"))
	assert.Equal(t, "", getHostTopicFromUniqueTopic("/foo%20bar/glum"))
	assert.Equal(t, "ooops", getHostTopicFromUniqueTopic("ooops/foo%20bar/glum"))
}

func TestGetTopicFromPath(t *testing.T) {

	assert.Equal(t, "shellID", getTopicFromPath("/connectionType/shellID"))
	assert.Equal(t, "", getTopicFromPath("NoLeadingSlash/A/B/C"))
	assert.Equal(t, "shell%20ID/connection%20ID", getTopicFromPath("/connectionType/shell%20ID/connection%20ID"))
	assert.Equal(t, "shellID/connectionID", getTopicFromPath("/connectionType/shellID/connectionID?QueryParams=Something"))
	assert.Equal(t, "shellID/connectionID", getTopicFromPath("/connectionType/shellID/connectionID?QueryParams=Something&SomeThing=Else"))
}
