package shellbar

import (
	"sync"
	"time"

	"github.com/practable/jump/internal/ttlcode"
)

//Config represents configuration of the relay & lets configuration be passed as argument to permit testing
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

// Shellbar runs ssh relay with the given configuration
func Shellbar(closed <-chan struct{}, parentwg *sync.WaitGroup, config Config) {

	var wg sync.WaitGroup

	messagesToDistribute := make(chan message, config.BufferSize) //TODO make buffer length separately configurable from incoming buffer size

	var topics topicDirectory

	topics.directory = make(map[string][]clientDetails)

	clientActionsChan := make(chan clientAction)

	wg.Add(2)

	go handleConnections(closed, &wg, clientActionsChan, messagesToDistribute, config)

	go handleClients(closed, &wg, &topics, clientActionsChan)

	wg.Wait()

	parentwg.Done()

}
