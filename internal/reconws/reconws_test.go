package reconws

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func init() {

	log.SetLevel(log.WarnLevel)

}

func TestBackoff(t *testing.T) {

	b := &backoff.Backoff{
		Min:    time.Second,
		Max:    30 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	lowerBound := []float64{0.5, 1.5, 3.5, 7.5, 15.5, 29.5, 29.5, 29.5}
	upperBound := []float64{1.5, 2.5, 4.5, 8.5, 16.5, 30.5, 30.5, 30.5}

	for i := 0; i < len(lowerBound); i++ {

		actual := big.NewFloat(b.Duration().Seconds())

		if actual.Cmp(big.NewFloat(upperBound[i])) > 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted <%f\n", i, actual, upperBound[i])
		}
		if actual.Cmp(big.NewFloat(lowerBound[i])) < 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted >%f\n", i, actual, lowerBound[i])
		}

	}

}

func TestReconnectAuth(t *testing.T) {

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

	// todo check token in header
	// send uri with code0
	// (check it connnects to the uri with that code)
	// (drop that connection)
	// repeat with new code, for uri with call-response handler
	// check messages sent and received

	token := "xx22yy33"
	codes := []string{"098mqv", "nb08gxx", "hgf32a"}
	path := "/connect/host00"
	msg0 := []byte(`apples`)
	msg1 := []byte(`bananas`)

	// relay checks path and code
	// then expects to receive the message apples,
	// then sends the message bananas
	// then disconnects

	idx := 0 //connection count at access

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))

	relay := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, path, r.URL.Path)

		q := r.URL.Query()

		if actualCode, ok := q["code"]; ok {
			assert.Equal(t, codes[idx-1], actualCode[0])
		} else {
			t.Error("no code in query")
		}

		if debug {
			t.Logf("relay url contained: path=%s code=%s\n", r.URL.Path, q["code"])
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		mt, message, err := c.ReadMessage()
		assert.NoError(t, err)

		assert.Equal(t, message, msg0)
		err = c.WriteMessage(mt, msg1)
		assert.NoError(t, err)

		if idx >= len(codes) {
			cancel() //stop the reconnecting client
		}
	}))
	defer relay.Close()

	// Create test server with the echo handler.
	access := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, token, r.Header.Get("Authorization"))

		URL, err := url.Parse(relay.URL)

		assert.NoError(t,err)

		URL.Scheme = "ws"
		URL.Path = path

		reply := Reply{
			URI: URL.String() + "?code=" + codes[idx], //update code each time
		}
		idx += 1
		rb, err := json.Marshal(reply)
		assert.NoError(t, err)
		fmt.Fprintln(w, string(rb))
		if debug {
			t.Log(string(rb))
		}
	}))
	defer access.Close()

	s0 := New()
	go s0.ReconnectAuth(ctx, access.URL+path, token)

	var timeout = time.Duration(time.Second)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second) //wait for reconnection
		select {
		case s0.Out <- WsMessage{Data: msg0, Type: websocket.BinaryMessage}:
		case <-time.After(timeout):
			t.Errorf("timeout sending message %d", i)
		}
		select {
		case msg := <-s0.In:
			assert.Equal(t, msg1, msg.Data)
		case <-time.After(timeout):
			t.Errorf("timeout receiving message %d", i)
		}
	}

	<-ctx.Done()

	if idx < len(codes) {
		t.Errorf("connection count was %d, expected %d", idx, len(codes))
	}

}

func TestWsEcho(t *testing.T) {

	r := New()

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go r.Reconnect(ctx, u)

	payload := []byte("Hello")
	mtype := int(websocket.TextMessage)

	r.Out <- WsMessage{Data: payload, Type: mtype}

	reply := <-r.In

	if !bytes.Equal(reply.Data, payload) {
		t.Errorf("Got unexpected response: %s, wanted %s\n", reply.Data, payload)
	}

	time.Sleep(2 * time.Second)

}

func TestRetryTiming(t *testing.T) {

	suppressLog()
	defer displayLog()

	r := New()

	c := make(chan int)

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deny(w, r, c)
	}))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	url := "ws" + strings.TrimPrefix(s.URL, "http")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Reconnect(ctx, url)

	// first failed connection should be immediate
	// backoff with jitter means we quite can't be sure what the timings are
	lowerBound := []float64{0.0, 0.9, 1.9, 3.9, 7.9, 9.9, 9.9, 9.9}
	upperBound := []float64{1.0, 1.1, 2.1, 4.1, 8.1, 10.1, 10.1, 10.1}

	iterations := len(lowerBound)

	if testing.Short() {
		fmt.Println("Reducing length of test in short mode")
		iterations = 3
	}

	if testing.Verbose() {
		fmt.Println("lower < actual < upper ok?")
	}

	for i := 0; i < iterations; i++ {

		start := time.Now()

		<-c // wait for deny handler to return a value (note: bad handshake due to use of deny handler)

		actual := big.NewFloat(time.Since(start).Seconds())
		ok := true

		if actual.Cmp(big.NewFloat(upperBound[i])) > 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted <%f\n", i, actual, upperBound[i])
			ok = false
		}
		if actual.Cmp(big.NewFloat(lowerBound[i])) < 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted >%f\n", i, actual, lowerBound[i])
			ok = false
		}

		if testing.Verbose() {
			fmt.Printf("%0.2f < %0.2f < %0.2f %s\n", lowerBound[i], actual, upperBound[i], okString(ok))
		}
	}

}

func TestReconnectAfterDisconnect(t *testing.T) {

	r := New()

	c := make(chan int)

	n := 0

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		connectAfterTrying(w, r, &n, 2, c)
	}))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	url := "ws" + strings.TrimPrefix(s.URL, "http")

	ctx, cancel := context.WithCancel(context.Background()) //, time.Second)
	go r.Reconnect(ctx, url)

	// first failed connection should be immediate
	// should connect on third try
	// then next attempt after that should fail immediately
	// backoff with jitter means we quite can't be sure what the timings are
	// fail immediately, wait retry and fail, wait retry and connect, fail immediately, wait retry and fail
	lowerBound := []float64{0.0, 0.9, 1.9, 0.0, 0.9, 1.9}
	upperBound := []float64{0.1, 1.1, 2.1, 0.1, 1.1, 2.1}

	iterations := len(lowerBound)

	if testing.Short() {
		fmt.Println("Reducing length of test in short mode")
		iterations = 6
	}

	if testing.Verbose() {
		fmt.Println("lower < actual < upper ok?")
	}

	for i := 0; i < iterations; i++ {

		start := time.Now()

		<-c // wait for deny handler to return a value (note: bad handshake due to use of deny handler)

		actual := big.NewFloat(time.Since(start).Seconds())
		ok := true
		if actual.Cmp(big.NewFloat(upperBound[i])) > 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted <%f\n", i, actual, upperBound[i])
			ok = false
		}
		if actual.Cmp(big.NewFloat(lowerBound[i])) < 0 {
			t.Errorf("retry timing was incorrect, iteration %d, elapsed %f, wanted >%f\n", i, actual, lowerBound[i])
			ok = false
		}
		if testing.Verbose() {
			fmt.Printf("%0.2f < %0.2f < %0.2f %s\n", lowerBound[i], actual, upperBound[i], okString(ok))
		}
	}
	cancel()
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func deny(w http.ResponseWriter, r *http.Request, c chan int) {
	c <- 0
}

func connectAfterTrying(w http.ResponseWriter, r *http.Request, n *int, connectAt int, c chan int) {

	defer func() { *n++ }()

	c <- *n

	if *n == connectAt {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer conn.Close()

		// immediately close
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	}
}

func suppressLog() {
	var ignore bytes.Buffer
	logignore := bufio.NewWriter(&ignore)
	log.SetOutput(logignore)
}

func displayLog() {
	log.SetOutput(os.Stdout)
}

func okString(ok bool) string {
	if ok {
		return "  ok"
	}
	return "  FAILED"
}
