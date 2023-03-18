// Package shellhost manages connections between
//  a shellrelay and a host machine accepting login shell
// connections, such as sshd. Such protocols are server speak first
// so the relay alerts shellhost when a new client has connected
// to the relay. Shellhost then makes a new dediated connection
// to the login shell port, and the relay.
package host

import (
	"context"
	"encoding/json"
	"net"

	"net/url"

	"github.com/gorilla/websocket"
	"github.com/practable/jump/internal/crossbar"
	"github.com/practable/jump/internal/reconws"
	"github.com/practable/jump/internal/ws"
	log "github.com/sirupsen/logrus"
)

// Host connects to remote relay, and makes a new connection
// to local (localhost:{port}) every time it is alerted to a new
// connection by a shellbar.ConnectionAction
func Run(ctx context.Context, local, remote, token string) {

	id := "host()"

	log.WithFields(log.Fields{"local": local, "remote": remote}).Infof("%s: starting", id)

	manager := reconws.New()
	go manager.ReconnectAuth(ctx, remote, token) // the manager works fine with reconws

	connections := make(map[string]context.CancelFunc)
	var ca crossbar.ConnectionAction

	for {
		select {
		case <-ctx.Done():
			log.Fatalf("%s: shutdown because context cancelled", id)
			for _, cancel := range connections {
				cancel()
			}
			return

		case msg, ok := <-manager.In:
			log.WithField("msg", string(msg.Data)).Debugf("%s: control message received", id)

			if !ok {
				log.Fatalf("%s: shutting down because control channel closed unexpectedly", id)
				for _, cancel := range connections {
					cancel()
				}
				return
			}

			err := json.Unmarshal(msg.Data, &ca)
			if err != nil {
				log.WithField("err", err.Error()).Warnf("%s: ignoring control message because unmarshal error", id)
				continue
			}

			_, err = url.ParseRequestURI(ca.URI)
			if err != nil {
				log.WithField("err", err.Error()).Warnf("%s: ignoring control message because URI parse error", id)
				continue
			}

			switch ca.Action {
			case "connect":
				uCtx, uCancel := context.WithCancel(ctx)
				connections[ca.UUID] = uCancel
				log.WithFields(log.Fields{"local": local, "uri": ca.URI, "uuid": ca.UUID}).Infof("%s: received connect command", id)
				go newConnection(uCtx, local, ca.URI, ca.UUID)

			case "disconnect":
				uCancel, ok := connections[ca.UUID]
				if !ok {
					log.WithFields(log.Fields{"local": local, "uri": ca.URI, "uuid": ca.UUID}).Warnf("%s: ignoring disconnect command because connection does not exist", id)
					continue
				}
				log.WithFields(log.Fields{"local": local, "uri": ca.URI, "uuid": ca.UUID}).Infof("%s: received disconnect command", id)
				uCancel()
			}
		}
	}

}

func newConnection(ctx context.Context, local, remote, uuid string) {

	id := "host.newConnection(" + uuid + ")"

	log.WithFields(log.Fields{"local": local, "remote": remote}).Debugf("%s: starting", id)

	// connect to websocket first, so that we can pass the hello message from ssh server
	// to the client at the far end
	wconn, _, err := websocket.DefaultDialer.Dial(remote, nil)
	if err != nil {
		log.Errorf("%s: ws dial error  %s", id, err)
	}

	defer wconn.Close()

	tconn, err := net.Dial("tcp", local)
	if err != nil {
		log.Errorf("%s ws dial error was %s", id, err.Error())
		return
	}

	defer tconn.Close()

	log.WithFields(log.Fields{"local": local, "remote": remote}).Infof("%s: started", id)

	ws.PipeBinaryIgnoreText(ctx, tconn, wconn)

	<-ctx.Done()
	log.WithFields(log.Fields{"local": local, "remote": remote}).Infof("%s: done", id)
}
