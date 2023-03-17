/* package shellbar is a crossbar-style relay
for connecting ssh connections transported
over websocket
*/

package shellbar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/practable/jump/internal/permission"
	"github.com/practable/jump/internal/ttlcode"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer (64kB)
	// Max MTU for ssh typically 65535 or 0xffff
	maxMessageSize = 0xffff
)

// When transferring files, messages are typically at max size
// So we can save some syscalls if we can fit them into the buffer
// null subprotocol required by Chrome
// TODO restrict CheckOrigin as required
var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	Subprotocols:    []string{"null"},
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//Config represents the configuration of a shellbar
type Config struct {

	// Audience must match the host in token
	Audience string

	// BufferSize is the channel buffer size for clients
	BufferSize int64

	// ExchangeCode swaps a code for the associated Token
	CodeStore *ttlcode.CodeStore

	// Listen is the listening port
	Listen int

	// Secret is used to validating statsTokens
	Secret string

	//StatsEvery controls how often to send stats reports
	StatsEvery time.Duration
}

// NewDefaultConfig returns a pointer to a new, default, Config
func NewDefaultConfig() *Config {
	c := &Config{}
	c.Listen = 3000
	c.CodeStore = ttlcode.NewDefaultCodeStore()
	return c
}

// WithListen sets the listening port in the Config
func (c *Config) WithListen(listen int) *Config {
	c.Listen = listen
	return c
}

// WithAudience sets the audience in the Config
func (c *Config) WithAudience(audience string) *Config {
	c.Audience = audience
	return c
}

// WithCodeStoreTTL sets the TTL of the codestore
func (c *Config) WithCodeStoreTTL(ttl int64) *Config {
	c.CodeStore = ttlcode.NewDefaultCodeStore().
		WithTTL(ttl)
	return c
}

// ConnectionAction represents an action happening on a  connection
type ConnectionAction struct {
	Action string `json:"action"`
	URI    string `json:"uri"`
	UUID   string `json:"uuid"`
}

// Client is a middleperson between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// connectedAt represents when the current connection started
	connectedAt time.Time

	// expiresAt represents when the connection's token expires
	expiresAt time.Time

	// Buffered channel of outbound messages.
	send chan message

	// string representing the path the client connected to
	topic string

	audience string

	stats *stats

	name string

	userAgent string

	remoteAddr string

	// existence of scopes to read, write
	canRead, canWrite bool

	// hostAlertUUID is the reference we sent to the host for our unique connection
	// store it so we can tell it which connection we are closing.
	hostAlertUUID string

	// prevent clients from sending before host has sent something, which
	// is what you need for server speaks first
	// some clients rush ahead and send their ssh identification before
	// receiving the hosts, so the host never gets it because they are
	// still connecting ....
	mustWaitToSend bool

	// closed once we've received something, or immediately if !mustWaitToSend
	clearToSend chan struct{}
}

// RxTx represents statistics for both receive and transmit
type RxTx struct {
	Tx ReportStats `json:"tx"`
	Rx ReportStats `json:"rx"`
}

// ReportStats represents statistics to be reported on a connection
type ReportStats struct {
	Last string `json:"last"`
}

// ClientReport represents statistics on a client, and omits non-serialisable internal references
type ClientReport struct {
	CanRead bool `json:"can_read"`

	CanWrite bool `json:"can_write"`

	ConnectedAt string `json:"connected_at"`

	ExpiresAt string `json:"expires_at"`

	RemoteAddr string `json:"remote_address"`

	Statistics ReportStats `json:"statistcs"`

	Topic string `json:"topic"`

	UserAgent string `json:"user_agent"`
}

// StatsCommand represents a command relating to collection of statistics
type StatsCommand struct {
	Command string `json:"cmd"`
}

// Stats represents statistics about when the last transmission was made
type stats struct {
	mu   *sync.RWMutex
	last time.Time
}

// messages will be wrapped in this struct for muxing
type message struct {
	sender Client
	mt     int
	data   []byte //text data are converted to/from bytes as needed
}

type clientDetails struct {
	name         string
	topic        string
	messagesChan chan message
}

// requests to add or delete subscribers are represented by this struct
type clientAction struct {
	action clientActionType
	client clientDetails
}

// userActionType represents the type of of action requested
type clientActionType int

// clientActionType constants
const (
	clientAdd clientActionType = iota
	clientDelete
)

type topicDirectory struct {
	sync.Mutex
	directory map[string][]clientDetails
}

// Shellbar runs ssh relay with the given configuration
func Shellbar(closed <-chan struct{}, parentwg *sync.WaitGroup, config Config) {

	var wg sync.WaitGroup

	var topics topicDirectory

	topics.directory = make(map[string][]clientDetails)

	clientActionsChan := make(chan clientAction)

	wg.Add(2)

	go handleConnections(closed, &wg, config)

	go handleClients(closed, &wg, &topics, clientActionsChan)

	wg.Wait()

	parentwg.Done()

	log.Trace("Shellbar finished")

}

func fpsFromNs(ns float64) float64 {
	return 1 / (ns * 1e-9)
}

func handleConnections(closed <-chan struct{}, parentwg *sync.WaitGroup, config Config) {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(closed, hub, w, r, config)
	})

	var wg sync.WaitGroup
	wg.Add(1)

	go statsClient(closed, &wg, hub, config)

	addr := ":" + strconv.Itoa(config.Listen)

	h := &http.Server{Addr: addr, Handler: nil}

	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Info("ListenAndServe: ", err) //TODO upgrade to fatal once httptest is supported
		}
	}()

	<-closed

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := h.Shutdown(ctx)
	if err != nil {
		log.Infof("ListenAndServe.Shutdown(): %s", err.Error())
	}
	wg.Wait()
	parentwg.Done()
	log.Trace("handleConnections is done")
}

func handleClients(closed <-chan struct{}, wg *sync.WaitGroup, topics *topicDirectory, clientActionsChan chan clientAction) {

	defer func() {

		wg.Done()

		log.WithFields(log.Fields{
			"func": "HandleClients",
			"verb": "closed",
		}).Trace("HandleClients closed")

	}()

	for {
		select {
		case <-closed:
			return
		case request := <-clientActionsChan:

			log.WithField("clientAction", request).Trace("handleClients")

			if request.action == clientAdd {

				addClientToTopic(topics, request.client)

			} else if request.action == clientDelete {
				deleteClientFromTopic(topics, request.client)

			}
		}
	}
}

func addClientToTopic(topics *topicDirectory, client clientDetails) {

	_, exists := topics.directory[client.topic]

	if !exists {
		topics.Lock()
		topics.directory[client.topic] = []clientDetails{client}
		topics.Unlock()

		log.WithFields(log.Fields{
			"topic":  client.topic,
			"client": client,
			"action": clientAdd,
			"verb":   "add",
			"count":  1,
		}).Debug("Added first client to new topic")

	} else {
		topics.Lock()
		topics.directory[client.topic] = append(topics.directory[client.topic], client)
		count := len(topics.directory[client.topic])
		topics.Unlock()

		log.WithFields(log.Fields{
			"topic":  client.topic,
			"client": client,
			"action": clientAdd,
			"verb":   "add",
			"count":  count,
		}).Debug("Added client to existing topic")

	}
}

func deleteClientFromTopic(topics *topicDirectory, client clientDetails) {

	_, exists := topics.directory[client.topic]
	if exists {
		topics.Lock()
		existingClients := topics.directory[client.topic]
		topics.directory[client.topic] = filterClients(existingClients, client)
		count := len(topics.directory[client.topic])
		topics.Unlock()

		log.WithFields(log.Fields{
			"topic":  client.topic,
			"client": client,
			"action": clientDelete,
			"verb":   "delete",
			"count":  count,
		}).Debug("Deleting client from existing topic")

	} else {

		log.WithFields(log.Fields{
			"topic":  client.topic,
			"client": client,
			"action": clientDelete,
			"verb":   "delete",
			"count":  0,
		}).Debug("Ignoring: can't delete client from non-existent topic")
	}
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool

	mu *sync.RWMutex

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		mu:         &sync.RWMutex{},
		broadcast:  make(chan message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.topic]; !ok {

				h.clients[client.topic] = make(map[*Client]bool)
			}
			h.clients[client.topic][client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.topic]; ok {
				delete(h.clients[client.topic], client)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			topic := message.sender.topic
			h.mu.RLock()
			for client := range h.clients[topic] {
				if client.name != message.sender.name {
					select {
					case client.send <- message:
					default:
						h.unregister <- client
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {

	id := "shellbar.readPump(" + c.topic + "/" + c.name + ")"

	defer func() {
		log.Tracef("%s.defer(): about to disconnect", id)
		// Tell the host that we have gone ...

		// alert SSH host agent to make a new connection to relay at the same address
		// No stats needed because we are not registering to receive messages
		adminClient := &Client{
			topic: getHostTopicFromUniqueTopic(c.topic),
			name:  uuid.New().String(),
		}

		ca := ConnectionAction{
			Action: "disconnect",
			UUID:   c.hostAlertUUID,
		}

		camsg, err := json.Marshal(ca)

		if err != nil {
			log.WithFields(log.Fields{"error": err, "uuid": c.hostAlertUUID}).Errorf("%s.defer(): Failed to make disconnect connectionAction message because %s", id, err.Error())
			return
		}

		c.hub.broadcast <- message{sender: *adminClient, data: camsg, mt: websocket.TextMessage}
		log.Tracef("%s.defer(): broadcast disconnect of UUID %s", id, c.hostAlertUUID)

		c.hub.unregister <- c
		log.Tracef("%s.defer(): client unregistered", id)

		c.conn.Close()
		log.Tracef("%s.defer(): DONE", id)

	}()

	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Errorf("readPump deadline error: %v", err)
		return
	}

	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})

	for {

		mt, data, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Tracef("Unexpected websocket close: %v", err)
			}
			break
		}

		size := len(data)

		if c.mustWaitToSend {
			<-c.clearToSend
		}

		if c.canWrite {

			c.hub.broadcast <- message{sender: *c, data: data, mt: mt}

			log.WithFields(log.Fields{"topic": c.topic, "size": size}).Tracef("%s: broadacast %d-byte message to topic %s", id, size, c.topic)

			c.stats.mu.Lock()
			c.stats.last = time.Now()
			c.stats.mu.Unlock()

		} else {
			log.WithFields(log.Fields{"topic": c.topic, "size": size}).Tracef("%s: ignored %d-byte message intended for broadcast to topic %s", id, size, c.topic)

		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(closed <-chan struct{}, cancelled <-chan struct{}) {

	id := "shellbar.writePump(" + c.topic + "/" + c.name + ")"

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		log.Tracef("%s: done", id)
	}()
	log.Tracef("%s: starting", id)

	awaitingFirstMessage := true

	for {

		select {

		case message, ok := <-c.send:

			if awaitingFirstMessage {
				close(c.clearToSend)
				awaitingFirstMessage = false
			}

			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Errorf("%s: writePump deadline error: %s", id, err.Error())
				return
			}

			if !ok {
				// The hub closed the channel.
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Errorf("%s: writePump closeMessage error: %s", id, err.Error())
				}
				return
			}

			if c.canRead { //only send if authorised to read

				w, err := c.conn.NextWriter(message.mt)
				if err != nil {
					return
				}

				n, err := w.Write(message.data)

				if err != nil {
					log.Errorf("writePump writing error: %v", err)
				}

				size := len(message.data)

				if n != size {
					log.Errorf("writePump incomplete write %d of %d", n, size)
				}

				log.WithFields(log.Fields{"topic": c.topic, "size": size}).Tracef("%s: wrote %d-byte message from topic %s", id, size, c.topic)

				// don't queue chunks; makes reading JSON objects on the host connectAction channel fail if two connects happen together
				// don't record stats on messages received (we already recorded them on what was transmitting the message)

				if err := w.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Errorf("%s: writePump ping deadline error: %s", id, err.Error())
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Warnf("%s: done because conn error %s", id, err.Error())
				return
			}
		case <-closed:
			log.Tracef("%s: done because closed channel closed", id)
			return
		case <-cancelled:
			log.Tracef("%s: done because cancelled channel closed", id)
			return
		}
	}
}

// ConnectionType represents whether the connection is for Session, Shell or Unsupported
type ConnectionType int

// Enumerated connection types
const (
	Session ConnectionType = iota
	Shell
	Unsupported
)

// serveWs handles websocket requests from clients.
func serveWs(closed <-chan struct{}, hub *Hub, w http.ResponseWriter, r *http.Request, config Config) {

	id := "shellbar.serveWs(" + uuid.New().String()[0:6] + ")"

	// check if topic is of a supported type before we go any further
	ct := Unsupported

	path := slashify(r.URL.Path)

	connectionType := getConnectionTypeFromPath(path)
	topic := getTopicFromPath(path)

	if connectionType == "shell" {
		ct = Shell
	}

	if ct == Unsupported {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.WithField("connectionType", connectionType).Errorf("%s: connectionType %s unsupported", id, connectionType)
		return
	}

	log.WithFields(log.Fields{"path": r.URL.Path}).Infof("%s: received %s connection to topic %s at %s", id, connectionType, topic, r.URL.Path)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{"path": r.URL.Path, "error": err}).Errorf("%s: failed to upgrade to websocket at %s", id, r.URL.Path)
		return
	}

	//Cannot return any http responses from here on

	// Enforce permissions by exchanging the authcode for a connection ticket
	// which contains expiry time, route, and permissions

	// Get the first code query param, lowercase only
	var code string

	code = r.URL.Query().Get("code")

	// if no code or empty, return 401
	if code == "" {
		log.WithField("url", r.URL.String()).Infof("%s: Unauthorized - No Code in %s", id, r.URL.String())
		return
	}

	// Exchange code for token

	token, err := config.CodeStore.ExchangeCode(code)

	if err != nil {
		log.WithFields(log.Fields{"topic": topic, "error": err}).Infof("%s: Unauthorized - Invalid Code because %s", id, err.Error())
		return
	}

	// check token is a permission token so we can process it properly
	// It's been validated so we don't need to re-do that
	if !permission.HasRequiredClaims(token) {
		log.WithField("topic", topic).Infof("%s: Unauthorized - original token missing claims", id)
		return
	}

	now := config.CodeStore.GetTime()

	if token.NotBefore.After(time.Unix(now, 0)) {
		log.WithField("topic", topic).Infof("%s: Unauthorized - Too early", id)
		return
	}

	ttl := token.ExpiresAt.Unix() - now

	log.WithFields(log.Fields{"ttl": ttl, "topic": topic}).Trace()

	audok := false

	for _, aud := range token.Audience {
		if aud == config.Audience {
			audok = true
		}
	}

	topicBad := (topic != token.Topic)
	expired := ttl < 0

	if (!audok) || topicBad || expired {
		log.WithFields(log.Fields{"audienceOK": audok, "topicOK": !topicBad, "expired": expired, "topic": topic}).Trace("Token invalid")
		return
	}

	// check permissions

	var canRead, canWrite bool

	for _, scope := range token.Scopes {
		if scope == "read" {
			canRead = true
		}
		if scope == "write" {
			canWrite = true
		}
	}

	if !(canRead || canWrite) {
		log.WithFields(log.Fields{"topic": topic, "scopes": token.Scopes}).Tracef("%s: No valid scopes", id)
		return
	}

	cancelled := make(chan struct{})

	// cancel the connection when the token has expired
	go func() {
		time.Sleep(time.Duration(ttl) * time.Second)
		close(cancelled)
	}()

	if ct == Shell {
		// initialise statistics
		stats := &stats{mu: &sync.RWMutex{}} //Leave last at default value

		client := &Client{hub: hub,
			audience:       config.Audience,
			canRead:        canRead,
			canWrite:       canWrite,
			clearToSend:    make(chan struct{}),
			conn:           conn,
			connectedAt:    time.Now(),
			expiresAt:      time.Unix((*token.ExpiresAt).Unix(), 0),
			hostAlertUUID:  uuid.New().String(),
			mustWaitToSend: token.AlertHost,
			name:           uuid.New().String(),
			remoteAddr:     r.Header.Get("X-Forwarded-For"),
			send:           make(chan message, config.BufferSize),
			stats:          stats,
			topic:          topic + token.TopicSalt,
			userAgent:      r.UserAgent(),
		}
		client.hub.register <- client

		log.WithField("Topic", client.topic).Tracef("%s: registering client at topic %s with name %s", id, client.topic, client.name)

		go client.writePump(closed, cancelled)
		go client.readPump()

		log.WithField("topic", topic+token.TopicSalt).Tracef("%s: started shellrelay client on topic %s", id, topic+token.TopicSalt)

		if token.AlertHost {
			log.WithField("topic", topic+token.TopicSalt).Tracef("%s: alert host of topic %s to new client %s with salt %s", id, topic, client.name, token.TopicSalt)

			// alert SSH host agent to make a new connection to relay at the same address
			// no stats required because we are not registering to receive messages
			adminClient := &Client{
				topic: getHostTopicFromUniqueTopic(topic),
				name:  uuid.New().String(),
			}

			permission.SetAlertHost(&token, false) //turn off host alert
			code = config.CodeStore.SubmitToken(token)

			if code == "" {
				log.Errorf("%s: failed to submit host connect token in exchange for a code", id)
				return
			}

			// same URL as client used, but different code (and leave out the salt)
			// note the token could have multiple audiences, whereas we are checking validity against
			// only one audience, the config audience, so use that to generate our URI
			hostAlertURI := config.Audience + "/" + token.ConnectionType + "/" + token.Topic + "?code=" + code
			ca := ConnectionAction{
				Action: "connect",
				URI:    hostAlertURI,
				UUID:   client.hostAlertUUID,
			}

			camsg, err := json.Marshal(ca)

			if err != nil {
				log.WithFields(log.Fields{"uuid": client.hostAlertUUID, "uri": hostAlertURI, "error": err}).Errorf("%s: Failed to make connectionAction message", id)
				return
			}

			hub.broadcast <- message{sender: *adminClient, data: camsg, mt: websocket.TextMessage}
			log.WithFields(log.Fields{"uuid": client.hostAlertUUID, "uri": hostAlertURI, "code": code}).Debugf("%s: sent host CONNECT for topic %s with UUID:%s at URI:%s", id, topic, client.hostAlertUUID, hostAlertURI)

		}

		return
	}

}

// StatsClient starts a routine which sends stats reports on demand.
func statsClient(closed <-chan struct{}, wg *sync.WaitGroup, hub *Hub, config Config) {

	stats := &stats{mu: &sync.RWMutex{}}

	client := &Client{hub: hub,
		connectedAt: time.Now(),
		send:        make(chan message, 256),
		topic:       "stats",
		stats:       stats,
		name:        "stats-generator-" + uuid.New().String(),
		audience:    config.Audience,
		userAgent:   "shellbar",
		remoteAddr:  "internal",
		canRead:     true,
		canWrite:    true,
	}
	client.hub.register <- client

	go client.statsReporter(closed, wg, config)

}

// StatsReporter sends a stats update in response to {"cmd":"update"}.
func (c *Client) statsReporter(closed <-chan struct{}, wg *sync.WaitGroup, config Config) {

	defer wg.Done()

	var sc StatsCommand

	for {

		select {
		case <-closed:
			log.Trace("StatsReporter closed")
			return
		case msg, ok := <-c.send: // received a message from hub

			if !ok {
				return //send is closed, so we are finished
			}

			err := json.Unmarshal(msg.data, &sc)

			if err != nil {
				log.WithFields(log.Fields{"error": err, "msg": string(msg.data)}).Error("statsReporter could not unmarshal into json")
			}

			log.WithField("cmd", sc.Command).Trace("statsReporter received command")

			doUpdate := false

			if sc.Command == "update" {
				doUpdate = true
			}

			// drain the channel to avoid stale requests on next iteration of the loop
			n := len(c.send)
			for i := 0; i < n; i++ {
				msg, ok = <-c.send
				if !ok {
					return //send is closed, so we are finished
				}

				err = json.Unmarshal(msg.data, &sc)

				if err != nil {
					log.WithFields(log.Fields{"error": err, "msg": string(msg.data)}).Error("statsReporter could not marshall into json")
				}

				log.WithField("cmd", sc.Command).Trace("statsReporter received command")

				if sc.Command == "update" {
					doUpdate = true
				}
			}

			log.WithField("doUpdate", doUpdate).Trace("statsReporter do update?")

			if !doUpdate { //don't send updated stats, because no command was valid
				continue
			}

		case <-time.After(config.StatsEvery):
			log.Trace("StatsReporter routine send...")
		}

		var reports []*ClientReport

		c.hub.mu.RLock()
		for _, topic := range c.hub.clients {
			for client := range topic {

				client.stats.mu.RLock()
				rs := ReportStats{
					Last: time.Since(client.stats.last).String(),
				}

				client.stats.mu.RUnlock()

				c, err := client.connectedAt.UTC().MarshalText()
				if err != nil {
					log.WithFields(log.Fields{"error": err.Error(), "topic": client.topic, "connectedAt": client.connectedAt}).Error("stats cannot marshal connectedAt time to string")
				}
				ea, err := client.expiresAt.UTC().MarshalText()
				if err != nil {
					log.WithFields(log.Fields{"error": err.Error(), "topic": client.topic, "expiresAt": client.expiresAt}).Error("stats cannot marshal expiresAt time to string")
				}

				report := &ClientReport{
					Topic:       client.topic,
					CanRead:     client.canRead,
					CanWrite:    client.canWrite,
					ConnectedAt: string(c),
					ExpiresAt:   string(ea),
					RemoteAddr:  client.remoteAddr,
					UserAgent:   client.userAgent,
					Statistics:  rs,
				}

				reports = append(reports, report)

			} //for client in topic
		} // for topic in hub
		c.hub.mu.RUnlock()
		reportsData, err := json.Marshal(reports)
		if err != nil {
			log.WithField("error", err).Error("statsReporter marshalling JSON")
			return
		}
		// broadcast stats back to the hub (i.e. and anyone listening to this topic)
		c.hub.broadcast <- message{sender: *c, data: reportsData, mt: websocket.TextMessage}

	}
}

func filterClients(clients []clientDetails, filter clientDetails) []clientDetails {
	filteredClients := clients[:0]
	for _, client := range clients {
		if client.name != filter.name {
			filteredClients = append(filteredClients, client)
		}
	}
	return filteredClients
}

func slashify(path string) string {

	//remove trailing slash (that's for directories)
	path = strings.TrimSuffix(path, "/")

	//ensure leading slash without needing it in config
	path = strings.TrimPrefix(path, "/")
	path = fmt.Sprintf("/%s", path)

	return path

}

func getHostTopicFromUniqueTopic(topic string) string {

	re := regexp.MustCompile(`^([\w\%-]*)`)

	matches := re.FindStringSubmatch(topic)

	if len(matches) < 2 {
		return ""
	}

	// matches[0] = "/{prefix}/"
	// matches[1] = "{prefix}"
	return matches[1]
}

func getConnectionTypeFromPath(path string) string {

	re := regexp.MustCompile(`^\/([\w\%-]*)`)

	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		return ""
	}

	// matches[0] = "/{prefix}/"
	// matches[1] = "{prefix}"
	return matches[1]
}

func getTopicFromPath(path string) string {

	re := regexp.MustCompile(`^\/[\w\%-]*\/([\w\%-\/]*)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func getShellIDFromPath(path string) string {

	re := regexp.MustCompile(`^\/[\w\%-]*\/([\w\%-]*)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func getConnectionIDFromPath(path string) string {

	re := regexp.MustCompile(`^\/(?:([\w\%-]*)\/){2}([\w\%-]*)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		return ""
	}

	return matches[2]
}
