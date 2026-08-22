package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestBindings(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/rpc-client.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	bindings := doc.Channels["queue"].Value.Bindings
	if bindings == nil {
		t.Fatal("expected bindings")
	}

	amqp := (*bindings.Value)[asyncapi.ProtocolAMQP]
	if amqp == nil {
		t.Fatal("expected an AMQP binding")
	}

	// the protocol-specific definitions are kept as they are
	want := `{"is":"queue","queue":{"exclusive":true}}`
	if got := string(amqp.Value); got != want {
		t.Fatalf("got: %s, want: %s", got, want)
	}
}

func TestBindings_Validate(t *testing.T) {
	t.Parallel()

	t.Run("unknown protocol", func(t *testing.T) {
		t.Parallel()

		doc := minimalDocument()
		doc.Channels["userSignedup"].Value.Bindings = &asyncapi.BindingsRef{
			Value: &asyncapi.Bindings{"carrierPigeon": {Value: []byte(`{}`)}},
		}

		err := doc.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		want := `channels["userSignedup"].bindings["carrierPigeon"] ("carrierPigeon") is invalid, ` +
			`must be one of: "http", "ws", "kafka", "anypointmq", "amqp", "amqp1", "mqtt", "mqtt5", ` +
			`"nats", "jms", "sns", "solace", "sqs", "stomp", "redis", "mercure", "ibmmq", ` +
			`"googlepubsub", "pulsar", "ros2"`
		if err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("extension", func(t *testing.T) {
		t.Parallel()

		// the bindings object may be extended with specification extensions
		doc := minimalDocument()
		doc.Channels["userSignedup"].Value.Bindings = &asyncapi.BindingsRef{
			Value: &asyncapi.Bindings{"x-custom": {Value: []byte(`{}`)}},
		}

		if err := doc.Validate(); err != nil {
			t.Fatal(err)
		}
	})
}
