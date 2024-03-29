package relay

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/crossbar"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/reconws"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRelay(t *testing.T) {

	// Setup logging
	debug := false

	if debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	// Setup relay on local (free) port
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ports, err := freeport.GetFreePorts(2)
	assert.NoError(t, err)

	relayPort := ports[0]
	accessPort := ports[1]

	audience := "http://[::]:" + strconv.Itoa(accessPort)
	target := "ws://127.0.0.1:" + strconv.Itoa(relayPort)
	base_path := "/api/v1"

	fmt.Printf("audience:%s\n", audience)
	fmt.Printf("target:%s\n", target)

	secret := "testsecret"

	config := Config{
		AccessPort:     accessPort,
		Audience:       audience,
		ConnectionType: "connect",
		RelayPort:      relayPort,
		Secret:         secret,
		Target:         target,
	}

	go Run(ctx, config)

	time.Sleep(time.Second) // big safety margin to get crossbar running

	// Start tests

	// TestBidirectionalChat

	var claims permission.Token

	start := jwt.NewNumericDate(time.Now().Add(-time.Second))
	after := jwt.NewNumericDate(time.Now().Add(30 * time.Second))
	claims.IssuedAt = start
	claims.NotBefore = start
	claims.ExpiresAt = after

	claims.Audience = jwt.ClaimStrings{audience}
	claims.Topic = "123"
	claims.ConnectionType = "connect"
	claims.Scopes = []string{"host"}

	hostToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	hostBearer, err := hostToken.SignedString([]byte(secret))
	assert.NoError(t, err)
	hostURI := audience + base_path + "/connect/123"

	h := reconws.New()
	go h.ReconnectAuth(ctx, hostURI, hostBearer)

	//hold until connected
	h.Out <- reconws.WsMessage{Type: websocket.TextMessage}

	// now connect a client
	claims.Scopes = []string{"client"}
	clientToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	clientBearer, err := clientToken.SignedString([]byte(secret))
	assert.NoError(t, err)

	c0 := reconws.New()
	go c0.ReconnectAuth(ctx, hostURI, clientBearer)

	// wait for client connection message

	var ca crossbar.ConnectionAction
	select {
	case msg, ok := <-h.In:
		assert.True(t, ok)
		err = json.Unmarshal(msg.Data, &ca)
		assert.NoError(t, err)
		assert.Equal(t, "connect", ca.Action)
	case <-time.After(time.Second):
		t.Fatal("Failed to get ConnectAction")
	}

	h1 := reconws.New()

	go func() {
		err := h1.Dial(ctx, ca.URI)
		assert.NoError(t, err)
	}()

	data := []byte("ping")
	h1.Out <- reconws.WsMessage{Data: data, Type: websocket.TextMessage}
	var timeout = 100 * time.Millisecond
	select {
	case msg, ok := <-c0.In:
		assert.True(t, ok)
		assert.Equal(t, data, msg.Data)
	case <-time.After(timeout):
		t.Fatal("Timed out getting ping")

	}

	c1 := reconws.New()
	go c1.ReconnectAuth(ctx, hostURI, clientBearer)

	select {
	case msg, ok := <-h.In:
		assert.True(t, ok)
		err = json.Unmarshal(msg.Data, &ca)
		assert.NoError(t, err)
		assert.Equal(t, "connect", ca.Action)
	case <-time.After(time.Second):
		t.Fatal("Failed to get ConnectAction")
	}

	h2 := reconws.New()

	go func() {
		err := h2.Dial(ctx, ca.URI)
		assert.NoError(t, err)
	}()

	time.Sleep(timeout)

	data = []byte("boo")

	h2.Out <- reconws.WsMessage{Data: data, Type: websocket.TextMessage}

	// c0 must not get this message
	select {
	case <-c0.In:
		t.Fatal("Got unexpected message")
	case <-time.After(timeout):
	}

	select {
	case msg, ok := <-c1.In:
		assert.True(t, ok)
		assert.Equal(t, data, msg.Data)
	case <-time.After(timeout):
		t.Fatal("Timed out getting boo")
	}

	// h admin
	// h1 services c0
	// h2 services c1
	// send message from c1, h1 must not get it
	data = []byte("far")

	c1.Out <- reconws.WsMessage{Data: data, Type: websocket.TextMessage}

	select {
	case <-h1.In:
		t.Fatal("Got unexpected message")
	case <-time.After(timeout):
	}

	select {
	case msg, ok := <-h2.In:
		assert.True(t, ok)
		assert.Equal(t, data, msg.Data)
	case <-time.After(timeout):
		t.Fatal("Timed out getting boo")
	}

}
