package shellbar

import (
	"sync"
	"time"

	"github.com/eclesh/welford"
	"github.com/gorilla/websocket"
	"github.com/practable/jump/internal/ttlcode"
)

//Config represents configuration of the relay & lets configuration be passed as argument to permit testing
type Config struct {

	// Listen is the listening port
	Listen int

	// Audience must match the host in token
	Audience string

	// Secret is used to validating statsTokens
	Secret string

	// ExchangeCode swaps a code for the associated Token
	CodeStore *ttlcode.CodeStore
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

	// Buffered channel of outbound messages.
	send chan message

	// string representing the path the client connected to
	topic string

	audience string

	stats *Stats

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
	Last string `json:"last"` //how many seconds ago...

	Size float64 `json:"size"`

	Fps float64 `json:"fps"`
}

// ClientReport represents statistics on a client, and omits non-serialisable internal references
type ClientReport struct {
	Topic string `json:"topic"`

	CanRead bool `json:"canRead"`

	CanWrite bool `json:"canWrite"`

	Connected string `json:"connected"`

	RemoteAddr string `json:"remoteAddr"`

	UserAgent string `json:"userAgent"`

	Stats RxTx `json:"stats"`
}

// StatsCommand represents a command relating to collection of statistics
type StatsCommand struct {
	Command string `json:"cmd"`
}

// Stats represents statistics for (video) frames received and transmitted
type Stats struct {
	connectedAt time.Time

	rx *Frames

	tx *Frames
}

// Frames represents statistics for (video) frames
type Frames struct {
	last time.Time

	size *welford.Stats

	ns *welford.Stats

	mu *sync.RWMutex
}

// messages will be wrapped in this struct for muxing
type message struct {
	sender Client
	mt     int
	data   []byte //text data are converted to/from bytes as needed
}

// TODO - remove unused types below this line (some still in use)

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
