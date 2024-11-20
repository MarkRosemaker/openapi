package openapi

// Webhooks describes requests initiated other than by an API call, for example by an out of band registration.
// The key name is a unique string to refer to each webhook, while the (optionally referenced) Path Item Object describes a request that may be initiated by the API provider and the expected responses.
type Webhooks map[string]*PathItemRef

// Validate checks the Webhooks for correctness.
func (ws Webhooks) Validate() error {
	for name, w := range ws {
		if err := w.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveWebhooks(ws Webhooks) error {
	for name, w := range ws {
		if err := l.resolvePathItemRef(w); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
