package client

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
	"github.com/practable/jump/internal/host"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/reconws"
	"github.com/practable/jump/internal/relay"
	"github.com/practable/jump/internal/tcpconnect"
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
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}
}

func TestClient(t *testing.T) {

	// Setup logging

	timeout := 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup relay on local (free) port
	ports, err := freeport.GetFreePorts(3)
	assert.NoError(t, err)

	relayPort := ports[0]
	accessPort := ports[1]
	clientPort := ports[2]

	audience := "http://[::]:" + strconv.Itoa(accessPort)
	target := "ws://127.0.0.1:" + strconv.Itoa(relayPort)
	clientURI := "http://[::]:" + strconv.Itoa(clientPort)

	log.Debug(fmt.Sprintf("audience:%s\n", audience))
	log.Debug(fmt.Sprintf("target:%s\n", target))
	log.Debug(fmt.Sprintf("clientURI:%s\n", clientURI))

	secret := "testsecret"

	config := relay.Config{
		AccessPort:     accessPort,
		Audience:       audience,
		ConnectionType: "connect",
		RelayPort:      relayPort,
		Secret:         secret,
		Target:         target,
		StatsEvery:     time.Duration(time.Second),
	}

	go relay.Run(ctx, config)

	// setup mock sshd

	defer cancel()

	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	sshduri := ":" + strconv.Itoa(port)

	echo := tcpconnect.New()
	go echo.Listen(ctx, sshduri, tcpconnect.SpeakThenEchoHandler)

	time.Sleep(2 * time.Second)

	// setup host
	basePath := "/api/v1"
	ct := "connect"
	session := "11014d77-e36e-40b7-9864-5a9239d1a071"
	scopes := []string{"host"} //host, client scopes are known only to access

	begin := time.Now().Unix() - 1 //ensure it's in the past
	end := begin + 180
	claims := permission.NewToken(audience, ct, session, scopes, begin, begin, end)
	hostToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hostBearer, err := hostToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	go host.Run(ctx, "localhost"+sshduri, audience+basePath+"/connect/"+session, hostBearer)

	time.Sleep(time.Second)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	// *** Setup client ***

	scopes = []string{"client"} //host, client scopes are known only to access
	claims = permission.NewToken(audience, "connect", session, scopes, begin, begin, end)
	clientToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	clientBearer, err := clientToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	go Run(ctx, clientPort, audience+basePath+"/connect/"+session, clientBearer)

	time.Sleep(2 * time.Second)

	// ============================= START  TESTS ======================================

	clientURI = "localhost:" + strconv.Itoa(clientPort)

	c0 := tcpconnect.New()
	go c0.Dial(ctx, clientURI)

	c1 := tcpconnect.New()
	go c1.Dial(ctx, clientURI)

	// Send messages, get echos...
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)
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
		assert.Equal(t, greeting, msg)
	}

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	select {
	case <-time.After(timeout):
		t.Error("timeout on greeting")
	case msg, ok := <-c1.In:
		assert.True(t, ok)
		assert.Equal(t, greeting, msg)
	}

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	data0 := []byte("ping")

	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case c0.Out <- data0:

	}

	time.Sleep(timeout)
	time.Sleep(timeout)
	time.Sleep(timeout)

	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case data, ok := <-c0.In:
		assert.True(t, ok)
		assert.Equal(t, data0, data)
		t.Log("TestConnectToLocal...PASS")
	}

	select {
	case <-time.After(timeout):
	case <-c1.In:
		t.Fatal("unexpected")
	}

	data1 := []byte("foo")
	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case c1.Out <- data1:
	}

	select {
	case <-time.After(timeout):
		t.Fatal("timeout")
	case data, ok := <-c1.In:
		assert.True(t, ok)
		assert.Equal(t, data1, data)
	}

	select {
	case <-time.After(timeout):
	case <-c0.In:
		t.Fatal("unexpected")
	}

	// while connected, get stats
	scopes = []string{"stats"}
	claims = permission.NewToken(audience, "connect", "stats", scopes, begin, begin, end)
	statsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	statsBearer, err := statsToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	stats := reconws.New()
	// in production use pkg/status
	go stats.ReconnectAuth(ctx, audience+basePath+"/connect/stats", statsBearer)

	cmd, err := json.Marshal(crossbar.StatsCommand{Command: "update"})

	assert.NoError(t, err)

	stats.Out <- reconws.WsMessage{Data: cmd, Type: websocket.TextMessage}

	select {
	case msg := <-stats.In:
		if debug {
			t.Log("Stats message received")
		}
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
	cancel()
}
