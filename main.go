package main

import (
	"fmt"
	"time"

	"github.com/lazy_kafka/lazykafka"
)

func main() {
	conf := lazykafka.NewConfig()
	b := lazykafka.NewBroker("127.0.0.1:9092")
	b.Open(conf)
	fmt.Println(b.Connected())

	retention := "-1"

	req := &lazykafka.CreateTopicsRequest{
		TopicDetails: map[string]*lazykafka.TopicDetail{
			"topic_via_tcp": {
				NumPartitions:     -1,
				ReplicationFactor: -1,
				ReplicaAssignment: map[int32][]int32{
					0: []int32{0, 1, 2},
				},
				ConfigEntries: map[string]*string{
					"retention.ms": &retention,
				},
			},
		},
		Timeout: 100 * time.Millisecond,
	}
	b.CreateTopics(req)

	b.Close()

}
