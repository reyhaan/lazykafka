package main

import "github.com/lazy_kafka/lazykafka"

func main() {
	conf := lazykafka.NewConfig()
	b := lazykafka.NewBroker("127.0.0.1:9092")
	b.Open(conf)
}
