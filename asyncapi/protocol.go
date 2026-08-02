package asyncapi

// Protocol is the protocol a server supports for connection,
// respectively the name of the protocol a binding applies to.
type Protocol string

const (
	// ProtocolHTTP is the HTTP protocol.
	ProtocolHTTP Protocol = "http"
	// ProtocolWebSockets is the WebSockets protocol.
	ProtocolWebSockets Protocol = "ws"
	// ProtocolKafka is the Kafka protocol.
	ProtocolKafka Protocol = "kafka"
	// ProtocolAnypointMQ is the Anypoint MQ protocol.
	ProtocolAnypointMQ Protocol = "anypointmq"
	// ProtocolAMQP is the AMQP 0-9-1 protocol.
	ProtocolAMQP Protocol = "amqp"
	// ProtocolAMQP1 is the AMQP 1.0 protocol.
	ProtocolAMQP1 Protocol = "amqp1"
	// ProtocolMQTT is the MQTT protocol.
	ProtocolMQTT Protocol = "mqtt"
	// ProtocolMQTT5 is the MQTT 5 protocol.
	ProtocolMQTT5 Protocol = "mqtt5"
	// ProtocolNATS is the NATS protocol.
	ProtocolNATS Protocol = "nats"
	// ProtocolJMS is the JMS protocol.
	ProtocolJMS Protocol = "jms"
	// ProtocolSNS is the SNS protocol.
	ProtocolSNS Protocol = "sns"
	// ProtocolSolace is the Solace protocol.
	ProtocolSolace Protocol = "solace"
	// ProtocolSQS is the SQS protocol.
	ProtocolSQS Protocol = "sqs"
	// ProtocolSTOMP is the STOMP protocol.
	ProtocolSTOMP Protocol = "stomp"
	// ProtocolRedis is the Redis protocol.
	ProtocolRedis Protocol = "redis"
	// ProtocolMercure is the Mercure protocol.
	ProtocolMercure Protocol = "mercure"
	// ProtocolIBMMQ is the IBM MQ protocol.
	ProtocolIBMMQ Protocol = "ibmmq"
	// ProtocolGooglePubSub is the Google Cloud Pub/Sub protocol.
	ProtocolGooglePubSub Protocol = "googlepubsub"
	// ProtocolPulsar is the Pulsar protocol.
	ProtocolPulsar Protocol = "pulsar"
	// ProtocolROS2 is the ROS 2 protocol.
	ProtocolROS2 Protocol = "ros2"
)

// allBindingProtocols are the protocols a bindings object may describe.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#serverBindingsObject
var allBindingProtocols = []Protocol{
	ProtocolHTTP,
	ProtocolWebSockets,
	ProtocolKafka,
	ProtocolAnypointMQ,
	ProtocolAMQP,
	ProtocolAMQP1,
	ProtocolMQTT,
	ProtocolMQTT5,
	ProtocolNATS,
	ProtocolJMS,
	ProtocolSNS,
	ProtocolSolace,
	ProtocolSQS,
	ProtocolSTOMP,
	ProtocolRedis,
	ProtocolMercure,
	ProtocolIBMMQ,
	ProtocolGooglePubSub,
	ProtocolPulsar,
	ProtocolROS2,
}
