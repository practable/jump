// Package shellclient provides a client which listens on a local tcp port and for each incoming
// connection makes a unique connection to a remote shellrelay
package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/practable/jump/internal/ws"
	log "github.com/sirupsen/logrus"
)

type Reply struct {
	URI string `json:"uri"`
}

// Client runs a listener for local ssh connections on a local tcp port
// and for each incoming connection makes a unique connection to a remote shellrelay
// Run one instance per host_id you wish to connect to (i.e. assign a unique local port to each host_id)
func Run(ctx context.Context, listen int, remote, token string) {

	id := "client.Run()"

	log.WithFields(log.Fields{"listen": listen, "remote": remote}).Infof("%s: starting", id)

	uri := ":" + strconv.Itoa(listen)

	lc := &net.ListenConfig{}

	l, err := lc.Listen(ctx, "tcp", uri)

	if err != nil {
		log.WithFields(log.Fields{"uri": uri, "err": err.Error()}).Fatalf("%s failed to create listener", id)
		return
	}

	defer l.Close()

	log.WithField("uri", uri).Infof("%s awaiting connections", id)
	defer log.WithFields(log.Fields{"listen": listen, "remote": remote}).Infof("%s: done", id)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Wait for a connection.
			tconn, err := l.Accept()

			if err != nil {
				log.WithFields(log.Fields{"uri": uri, "err": err.Error()}).Errorf("%s failed to accept connection", id)
				continue //the context is probably cancelled if we get to here, but don't return here, in case context ok and we have other live connections
			}
			// Handle the connection in a new goroutine.

			log.WithField("uri", uri).Tracef("%s got a new connection", id)

			// make websocket connection and relay to/from conn

			go func() {

				// first we make a request to the access point
				var client = &http.Client{
					Timeout: time.Second * 10,
				}

				req, err := http.NewRequest("POST", remote, nil)
				if err != nil {
					log.WithFields(log.Fields{"err": err.Error(), "remote": remote}).Warnf("%s: failed to create request", id)
					return
				}

				req.Header.Add("Authorization", token)

				resp, err := client.Do(req)

				if err != nil {
					log.WithField("err", err).Warnf("%s: failed request to access endpoint", id)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.WithField("err", err).Warnf("%s: failed reading access response body", id)
					return
				}

				var reply Reply

				err = json.Unmarshal(body, &reply)

				if err != nil {

					log.WithFields(log.Fields{"err": err, "body": string(body)}).Debugf("%s: failed marshalling access response into struct", id)
					return
				}

				wsURL, err := url.Parse(reply.URI)

				if err != nil {
					log.WithFields(log.Fields{"err": err, "body": string(body)}).Debugf("%s: failed marshalling access response into struct", id)
					return
				}

				log.Infof("%s: successful relay access request", id)

				wconn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
				if err != nil {
					log.WithField("err", err.Error()).Errorf("%s dial error", id)
					return
				}

				defer wconn.Close()

				ws.PipeBinaryIgnoreText(ctx, tconn, wconn)
			}()

		}
	}

}
