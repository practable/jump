package access

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/phayes/freeport"
	"github.com/practable/jump/internal/access/restapi/operations"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/ttlcode"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetPrefixFromPath(t *testing.T) {

	assert.Equal(t, "foo%20bar", getPrefixFromPath("/foo%20bar/glum"))
	assert.Equal(t, "", getPrefixFromPath("ooops/foo%20bar/glum"))

}

func TestTokenGeneration(t *testing.T) {

	iat := int64(1609329233)
	nbf := int64(1609329233)
	exp := int64(1609330233)
	audience := "https://relay-access.example.io"
	ct := "connect"
	topic := "f7558de0-cb0d-4cb5-9518-ac71d044800b"
	scopes := []string{"host"}
	secret := "somesecret"

	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b3BpYyI6ImY3NTU4ZGUwLWNiMGQtNGNiNS05NTE4LWFjNzFkMDQ0ODAwYiIsInByZWZpeCI6ImNvbm5lY3QiLCJzY29wZXMiOlsiaG9zdCJdLCJhdWQiOlsiaHR0cHM6Ly9yZWxheS1hY2Nlc3MuZXhhbXBsZS5pbyJdLCJleHAiOjE2MDkzMzAyMzMsIm5iZiI6MTYwOTMyOTIzMywiaWF0IjoxNjA5MzI5MjMzfQ.AkwDAejQx2_5WJpOxbGLQqAV_6oGlBs1hFTaL4WKu8o"

	bearer, err := Token(audience, ct, topic, secret, scopes, iat, nbf, exp)

	assert.NoError(t, err)
	assert.Equal(t, expected, bearer)

}

func TestAPI(t *testing.T) {

	debug := true

	if debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
		defer log.SetOutput(os.Stdout)

	} else {
		var ignore bytes.Buffer
		logignore := bufio.NewWriter(&ignore)
		log.SetOutput(logignore)
	}

	closed := make(chan struct{})
	var wg sync.WaitGroup

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	base_path := "/api/v1"
	secret := "testsecret"

	audience := "http://[::]:" + strconv.Itoa(port)
	cs := ttlcode.NewDefaultCodeStore()
	target := "wss://relay.example.io"

	wg.Add(1)

	config := Config{
		CodeStore: cs,
		Listen:    port,
		Audience:  audience,
		Secret:    secret,
		Target:    target,
	}

	go API(closed, &wg, config) //port, audience, secret, target, cs)

	time.Sleep(100 * time.Millisecond)

	client := &http.Client{}

	// Start tests
	req, err := http.NewRequest("POST", audience+base_path+"/connect/123", nil)
	assert.NoError(t, err)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string([]byte(body))
	assert.Equal(t, `{"code":401,"message":"unauthenticated for invalid credentials"}`, bodyStr)

	var claims permission.Token

	start := jwt.NewNumericDate(time.Now().Add(-time.Second))
	after5 := jwt.NewNumericDate(time.Now().Add(5 * time.Second))
	claims.IssuedAt = start
	claims.NotBefore = start
	claims.ExpiresAt = after5

	claims.Audience = jwt.ClaimStrings{audience}
	claims.Topic = "123"
	claims.ConnectionType = "connect"
	claims.Scopes = []string{"read", "write"} //Wrong scopes for shell, deliberately....

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	bearer, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", audience+base_path+"/connect/123", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", bearer)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, `{"code":"401","message":"missing client, host or stats scope"}`+"\n", string(body))

	// now try with correct scopes :-)
	claims.Scopes = []string{"host"}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	bearer, err = token.SignedString([]byte(secret))
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", audience+base_path+"/connect/123", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", bearer)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	body, _ = ioutil.ReadAll(resp.Body)

	var p operations.ConnectOKBody

	err = json.Unmarshal(body, &p)
	assert.NoError(t, err)
	if err != nil {
		t.Log(string(body))
	}

	expected := "wss://relay.example.io/connect/123?code="

	if len(p.URI) < len(expected) {
		t.Log(p.URI)
		t.Fatal("URI too short")
	} else {
		assert.Equal(t, expected, p.URI[0:len(expected)])
	}

	// Now repeat with client, expecting to get a connectionID added to the uri ...

	claims.Scopes = []string{"client"}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	bearer, err = token.SignedString([]byte(secret))
	assert.NoError(t, err)

	req, err = http.NewRequest("POST", audience+base_path+"/connect/123", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", bearer)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	body, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &p)
	assert.NoError(t, err)
	if err != nil {
		t.Log(string(body))
	}

	expected = "wss://relay.example.io/connect/123/" //has a unique element after this, not the word code

	if len(p.URI) < len(expected) {
		t.Fatal("URI too short")
	} else {
		assert.Equal(t, expected, p.URI[0:len(expected)])
	}

	re := regexp.MustCompile(`wss:\/\/relay\.example\.io\/connect\/123\/([\w-\%]*)\?code=.*`)
	matches := re.FindStringSubmatch(p.URI)
	assert.Equal(t, 36, len(matches[1])) //length of a UUID

	uniqueConnection0 := matches[1]

	// repeat the request, and check we get a different connection id

	req, err = http.NewRequest("POST", audience+base_path+"/connect/123", nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", bearer)

	resp, err = client.Do(req)
	assert.NoError(t, err)
	body, _ = ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &p)
	assert.NoError(t, err)
	if err != nil {
		t.Log(string(body))
	}

	expected = "wss://relay.example.io/connect/123/"

	if len(p.URI) < len(expected) {
		t.Fatal("URI too short")
	} else {
		assert.Equal(t, expected, p.URI[0:len(expected)])
	}

	matches = re.FindStringSubmatch(p.URI)
	assert.Equal(t, 36, len(matches[1])) //length of a UUID

	assert.NotEqual(t, uniqueConnection0, matches[1])

	// End tests
	close(closed)
	wg.Wait()
}
