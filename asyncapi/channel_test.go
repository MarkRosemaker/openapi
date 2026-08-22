package asyncapi_test

import (
	"slices"
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestChannel_AddressExpressions(t *testing.T) {
	t.Parallel()

	for address, want := range map[string][]string{
		"user/signedup":                     {},
		"user/{userId}/signedup":            {"userId"},
		"{env}/user/{userId}/{action}":      {"env", "userId", "action"},
		"smartylighting.{streetlightId}.on": {"streetlightId"},
	} {
		t.Run(address, func(t *testing.T) {
			t.Parallel()

			c := &asyncapi.Channel{Address: address}
			if got := c.AddressExpressions(); !slices.Equal(got, want) {
				t.Fatalf("got: %v, want: %v", got, want)
			}
		})
	}
}

func TestChannel_Validate_MissingParameter(t *testing.T) {
	t.Parallel()

	doc := minimalDocument()
	doc.Channels["userSignedup"].Value.Address = "user/{userId}/signedup"

	err := doc.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	want := `channels["userSignedup"].parameters["userId"] is required`
	if err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}

	// the parameter is defined, so the channel is valid
	doc.Channels["userSignedup"].Value.Parameters = asyncapi.Parameters{
		"userId": {Value: &asyncapi.Parameter{Description: "Id of the user."}},
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestParameter_Validate(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		param *asyncapi.Parameter
		want  string
	}{
		"empty enum": {
			&asyncapi.Parameter{Enum: []string{}},
			"enum array must not be empty",
		},
		"default not in enum": {
			&asyncapi.Parameter{Enum: []string{"a", "b"}, Default: "c"},
			`default value "c" must exist in the enum's values`,
		},
		"invalid location": {
			&asyncapi.Parameter{Location: "$request.header#/foo"},
			`location ("$request.header#/foo") is invalid: ` +
				`must be a runtime expression, e.g. "$message.header#/correlationId"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.param.Validate()
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("got: %v, want: %v", err, tc.want)
			}
		})
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		p := &asyncapi.Parameter{
			Enum:     []string{"a", "b"},
			Default:  "a",
			Examples: []string{"a"},
			Location: "$message.payload#/user/id",
		}

		if err := p.Validate(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestOperation_Validate_Errors(t *testing.T) {
	t.Parallel()

	channelRef := func() *asyncapi.ChannelRef {
		return &asyncapi.ChannelRef{
			Ref:   &asyncapi.Reference{Identifier: "#/channels/userSignedup"},
			Value: &asyncapi.Channel{Address: "user/signedup"},
		}
	}

	for name, tc := range map[string]struct {
		op   *asyncapi.Operation
		want string
	}{
		"no action": {
			&asyncapi.Operation{Channel: channelRef()},
			`operations["test"].action is required`,
		},
		"unknown action": {
			&asyncapi.Operation{Action: "publish", Channel: channelRef()},
			`operations["test"].action ("publish") is invalid, must be one of: "send", "receive"`,
		},
		"no channel": {
			&asyncapi.Operation{Action: asyncapi.OperationActionSend},
			`operations["test"].channel is required`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			doc := minimalDocument()
			doc.Operations = asyncapi.Operations{"test": {Value: tc.op}}

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
