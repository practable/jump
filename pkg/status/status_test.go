package status

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

func TestUnmarshal(t *testing.T) {

	m := []byte(`{"topic":"stats","canRead":true,"canWrite":true,"connected":"2023-03-10T14:04:45.294633437Z","expiresAt":"2023-03-10T14:04:45.294633437Z","remoteAddr":"internal","userAgent":"crossbar","stats":{"tx":{"last":"Never","size":0,"fps":0},"rx":{"last":"Never","size":0,"fps":0}}}`)

	ma := []byte(`[{"topic":"stats","canRead":true,"canWrite":true,"connected":"2023-03-10T14:04:45.294633437Z","expiresAt":"2023-03-10T14:04:45.294633437Z","remoteAddr":"internal","userAgent":"crossbar","stats":{"tx":{"last":"Never","size":0,"fps":0},"rx":{"last":"Never","size":0,"fps":0}}},{"topic":"stats","canRead":true,"canWrite":false,"connected":"2023-03-10T14:04:48.102847403Z","expiresAt":"2023-03-10T14:04:48.102847403Z","remoteAddr":"","userAgent":"Go-http-client/1.1","stats":{"tx":{"last":"Never","size":0,"fps":0},"rx":{"last":"Never","size":0,"fps":0}}},{"topic":"123","canRead":true,"canWrite":true,"connected":"2023-03-10T14:04:46.294717207Z","expiresAt":"2023-03-10T14:04:46.294717207Z","remoteAddr":"","userAgent":"Go-http-client/1.1","stats":{"tx":{"last":"2.90373838s","size":5,"fps":20.173318460352565},"rx":{"last":"2.903832419s","size":5,"fps":20.192179268105992}}},{"topic":"123","canRead":true,"canWrite":true,"connected":"2023-03-10T14:04:46.29484872Z","expiresAt":"2023-03-10T14:04:46.29484872Z","remoteAddr":"","userAgent":"Go-http-client/1.1","stats":{"tx":{"last":"2.903868479s","size":5,"fps":20.225188461845494},"rx":{"last":"2.903726512s","size":4,"fps":10.097985601322723}}}]`)

	mb := []byte(`{"topic":"stats","canRead":true,"canWrite":true,"connected":"2023-03-10T14:04:45.294633437Z","expiresAt":"2023-03-10T14:04:45.294633437Z","remoteAddr":"internal","userAgent":"crossbar","stats":{"tx":{"last":3600000000000,"size":0,"fps":0},"rx":{"last":100000000,"size":0,"fps":0}}}`)

	assert.True(t, json.Valid(m))
	assert.True(t, json.Valid(ma))
	assert.True(t, json.Valid(mb))

	var report Report

	err := json.Unmarshal(m, &report)

	assert.NoError(t, err)

	var reports []Report

	err = json.Unmarshal(ma, &reports)

	assert.NoError(t, err)

	err = json.Unmarshal(mb, &report)

	assert.NoError(t, err)

	assert.Equal(t, time.Hour, report.Stats.Tx.Last)

}

func TestStatus(t *testing.T) {

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
	ports, err := freeport.GetFreePorts(2)
	assert.NoError(t, err)

	relayPort := ports[0]
	accessPort := ports[1]

	audience := "http://[::]:" + strconv.Itoa(accessPort)
	target := "ws://127.0.0.1:" + strconv.Itoa(relayPort)

	secret := "testsecret"
	base_path := "/api/v1"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	time.Sleep(3 * time.Second)

	//TODO

	// consider adding some clients and sending some messages to check the stats of actual
	// active clients (or check we test this elsewhere)

	assert.NoError(t, err)

	var claims permission.Token

	start := jwt.NewNumericDate(time.Now().Add(-time.Second))
	afterTTL := jwt.NewNumericDate(time.Now().Add(time.Duration(60) * time.Second))
	claims.IssuedAt = start
	claims.NotBefore = start
	claims.ExpiresAt = afterTTL
	claims.Audience = jwt.ClaimStrings{audience}
	claims.Topic = "stats"
	claims.ConnectionType = "connect"
	claims.Scopes = []string{"stats"}

	rawtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	token, err := rawtoken.SignedString([]byte(secret))

	assert.NoError(t, err)

	s := New()
	to := audience + base_path + "/connect/stats"
	go s.Connect(ctx, to, token)

	select {
	case <-time.After(2 * time.Second):
		t.Error("did not receive status report in time")
	case report := <-s.Status:

		if debug {
			fmt.Printf("status: %+v\n", report)
		}

		// check we got a real set of reports, by checking for the existence
		// of the optics that currently apply
		ta := make(map[string]bool)

		for _, r := range report {
			ta[r.Topic] = true
		}

		te := map[string]bool{
			"stats": true,
		}

		assert.Equal(t, te, ta)

	}

}

func TestMocking(t *testing.T) {

	s := New()

	tr := []Report{Report{Topic: "test00"}}
	go func() {
		s.Status <- tr
	}()

	select {
	case <-time.After(time.Second):
		t.Error("did not receive report")

	case ar := <-s.Status:

		assert.Equal(t, tr, ar)
	}

}
