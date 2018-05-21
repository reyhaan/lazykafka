# Client for Apache Kafka written in GO

### Example:

```
// get default config or define your own config
conf := lazykafka.NewConfig()

// Create a broker target for the given address
b := lazykafka.NewBroker("127.0.0.1:9092")

// Open a connection to the targeted broker
b.Open(conf)

// Check if you are connected to the broker
fmt.Println(b.Connected())
```

Once you connect to a broker, you can start by creating a new Topic using `CreateTopics` function on broker object
```
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
```