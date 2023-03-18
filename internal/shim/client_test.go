package shim

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/relay"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {

	log.SetLevel(log.WarnLevel)

}

func makeTestToken(audience, secret string, ttl int64) (string, error) {

	var claims permission.Token

	start := jwt.NewNumericDate(time.Now().Add(-time.Second))
	afterTTL := jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Second))
	claims.IssuedAt = start
	claims.NotBefore = start
	claims.ExpiresAt = afterTTL
	claims.Audience = jwt.ClaimStrings{audience}
	claims.Topic = "stats"
	claims.ConnectionType = "connect"
	claims.Scopes = []string{"stats"}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}

func TestClientConnect(t *testing.T) {

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Setup relay on local (free) port
	ports, err := freeport.GetFreePorts(2)
	assert.NoError(t, err)

	relayPort := ports[0]
	accessPort := ports[1]

	audience := "http://[::]:" + strconv.Itoa(accessPort)
	target := "ws://127.0.0.1:" + strconv.Itoa(relayPort)
	base_path := "/api/v1"

	secret := "testsecret"

	go func() {
		time.Sleep(2 * time.Second)
		config := relay.Config{
			AccessPort:     accessPort,
			RelayPort:      relayPort,
			Audience:       audience,
			ConnectionType: "connect",
			Secret:         secret,
			Target:         target,
			StatsEvery:     time.Duration(time.Second),
		}
		go relay.Run(ctx, config)
	}()

	// we sleep before starting the relay to help avoid issues with multiple
	// handlers registering with net/http when running all tests

	// Sign and get the complete encoded token as a string using the secret
	token, err := makeTestToken(audience, secret, 30)

	assert.NoError(t, err)

	// now clients connect using their uris...
	to := audience + base_path + "/connect/stats"

	// wait until relay has been up for about one second
	time.Sleep(3 * time.Second)

	c0 := NewClient()
	go c0.Connect(ctx, to, token)

	select {
	case <-time.After(5 * time.Second):
		t.Error("timed out waiting for message")
	case <-c0.Receive:
		//should get a stats message every second
	}

}
