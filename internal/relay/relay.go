package relay

import (
	"context"
	"time"

	"github.com/practable/jump/internal/access"
	"github.com/practable/jump/internal/crossbar"
	"github.com/practable/jump/internal/ttlcode"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AccessPort     int
	Audience       string
	BufferSize     int64
	ConnectionType string
	RelayPort      int
	Secret         string
	StatsEvery     time.Duration
	Target         string
}

// Relay runs a shellrelay instance that relays ssh connections between shellclient and shellhost
func Relay(ctx context.Context, config Config) {

	cs := ttlcode.NewDefaultCodeStore()

	if config.BufferSize < 1 || config.BufferSize > 512 {
		log.WithFields(log.Fields{"requested": config.BufferSize, "actual": 256}).Warn("Overriding configured buffer size because out of range 1-512")
		config.BufferSize = 256
	}

	if config.StatsEvery < time.Duration(time.Second) {
		log.WithFields(log.Fields{"requested": config.StatsEvery, "actual": "1s"}).Warn("Overriding configured stats every because smaller than 1s")
		config.StatsEvery = time.Duration(time.Second) //we have to balance fast testing vs high CPU load in production if too short
	}

	sc := bar.Config{
		Audience:       config.Target,
		BufferSize:     config.BufferSize,
		CodeStore:      cs,
		Listen:         config.RelayPort,
		ConnectionType: config.ConnectionType,
		StatsEvery:     config.StatsEvery,
	}

	go crossbar.CrossBar(ctx, sc)

	ac := access.Config{
		Audience:  config.Audience,
		CodeStore: cs,
		Listen:    config.AccessPort,
		Secret:    config.Secret,
		Target:    config.Target,
	}

	go access.API(context, ac)

	<-ctx.Done()
	log.Trace("Relay done")
}
