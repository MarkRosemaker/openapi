package asyncapi_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// TestValidate_Strict checks the rules that the specification states as a MUST
// but that a document can break without being syntactically wrong.
func TestValidate_Strict(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		doc  func(*asyncapi.Document)
		want string
	}{
		"identifier that is not a URI": {
			func(d *asyncapi.Document) { d.ID = mustParseURL("smartylighting") },
			`id ("smartylighting") is invalid: must conform to the URI format`,
		},
		"terms of service that is not an absolute URL": {
			func(d *asyncapi.Document) { d.Info.TermsOfService = &url.URL{Path: "terms"} },
			`info.termsOfService ("https://terms") is invalid: must be an absolute URL`,
		},
		"contact URL that is not an absolute URL": {
			func(d *asyncapi.Document) {
				d.Info.Contact = &asyncapi.Contact{URL: &url.URL{Path: "support"}}
			},
			`info.contact.url ("https://support") is invalid: must be an absolute URL`,
		},
		"license URL that is not an absolute URL": {
			func(d *asyncapi.Document) {
				d.Info.License = &asyncapi.License{Name: "Apache 2.0", URL: &url.URL{Path: "license"}}
			},
			`info.license.url ("https://license") is invalid: must be an absolute URL`,
		},
		"external docs URL that is not an absolute URL": {
			func(d *asyncapi.Document) {
				d.Info.ExternalDocs = &asyncapi.ExternalDocsRef{
					Value: &asyncapi.ExternalDocs{URL: &url.URL{Path: "docs"}},
				}
			},
			`info.externalDocs.url ("https://docs") is invalid: must be an absolute URL`,
		},
		"channel address with a query": {
			func(d *asyncapi.Document) {
				d.Channels["userSignedup"].Value.Address = "user/signedup?filter=all"
			},
			`channels["userSignedup"].address ("user/signedup?filter=all") is invalid: ` +
				"query parameters and fragments must not be used, use bindings instead",
		},
		"channel address with a fragment": {
			func(d *asyncapi.Document) {
				d.Channels["userSignedup"].Value.Address = "user/signedup#now"
			},
			`channels["userSignedup"].address ("user/signedup#now") is invalid: ` +
				"query parameters and fragments must not be used, use bindings instead",
		},
		"parameter name that doesn't match the pattern": {
			func(d *asyncapi.Document) {
				d.Channels["userSignedup"].Value.Parameters = asyncapi.Parameters{
					"user id": {Value: &asyncapi.Parameter{}},
				}
			},
			`channels["userSignedup"].parameters["user id"] ("user id") is invalid: ` +
				`must match the regular expression "^[A-Za-z0-9_\-]+$"`,
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

func TestValidate_Locations(t *testing.T) {
	t.Parallel()

	t.Run("server that is not in the root", func(t *testing.T) {
		t.Parallel()

		// the channel refers to a server of the components object
		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"channels":{"userSignedup":{"servers":[` +
				`{"$ref":"#/components/servers/production"}]}},` +
				`"components":{"servers":{"production":` +
				`{"host":"example.com","protocol":"kafka"}}}}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		validateErr := doc.Validate()
		if validateErr == nil {
			t.Fatal("expected error")
		}

		want := `channels["userSignedup"].servers[0]: ` + asyncapi.ErrServerNotInRoot.Error()
		if validateErr.Error() != want {
			t.Fatalf("got: %v, want: %v", validateErr, want)
		}
	})

	t.Run("channel that is not in the root", func(t *testing.T) {
		t.Parallel()

		// the operation refers to a channel of the components object
		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"operations":{"sendUserSignedup":{"action":"send",` +
				`"channel":{"$ref":"#/components/channels/userSignedup"}}},` +
				`"components":{"channels":{"userSignedup":{"address":"user/signedup"}}}}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		validateErr := doc.Validate()
		if validateErr == nil {
			t.Fatal("expected error")
		}

		want := `operations["sendUserSignedup"].channel: ` + asyncapi.ErrChannelNotInRoot.Error()
		if validateErr.Error() != want {
			t.Fatalf("got: %v, want: %v", validateErr, want)
		}
	})

	t.Run("message that is not of the channel", func(t *testing.T) {
		t.Parallel()

		// the operation refers to a message of another channel
		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"channels":{"userSignedup":{"address":"user/signedup"},` +
				`"userLoggedIn":{"address":"user/loggedin","messages":{"userLoggedIn":{}}}},` +
				`"operations":{"sendUserSignedup":{"action":"send",` +
				`"channel":{"$ref":"#/channels/userSignedup"},` +
				`"messages":[{"$ref":"#/channels/userLoggedIn/messages/userLoggedIn"}]}}}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		validateErr := doc.Validate()
		if validateErr == nil {
			t.Fatal("expected error")
		}

		want := `operations["sendUserSignedup"].messages[0]: ` +
			asyncapi.ErrMessageNotOfChannel.Error()
		if validateErr.Error() != want {
			t.Fatalf("got: %v, want: %v", validateErr, want)
		}
	})
}

func TestReference_IgnoresAdditionalProperties(t *testing.T) {
	t.Parallel()

	// "This object cannot be extended with additional properties
	// and any properties added SHALL be ignored."
	doc, err := asyncapi.LoadFromDataJSON([]byte(
		`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
			`"channels":{"userSignedup":{"address":"user/signedup","messages":{"userSignedUp":` +
			`{"$ref":"#/components/messages/userSignedUp","summary":"ignored"}}}},` +
			`"components":{"messages":{"userSignedUp":{"name":"userSignedUp"}}}}`,
	))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	msg := doc.Channels["userSignedup"].Value.Messages["userSignedUp"]
	if msg.Ref == nil {
		t.Fatal("expected the message to be given as a reference")
	}

	if msg.Value != doc.Components.Messages["userSignedUp"].Value {
		t.Fatal("the message was not resolved")
	}

	// the additional property is not written back
	got, err := doc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	if want := "ignored"; strings.Contains(string(got), want) {
		t.Fatalf("got:\n%s\nwant it to not contain: %s", got, want)
	}
}
