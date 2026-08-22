package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestMessage_Traits(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/streetlights-kafka.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	msg := doc.Components.Messages["lightMeasured"].Value
	if got, want := len(msg.Traits), 1; got != want {
		t.Fatalf("got: %d traits, want: %d", got, want)
	}

	// the trait was resolved
	if msg.Traits[0].Value != doc.Components.MessageTraits["commonHeaders"].Value {
		t.Fatal("the trait of the message was not resolved")
	}

	if got, want := msg.ContentType, asyncapi.MediaTypeJSON; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	// the payload was resolved
	if msg.Payload.Value != doc.Components.Schemas["lightMeasuredPayload"].Value {
		t.Fatal("the payload of the message was not resolved")
	}
}

func TestMessage_Validate_Errors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		msg  *asyncapi.Message
		want string
	}{
		"invalid content type": {
			&asyncapi.Message{ContentType: "not a media type"},
			`channels["userSignedup"].messages["test"].contentType: ` +
				"mime: expected slash after first token",
		},
		"example without headers and payload": {
			&asyncapi.Message{Examples: asyncapi.MessageExamples{{Name: "empty"}}},
			`channels["userSignedup"].messages["test"].examples[0]: ` +
				asyncapi.ErrEmptyMessageExample.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			doc := minimalDocument()
			doc.Channels["userSignedup"].Value.Messages = asyncapi.Messages{
				"test": {Value: tc.msg},
			}

			err := doc.Validate()
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("got: %v, want: %v", err, tc.want)
			}
		})
	}
}

func TestMessageExample(t *testing.T) {
	t.Parallel()

	ex := &asyncapi.MessageExample{
		Name:    "SimpleSignup",
		Summary: "A simple UserSignup example message",
		Headers: []byte(`{"correlationId":"my-correlation-id"}`),
		Payload: []byte(`{"user":{"someUserKey":"someUserValue"}}`),
	}

	if err := ex.Validate(); err != nil {
		t.Fatal(err)
	}
}
