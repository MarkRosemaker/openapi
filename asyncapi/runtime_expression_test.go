package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestRuntimeExpression_Validate(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		for _, expr := range []asyncapi.RuntimeExpression{
			"$message.header#/MQMD/CorrelId",
			"$message.payload#/messageId",
			"$message.header",
			"$message.payload",
		} {
			if err := expr.Validate(); err != nil {
				t.Fatalf("%s: %v", expr, err)
			}
		}
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		var expr asyncapi.RuntimeExpression

		err := expr.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "a value is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		for _, expr := range []asyncapi.RuntimeExpression{
			"$request.header#/foo",
			"$message.body#/foo",
			"message.header",
			"$message.header/foo",
		} {
			if err := expr.Validate(); err == nil {
				t.Fatalf("%s: expected error", expr)
			}
		}
	})
}

func TestCorrelationID_Validate(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/rpc-client.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	msg := doc.Channels["queue"].Value.Messages["receiveSumResult"].Value
	if got, want := msg.CorrelationID.Value.Location,
		asyncapi.RuntimeExpression("$message.header#/correlation_id"); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	// an invalid location is reported
	msg.CorrelationID.Value.Location = "$message.somewhere"

	err = doc.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	want := `channels["queue"].messages["receiveSumResult"].correlationId.location ` +
		`("$message.somewhere") is invalid: ` +
		`must be a runtime expression, e.g. "$message.header#/correlationId"`
	if err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}
}
