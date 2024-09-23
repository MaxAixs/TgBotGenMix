package event_consumer

import (
	"BotMixology/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	bathSize  int
}

func NewConsumer(fetcher events.Fetcher, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %v", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("Cant handle event: %v", err)
		}
	}

	return nil
}
