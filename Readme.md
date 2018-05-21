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