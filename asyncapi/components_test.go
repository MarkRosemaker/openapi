package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// allComponents is a document that makes use of every field of the components object.
const allComponents = `{
  "asyncapi": "3.1.0",
  "info": {
    "title": "Account Service",
    "version": "1.0.0",
    "contact": {
      "name": "API Support",
      "url": "https://www.example.com/support",
      "email": "support@example.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "https://www.apache.org/licenses/LICENSE-2.0"
    },
    "externalDocs": {
      "$ref": "#/components/externalDocs/infoDocs"
    }
  },
  "servers": {
    "production": {
      "$ref": "#/components/servers/production"
    }
  },
  "channels": {
    "userSignedup": {
      "$ref": "#/components/channels/userSignedup"
    }
  },
  "operations": {
    "sendUserSignedup": {
      "$ref": "#/components/operations/sendUserSignedup"
    }
  },
  "components": {
    "schemas": {
      "userSignedUpPayload": {
        "type": "object",
        "properties": {
          "displayName": {
            "type": "string"
          }
        }
      }
    },
    "servers": {
      "production": {
        "host": "rabbitmq.example.org",
        "protocol": "amqp",
        "variables": {
          "port": {
            "$ref": "#/components/serverVariables/port"
          }
        },
        "security": [
          {
            "$ref": "#/components/securitySchemes/user-password"
          }
        ],
        "tags": [
          {
            "$ref": "#/components/tags/user"
          }
        ],
        "bindings": {
          "$ref": "#/components/serverBindings/amqp"
        }
      }
    },
    "channels": {
      "userSignedup": {
        "address": "user/{userId}/signedup",
        "messages": {
          "userSignedUp": {
            "$ref": "#/components/messages/userSignedUp"
          }
        },
        "parameters": {
          "userId": {
            "$ref": "#/components/parameters/userId"
          }
        },
        "bindings": {
          "$ref": "#/components/channelBindings/amqp"
        }
      },
      "userSignedupReply": {
        "messages": {
          "userSignedUp": {
            "$ref": "#/components/messages/userSignedUp"
          }
        }
      }
    },
    "operations": {
      "sendUserSignedup": {
        "action": "send",
        "channel": {
          "$ref": "#/components/channels/userSignedup"
        },
        "bindings": {
          "$ref": "#/components/operationBindings/amqp"
        },
        "traits": [
          {
            "$ref": "#/components/operationTraits/kafka"
          }
        ],
        "messages": [
          {
            "$ref": "#/components/channels/userSignedup/messages/userSignedUp"
          }
        ],
        "reply": {
          "$ref": "#/components/replies/userSignedupReply"
        }
      }
    },
    "messages": {
      "userSignedUp": {
        "payload": {
          "$ref": "#/components/schemas/userSignedUpPayload"
        },
        "correlationId": {
          "$ref": "#/components/correlationIds/default"
        },
        "bindings": {
          "$ref": "#/components/messageBindings/amqp"
        },
        "traits": [
          {
            "$ref": "#/components/messageTraits/commonHeaders"
          }
        ]
      }
    },
    "securitySchemes": {
      "user-password": {
        "type": "userPassword"
      }
    },
    "serverVariables": {
      "port": {
        "enum": [
          "5672",
          "5673"
        ],
        "default": "5672"
      }
    },
    "parameters": {
      "userId": {
        "description": "Id of the user."
      }
    },
    "correlationIds": {
      "default": {
        "description": "Default Correlation ID",
        "location": "$message.header#/correlationId"
      }
    },
    "replies": {
      "userSignedupReply": {
        "address": {
          "$ref": "#/components/replyAddresses/userSignedupReply"
        },
        "channel": {
          "$ref": "#/components/channels/userSignedupReply"
        },
        "messages": [
          {
            "$ref": "#/components/channels/userSignedupReply/messages/userSignedUp"
          }
        ]
      }
    },
    "replyAddresses": {
      "userSignedupReply": {
        "description": "Consumer inbox",
        "location": "$message.header#/replyTo"
      }
    },
    "externalDocs": {
      "infoDocs": {
        "url": "https://example.com/docs"
      }
    },
    "tags": {
      "user": {
        "name": "user",
        "description": "User-related messages"
      }
    },
    "operationTraits": {
      "kafka": {
        "bindings": {
          "kafka": {
            "clientId": "my-app-id"
          }
        }
      }
    },
    "messageTraits": {
      "commonHeaders": {
        "headers": {
          "type": "object",
          "properties": {
            "my-app-header": {
              "type": "integer"
            }
          }
        }
      }
    },
    "serverBindings": {
      "amqp": {
        "amqp": {
          "vhost": "/"
        }
      }
    },
    "channelBindings": {
      "amqp": {
        "amqp": {
          "is": "queue"
        }
      }
    },
    "operationBindings": {
      "amqp": {
        "amqp": {
          "ack": true
        }
      }
    },
    "messageBindings": {
      "amqp": {
        "amqp": {
          "contentEncoding": "gzip"
        }
      }
    }
  }
}`

func TestComponents(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromDataJSON([]byte(allComponents))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	// every reference was resolved
	op := doc.Operations["sendUserSignedup"]
	if op.Value != doc.Components.Operations["sendUserSignedup"].Value {
		t.Fatal("the operation was not resolved")
	}

	reply := op.Value.Reply
	if reply.Value != doc.Components.Replies["userSignedupReply"].Value {
		t.Fatal("the reply was not resolved")
	}

	if reply.Value.Address.Value != doc.Components.ReplyAddresses["userSignedupReply"].Value {
		t.Fatal("the reply address was not resolved")
	}

	if got, want := reply.Value.Address.Value.Location,
		asyncapi.RuntimeExpression("$message.header#/replyTo"); got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	msg := doc.Components.Messages["userSignedUp"].Value
	if msg.CorrelationID.Value != doc.Components.CorrelationIDs["default"].Value {
		t.Fatal("the correlation ID was not resolved")
	}

	server := doc.Servers["production"].Value
	if server.Bindings.Value != doc.Components.ServerBindings["amqp"].Value {
		t.Fatal("the server bindings were not resolved")
	}

	if server.Tags[0].Value != doc.Components.Tags["user"].Value {
		t.Fatal("the tag was not resolved")
	}

	if doc.Info.ExternalDocs.Value != doc.Components.ExternalDocs["infoDocs"].Value {
		t.Fatal("the external documentation was not resolved")
	}

	// the document is written back as it was read
	got, err := doc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != allComponents {
		t.Fatalf("got:\n%s\nwant:\n%s", got, allComponents)
	}
}

func TestComponents_ReplyWithAddressAndChannelAddress(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromDataJSON([]byte(allComponents))
	if err != nil {
		t.Fatal(err)
	}

	// when the reply defines an address, the channel it refers to must not have one
	doc.Components.Channels["userSignedupReply"].Value.Address = "user/signedup/reply"

	err = doc.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	want := `components.replies["userSignedupReply"].channel.address ("user/signedup/reply") ` +
		`is invalid: must be empty when the reply defines an address`
	if err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}
}
