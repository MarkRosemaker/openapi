package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestReferences_Resolve(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/simple.yaml")
	if err != nil {
		t.Fatal(err)
	}

	channel := doc.Channels["userSignedup"]
	if channel.Value == nil {
		t.Fatal("the channel was not resolved")
	}

	// the message of the channel is a reference to a message of the components object
	msg := channel.Value.Messages["UserSignedUp"]
	if msg.Ref == nil {
		t.Fatal("expected the message to be given as a reference")
	}

	if msg.Value != doc.Components.Messages["UserSignedUp"].Value {
		t.Fatal("the message of the channel was not resolved")
	}

	// the operation refers to the message of the channel,
	// which in turn refers to the message of the components object
	op := doc.Operations["sendUserSignedup"]
	if got, want := len(op.Value.Messages), 1; got != want {
		t.Fatalf("got: %d messages, want: %d", got, want)
	}

	if op.Value.Messages[0].Value != msg.Value {
		t.Fatal("the message of the operation was not resolved")
	}

	// the operation refers to the channel
	if op.Value.Channel.Value != channel.Value {
		t.Fatal("the channel of the operation was not resolved")
	}
}

func TestReferences_MustBeReference(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		doc  func(*asyncapi.Document)
		want string
	}{
		"channel of an operation": {
			func(d *asyncapi.Document) {
				d.Operations = asyncapi.Operations{"sendUserSignedup": {Value: &asyncapi.Operation{
					Action:  asyncapi.OperationActionSend,
					Channel: &asyncapi.ChannelRef{Value: &asyncapi.Channel{Address: "user/signedup"}},
				}}}
			},
			`operations["sendUserSignedup"].channel: ` + asyncapi.ErrMustBeReference.Error(),
		},
		"messages of an operation": {
			func(d *asyncapi.Document) {
				d.Operations = asyncapi.Operations{"sendUserSignedup": {Value: &asyncapi.Operation{
					Action: asyncapi.OperationActionSend,
					Channel: &asyncapi.ChannelRef{
						Ref: &asyncapi.Reference{Identifier: "#/channels/userSignedup"},
						// the reference was resolved when the document was loaded
						Value: &asyncapi.Channel{Address: "user/signedup"},
					},
					Messages: asyncapi.MessageRefList{{Value: &asyncapi.Message{Name: "userSignedUp"}}},
				}}}
			},
			`operations["sendUserSignedup"].messages[0]: ` + asyncapi.ErrMustBeReference.Error(),
		},
		"servers of a channel": {
			func(d *asyncapi.Document) {
				d.Channels["userSignedup"].Value.Servers = asyncapi.ServerRefList{
					{Value: &asyncapi.Server{Host: "example.com", Protocol: asyncapi.ProtocolKafka}},
				}
			},
			`channels["userSignedup"].servers[0]: ` + asyncapi.ErrMustBeReference.Error(),
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			doc := minimalDocument()
			tc.doc(doc)

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

func TestReferences_Unresolved(t *testing.T) {
	t.Parallel()

	// a reference that was never resolved is reported when validating
	doc := minimalDocument()
	doc.Channels["userSignedup"].Value.Messages = asyncapi.Messages{
		"UserSignedUp": {Ref: &asyncapi.Reference{Identifier: "#/components/messages/UserSignedUp"}},
	}

	err := doc.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	want := `channels["userSignedup"].messages["UserSignedUp"]: ` +
		`#/components/messages/UserSignedUp (*asyncapi.Message) was not resolved`
	if err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}
}

func TestReference_Validate(t *testing.T) {
	t.Parallel()

	r := &asyncapi.Reference{}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := "$ref is required"; err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}

	r.Identifier = "#/components/messages/UserSignedUp"
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}
