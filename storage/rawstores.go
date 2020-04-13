package storage

import (
	"fmt"
	"net/http"

	"github.com/lyft/flytestdlib/promutils"
)

type dataStoreCreateFn func(cfg *Config, metricsScope promutils.Scope) (RawStore, error)

var stores = map[string]dataStoreCreateFn{
	TypeMemory: NewInMemoryRawStore,
	TypeLocal:  newLocalRawStore,
	TypeMinio:  newStowRawStore,
	TypeS3:     newStowRawStore,
	TypeStow:   newStowRawStore,
}

type proxyTransport struct {
	http.RoundTripper
	defaultHeaders map[string][]string
}

func (p proxyTransport) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	for key, values := range p.defaultHeaders {
		for _, val := range values {
			r.Header.Add(key, val)
		}
	}

	return p.RoundTripper.RoundTrip(r)
}

func createHttpClientWithDefaultHeaders(headers map[string][]string) *http.Client {
	c := &http.Client{}
	c.Transport = &proxyTransport{
		RoundTripper:   http.DefaultTransport,
		defaultHeaders: headers,
	}

	return c
}

// Creates a new Data Store with the supplied config.
func NewDataStore(cfg *Config, metricsScope promutils.Scope) (s *DataStore, err error) {
	// HACK: This sets http headers to the default http client. This is because
	// some underlying stores (e.g. S3 Stow Store) grabs the default http client
	// and doesn't allow configuration of default headers.
	if len(cfg.DefaultHttpClientHeaders) > 0 {
		defaultClient := http.DefaultClient
		defer func() {
			http.DefaultClient = defaultClient
		}()

		http.DefaultClient = createHttpClientWithDefaultHeaders(cfg.DefaultHttpClientHeaders)
	}

	var rawStore RawStore
	if fn, found := stores[cfg.Type]; found {
		rawStore, err = fn(cfg, metricsScope)
		if err != nil {
			return &emptyStore, err
		}

		protoStore := NewDefaultProtobufStore(newCachedRawStore(cfg, rawStore, metricsScope), metricsScope)
		return NewCompositeDataStore(URLPathConstructor{}, protoStore), nil
	}

	return &emptyStore, fmt.Errorf("type is of an invalid value [%v]", cfg.Type)
}

// Composes a new DataStore.
func NewCompositeDataStore(refConstructor ReferenceConstructor, composedProtobufStore ComposedProtobufStore) *DataStore {
	return &DataStore{
		ReferenceConstructor:  refConstructor,
		ComposedProtobufStore: composedProtobufStore,
	}
}
