package beater

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/visheyra/pbeat/config"
)

//Pbeat : Beater structure
type Pbeat struct {
	config config.Config
	client beat.Client
}

//New : Creates new beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Pbeat{
		config: c,
	}
	return bt, nil
}

//Run : run the beater
func (bt *Pbeat) Run(b *beat.Beat) error {
	logp.Info("pbeat is running! Hit CTRL-C to stop it.")

	client, err := b.Publisher.Connect()
	if err != nil {
		fmt.Println("Failed to establish connection to ES, ...")
		fmt.Println("Aborting")
		return err
	}

	//Create prom server
	srv := NewServer()
	chn := make(chan beat.Event)
	defer close(chn)

	go srv.StartServer(chn)

	// Wait for metrics
	for {
		select {
		case event := <-chn:
			fmt.Println("New event received")
			client.Publish(event)
		}
	}
}

//Stop : stop the beater
func (bt *Pbeat) Stop() {
	bt.client.Close()
}
