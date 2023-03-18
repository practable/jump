package host

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/crossbar"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/reconws"
	"github.com/practable/jump/internal/relay"
	"github.com/practable/jump/internal/tcpconnect"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {

	debug := false
	if debug {
		log.SetReportCaller(true)
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}
}

func TestHost(t *testing.T) {

	// Setup logging

	timeout := 100 * time.Millisecond

	// setup relay on local (free) port
	ports, err := freeport.GetFreePorts(2)
	assert.NoError(t, err)

	relayPort := ports[0]
	accessPort := ports[1]

	accessURI := "http://[::]:" + strconv.Itoa(accessPort)
	relayURI := "ws://127.0.0.1:" + strconv.Itoa(relayPort)

	base_path := "/api/v1"

	log.Debug(fmt.Sprintf("accessURI:%s\n", accessURI))
	log.Debug(fmt.Sprintf("relayURI:%s\n", relayURI))

	secret := "testsecret"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := relay.Config{
		AccessPort:     accessPort,
		Audience:       accessURI,
		ConnectionType: "connect",
		RelayPort:      relayPort,
		Secret:         secret,
		Target:         relayURI,
	}
	go relay.Run(ctx, config)

	// setup mock sshd

	defer cancel()

	sshPort, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	sshduri := ":" + strconv.Itoa(sshPort)

	echo := tcpconnect.New()
	go echo.Listen(ctx, sshduri, tcpconnect.SpeakThenEchoHandler)

	time.Sleep(2 * time.Second)

	// setup host
	ct := "connect"
	session := "11014d77-e36e-40b7-9864-5a9239d1a071"
	scopes := []string{"host"} //host, client scopes are known only to access

	begin := time.Now().Unix() - 1 //ensure it's in the past
	end := begin + 180
	claims := permission.NewToken(accessURI, ct, session, scopes, begin, begin, end)
	hostToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hostBearer, err := hostToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	go Run(ctx, "localhost"+sshduri, accessURI+base_path+"/connect/"+session, hostBearer)

	time.Sleep(time.Second)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	// ============================= START  TESTS ======================================

	// *** TestConnectToLocal ***

	scopes = []string{"client"} //host, client scopes are known only to access
	claims = permission.NewToken(accessURI, "connect", session, scopes, begin, begin, end)
	clientToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	clientBearer, err := clientToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	clientURI := accessURI + base_path + "/connect/" + session

	c0 := reconws.New()
	go c0.ReconnectAuth(ctx, clientURI, clientBearer)

	c1 := reconws.New()
	go c1.ReconnectAuth(ctx, clientURI, clientBearer)

	// Send messages, get echos...
	time.Sleep(3 * time.Second) //give host a chance to make new connections

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	data0 := []byte("ping")

	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case c0.Out <- reconws.WsMessage{Data: data0, Type: websocket.BinaryMessage}:

	}

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	// get greetings

	greeting := []byte("Echo Service")

	select {
	case <-time.After(timeout):
		t.Error("timeout on greeting")
	case msg, ok := <-c0.In:
		assert.True(t, ok)
		assert.Equal(t, greeting, msg.Data)
	}

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	select {
	case <-time.After(timeout):
		t.Error("timeout on greeting")
	case msg, ok := <-c1.In:
		assert.True(t, ok)
		assert.Equal(t, greeting, msg.Data)
	}

	// get echo
	select {
	case <-time.After(timeout):
		t.Error("timeout")
	case msg, ok := <-c0.In:
		assert.True(t, ok)
		assert.Equal(t, data0, msg.Data)
		t.Log("TestConnectToLocal...PASS")
	}

	// get nothing as planned
	select {
	case <-time.After(timeout):
	case <-c1.In:
		t.Fatal("unexpected")
	}

	// send on other client c1

	data1 := []byte("foo")
	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case c1.Out <- reconws.WsMessage{Data: data1, Type: websocket.BinaryMessage}:
	}

	// get echo
	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case msg, ok := <-c1.In:
		assert.True(t, ok)
		assert.Equal(t, data1, msg.Data)
	}

	// get nothing on other client c0 as expected
	select {
	case <-time.After(timeout):
	case <-c0.In:
		t.Fatal("unexpected")
	}

	// while connected, get stats
	scopes = []string{"stats"}
	claims = permission.NewToken(accessURI, "connect", "stats", scopes, begin, begin, end)
	statsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	statsBearer, err := statsToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	stats := reconws.New()
	go stats.ReconnectAuth(ctx, accessURI+base_path+"/connect/stats", statsBearer)

	cmd, err := json.Marshal(crossbar.StatsCommand{Command: "update"})

	assert.NoError(t, err)

	stats.Out <- reconws.WsMessage{Data: cmd, Type: websocket.TextMessage}

	select {
	case msg := <-stats.In:

		var reports []*crossbar.ClientReport

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

		sessionCount := 0
		for topic, count := range agents {
			log.Debug(topic)
			if strings.HasPrefix(topic, session+"/") {
				sessionCount = sessionCount + count
			}
		}
		expectedCount := 4
		assert.Equal(t, expectedCount, sessionCount)

		//TODO we can't know this because salted, so search for partial match to session
		if sessionCount == expectedCount {
			t.Log("TestGetStats...PASS")
		} else {
			pretty, err := json.MarshalIndent(reports, "", "\t")
			assert.NoError(t, err)
			t.Log(string(pretty))
			t.Fatalf("TestGetStats...FAIL (wrong agent count)")
		}

	case <-time.After(timeout):
		t.Fatalf("TestGetStats...FAIL (timeout)")
	}

	// ================================== Teardown  ===============================================

}
