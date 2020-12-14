package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type Producer struct {
	ChatProducer sarama.SyncProducer
}

func NewProducer() *Producer {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	c, err := sarama.NewSyncProducer([]string{
		"kafka:9092",
	}, config)
	if err != nil {
		panic(err)
	}
	return &Producer{ChatProducer: c}
}

func (p *Producer) Close() error {
	err := p.ChatProducer.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Producer) SendStringData(message string) error {
	partition, offset, err := p.ChatProducer.SendMessage(&sarama.ProducerMessage{
		Topic: "peter-topic",
		Value: sarama.StringEncoder(message),
	})
	if err != nil {
		return err
	}
	fmt.Printf("%d/%d\n", partition, offset)
	return nil
}
func main() {
	p := NewProducer()
	for i := 0; i < 30; i++ {
		err := p.SendStringData(fmt.Sprintf("Message #%d\n", i))
		if err != nil {
			return
		}
		fmt.Println("간다간다", i)
	}
}
