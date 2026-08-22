package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// fullDocument is a document that makes use of as many fields as possible.
const fullDocument = `{
  "asyncapi": "3.1.0",
  "id": "urn:example:com:smartylighting:streetlights:server",
  "info": {
    "title": "Account Service",
    "version": "1.0.0",
    "description": "This service is in charge of processing user signups.",
    "termsOfService": "https://example.com/terms",
    "contact": {
      "name": "API Support",
      "url": "https://www.example.com/support",
      "email": "support@example.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "https://www.apache.org/licenses/LICENSE-2.0"
    },
    "tags": [
      {
        "name": "user",
        "description": "User-related messages",
        "externalDocs": {
          "description": "Find more info here",
          "url": "https://example.com/docs/user"
        }
      }
    ],
    "externalDocs": {
      "url": "https://example.com/docs"
    }
  },
  "servers": {
    "production": {
      "host": "{env}.example.org:{port}",
      "protocol": "kafka-secure",
      "protocolVersion": "3.5.0",
      "pathname": "/events",
      "description": "The production broker.",
      "title": "Production",
      "summary": "Production broker",
      "variables": {
        "env": {
          "enum": [
            "prod",
            "staging"
          ],
          "default": "prod",
          "description": "The environment.",
          "examples": [
            "prod"
          ]
        },
        "port": {
          "default": "9092"
        }
      },
      "security": [
        {
          "type": "oauth2",
          "description": "Sign in with your account.",
          "flows": {
            "implicit": {
              "authorizationUrl": "https://example.com/oauth/authorize",
              "refreshUrl": "https://example.com/oauth/refresh",
              "availableScopes": {
                "user:read": "Read user data"
              }
            },
            "password": {
              "tokenUrl": "https://example.com/oauth/token",
              "availableScopes": {}
            },
            "clientCredentials": {
              "tokenUrl": "https://example.com/oauth/token",
              "availableScopes": {}
            },
            "authorizationCode": {
              "authorizationUrl": "https://example.com/oauth/authorize",
              "tokenUrl": "https://example.com/oauth/token",
              "refreshUrl": "https://example.com/oauth/refresh",
              "availableScopes": {}
            }
          },
          "scopes": [
            "user:read"
          ]
        }
      ],
      "tags": [
        {
          "name": "env:production"
        }
      ],
      "externalDocs": {
        "url": "https://example.com/docs/production"
      },
      "bindings": {
        "kafka": {
          "schemaRegistryUrl": "https://schema-registry.example.com/"
        }
      }
    }
  },
  "defaultContentType": "application/json",
  "channels": {
    "userSignedup": {
      "address": "user/{userId}/signedup",
      "messages": {
        "userSignedUp": {
          "headers": {
            "type": "object",
            "properties": {
              "correlationId": {
                "type": "string"
              }
            }
          },
          "payload": {
            "$ref": "#/components/schemas/userSignedUpPayload"
          },
          "correlationId": {
            "location": "$message.header#/correlationId"
          },
          "contentType": "application/json",
          "name": "userSignedUp",
          "title": "User signed up",
          "summary": "A user signed up.",
          "description": "This message is sent when a user signs up.",
          "tags": [
            {
              "name": "signup"
            }
          ],
          "externalDocs": {
            "url": "https://example.com/docs/user-signed-up"
          },
          "bindings": {
            "kafka": {
              "key": {
                "type": "string"
              }
            }
          },
          "examples": [
            {
              "headers": {
                "correlationId": "my-correlation-id"
              },
              "payload": {
                "displayName": "Lucas"
              },
              "name": "SimpleSignup",
              "summary": "A simple example"
            }
          ],
          "traits": [
            {
              "$ref": "#/components/messageTraits/commonHeaders"
            }
          ]
        }
      },
      "title": "User signed up",
      "summary": "The channel a user signup is announced on.",
      "description": "This channel is used to announce user signups.",
      "servers": [
        {
          "$ref": "#/servers/production"
        }
      ],
      "parameters": {
        "userId": {
          "enum": [
            "1",
            "2"
          ],
          "default": "1",
          "description": "Id of the user.",
          "examples": [
            "1"
          ],
          "location": "$message.payload#/user/id"
        }
      },
      "tags": [
        {
          "name": "user"
        }
      ],
      "externalDocs": {
        "url": "https://example.com/docs/user-signedup"
      },
      "bindings": {
        "kafka": {
          "topic": "user-signedup"
        }
      }
    }
  },
  "operations": {
    "sendUserSignedup": {
      "action": "send",
      "channel": {
        "$ref": "#/channels/userSignedup"
      },
      "title": "Send a user signup",
      "summary": "Announce that a user signed up.",
      "description": "This operation announces that a user signed up.",
      "security": [
        {
          "type": "userPassword"
        }
      ],
      "tags": [
        {
          "name": "user"
        }
      ],
      "externalDocs": {
        "url": "https://example.com/docs/send-user-signedup"
      },
      "bindings": {
        "kafka": {
          "groupId": {
            "type": "string"
          }
        }
      },
      "traits": [
        {
          "title": "Kafka operation",
          "summary": "An operation on Kafka.",
          "description": "This trait is applied to all Kafka operations.",
          "security": [
            {
              "type": "userPassword"
            }
          ],
          "tags": [
            {
              "name": "kafka"
            }
          ],
          "externalDocs": {
            "url": "https://example.com/docs/kafka"
          },
          "bindings": {
            "kafka": {
              "clientId": "my-app-id"
            }
          }
        }
      ],
      "messages": [
        {
          "$ref": "#/channels/userSignedup/messages/userSignedUp"
        }
      ]
    }
  },
  "components": {
    "schemas": {
      "userSignedUpPayload": {
        "$id": "https://example.com/schemas/userSignedUpPayload",
        "$schema": "http://json-schema.org/draft-07/schema",
        "$comment": "The payload of a user signup.",
        "title": "User signup payload",
        "description": "The payload of a user signup.",
        "type": "object",
        "allOf": [
          {
            "type": "object"
          }
        ],
        "oneOf": [
          {
            "required": [
              "displayName"
            ]
          }
        ],
        "anyOf": [
          {
            "required": [
              "email"
            ]
          }
        ],
        "not": {
          "type": "null"
        },
        "if": {
          "required": [
            "email"
          ]
        },
        "then": {
          "required": [
            "displayName"
          ]
        },
        "else": true,
        "minProperties": 1,
        "maxProperties": 10,
        "required": [
          "displayName"
        ],
        "properties": {
          "displayName": {
            "type": "string",
            "minLength": 1,
            "maxLength": 100,
            "pattern": "^[a-zA-Z ]+$",
            "examples": [
              "Lucas"
            ]
          },
          "email": {
            "type": "string",
            "format": "email"
          },
          "age": {
            "type": "integer",
            "multipleOf": 1,
            "minimum": 0,
            "exclusiveMinimum": -1,
            "maximum": 200,
            "exclusiveMaximum": 201,
            "default": 0
          },
          "roles": {
            "type": "array",
            "minItems": 1,
            "maxItems": 10,
            "uniqueItems": true,
            "items": {
              "type": "string",
              "enum": [
                "admin",
                "user"
              ]
            },
            "additionalItems": false,
            "contains": {
              "const": "user"
            }
          },
          "avatar": {
            "type": "string",
            "contentEncoding": "base64",
            "contentMediaType": "image/png"
          },
          "createdAt": {
            "type": "string",
            "format": "date-time",
            "readOnly": true
          },
          "password": {
            "type": "string",
            "format": "password",
            "writeOnly": true
          },
          "legacyId": {
            "type": [
              "string",
              "null"
            ],
            "deprecated": true
          }
        },
        "patternProperties": {
          "^x-": {
            "type": "string"
          }
        },
        "additionalProperties": false,
        "propertyNames": {
          "pattern": "^[a-zA-Z]+$"
        },
        "definitions": {
          "empty": {}
        },
        "externalDocs": {
          "url": "https://example.com/docs/user-signed-up-payload"
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
        },
        "correlationId": {
          "location": "$message.header#/correlationId"
        },
        "contentType": "application/json",
        "name": "commonHeaders",
        "title": "Common headers",
        "summary": "The headers all messages have in common.",
        "description": "These headers are applied to all messages.",
        "tags": [
          {
            "name": "headers"
          }
        ],
        "externalDocs": {
          "url": "https://example.com/docs/common-headers"
        },
        "bindings": {
          "kafka": {
            "bindingVersion": "0.5.0"
          }
        },
        "examples": [
          {
            "headers": {
              "my-app-header": 12
            }
          }
        ]
      }
    }
  }
}`

func TestFullDocument(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromDataJSON([]byte(fullDocument))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	// the document is written back as it was read
	got, err := doc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != fullDocument {
		t.Fatalf("got:\n%s\nwant:\n%s", got, fullDocument)
	}
}
