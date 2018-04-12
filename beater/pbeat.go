package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/visheyra/pbeat/config"
)

type Pbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Pbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

func (bt *Pbeat) Run(b *beat.Beat) error {
	logp.Info("pbeat is running! Hit CTRL-C to stop it.")

	srv := NewServer()
	chn := make(chan common.MapStr)
	go srv.StartServer(chn)
	return nil
}

func (bt *Pbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
