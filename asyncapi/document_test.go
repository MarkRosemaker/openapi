package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// minimalDocument returns the smallest document that is valid.
func minimalDocument() *asyncapi.Document {
	return &asyncapi.Document{
		AsyncAPI: "3.1.0",
		Info:     &asyncapi.Info{Title: "Account Service", Version: "1.0.0"},
		Channels: asyncapi.Channels{
			"userSignedup": {Value: &asyncapi.Channel{Address: "user/signedup"}},
		},
	}
}

func TestDocument_Validate(t *testing.T) {
	t.Parallel()

	if err := minimalDocument().Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDocument_Validate_Errors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		doc  func(*asyncapi.Document)
		want string
	}{
		"no version": {
			func(d *asyncapi.Document) { d.AsyncAPI = "" },
			"asyncapi is required",
		},
		"version of another major version": {
			func(d *asyncapi.Document) { d.AsyncAPI = "2.6.0" },
			`asyncapi ("2.6.0") is invalid: must be a valid version (3.x.y)`,
		},
		"version is not a version": {
			func(d *asyncapi.Document) { d.AsyncAPI = "three" },
			`asyncapi ("three") is invalid: must be a valid version (3.x.y)`,
		},
		"no info": {
			func(d *asyncapi.Document) { d.Info = nil },
			"info is required",
		},
		"no title": {
			func(d *asyncapi.Document) { d.Info.Title = "" },
			"info.title is required",
		},
		"no version of the API": {
			func(d *asyncapi.Document) { d.Info.Version = "" },
			"info.version is required",
		},
		"nothing is described": {
			func(d *asyncapi.Document) { d.Channels = nil },
			asyncapi.ErrEmptyDocument.Error(),
		},
		"invalid server name": {
			func(d *asyncapi.Document) {
				d.Servers = asyncapi.Servers{"my server": {Value: &asyncapi.Server{
					Host: "example.com", Protocol: asyncapi.ProtocolKafka,
				}}}
			},
			`servers["my server"] ("my server") is invalid: must match the regular expression "^[A-Za-z0-9_\-]+$"`,
		},
		"server without a host": {
			func(d *asyncapi.Document) {
				d.Servers = asyncapi.Servers{"production": {Value: &asyncapi.Server{
					Protocol: asyncapi.ProtocolKafka,
				}}}
			},
			`servers["production"].host is required`,
		},
		"server without a protocol": {
			func(d *asyncapi.Document) {
				d.Servers = asyncapi.Servers{"production": {Value: &asyncapi.Server{
					Host: "example.com",
				}}}
			},
			`servers["production"].protocol is required`,
		},
		"invalid default content type": {
			func(d *asyncapi.Document) { d.DefaultContentType = "not a media type" },
			"defaultContentType: mime: expected slash after first token",
		},
		"invalid component name": {
			func(d *asyncapi.Document) {
				d.Components.Schemas = asyncapi.Schemas{
					"not a valid name": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{}}},
				}
			},
			`components.schemas["not a valid name"] ("not a valid name") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
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

func TestDocument_SortMaps(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/streetlights-kafka.yaml")
	if err != nil {
		t.Fatal(err)
	}

	doc.SortMaps()

	want := []string{"lightTurnOff", "lightTurnOn", "lightingMeasured", "lightsDim"}

	i := 0
	for name := range doc.Channels.ByIndex() {
		if name != want[i] {
			t.Fatalf("got: %v, want: %v", name, want[i])
		}

		i++
	}

	wantSchemas := []string{
		"dimLightPayload", "lightMeasuredPayload", "sentAt", "turnOnOffPayload",
	}

	i = 0
	for name := range doc.Components.Schemas.ByIndex() {
		if name != wantSchemas[i] {
			t.Fatalf("got: %v, want: %v", name, wantSchemas[i])
		}

		i++
	}
}
