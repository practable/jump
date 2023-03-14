package shellrelay

import (
	"sync"
	"time"

	"github.com/practable/jump/internal/shellaccess"
	"github.com/practable/jump/internal/shellbar"
	"github.com/practable/jump/internal/ttlcode"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AccessPort int
	Audience   string
	BufferSize int64
	RelayPort  int
	Secret     string
	StatsEvery time.Duration
	Target     string
}

// Relay runs a shellrelay instance that relays ssh connections between shellclient and shellhost
func Relay(closed <-chan struct{}, parentwg *sync.WaitGroup, config Config) { // accessPort, relayPort int, audience, secret, target string) {

	var wg sync.WaitGroup

	cs := ttlcode.NewDefaultCodeStore()

	if config.BufferSize < 1 || config.BufferSize > 512 {
		log.WithFields(log.Fields{"requested": config.BufferSize, "actual": 256}).Warn("Overriding configured buffer size because out of range 1-512")
		config.BufferSize = 256
	}

	if config.StatsEvery < time.Duration(time.Second) {
		log.WithFields(log.Fields{"requested": config.StatsEvery, "actual": "1s"}).Warn("Overriding configured stats every because smaller than 1s")
		config.StatsEvery = time.Duration(time.Second) //we have to balance fast testing vs high CPU load in production if too short
	}

	sc := shellbar.Config{
		Audience:   config.Target,
		BufferSize: config.BufferSize,
		CodeStore:  cs,
		Listen:     config.RelayPort,
		StatsEvery: config.StatsEvery,
	}

	wg.Add(1)
	go shellbar.Shellbar(closed, &wg, sc)

	ac := shellaccess.Config{
		Audience:  config.Audience,
		CodeStore: cs,
		Listen:    config.AccessPort,
		Secret:    config.Secret,
		Target:    config.Target,
	}

	wg.Add(1)
	go shellaccess.API(closed, &wg, ac)

	wg.Wait()
	parentwg.Done()
	log.Trace("Relay done")
}
