package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Consumer struct {
	Consumer          sarama.Consumer
	PartitionConsumer [2]sarama.PartitionConsumer
}

func NewConsumer() (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Retry.Max = 0
	consumer, err := sarama.NewConsumer([]string{
		"127.0.0.1:9092"}, config)
	if err != nil {
		panic(err)
	}
	c := &Consumer{Consumer: consumer}

	partitions, err := consumer.Partitions("peter-topic")
	if err != nil {
		panic(err)
	}

	for i, v := range partitions {
		c.PartitionConsumer[i], err = consumer.ConsumePartition("peter-topic", v, sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
	}

	return c, nil
}

func (c *Consumer) Run() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signals
		fmt.Println()
		fmt.Println(sig)
		panic(sig)
	}()

	for _, v := range c.PartitionConsumer {
		if v == nil {
			break
		}
		go func(v sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-v.Messages():
					fmt.Println("Message   : ", string(msg.Value))
					fmt.Println("Partition : ", msg.Partition)
					fmt.Println("Offset    : ", msg.Offset)
				}
			}
		}(v)
	}
	time.Sleep(time.Second * 3600)
}

func main() {
	consumer, err := NewConsumer()
	if err != nil {
		panic(err)
	}
	consumer.Run()
}
