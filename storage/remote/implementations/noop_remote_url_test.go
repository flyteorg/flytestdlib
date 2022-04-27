package implementations

import (
	"context"
	"testing"

	"github.com/flyteorg/flyteadmin/pkg/common/mocks"
	"github.com/flyteorg/flytestdlib/storage"

	"github.com/stretchr/testify/assert"
)

const noopFileSize = int64(1256)

type MockMetadata struct{}

func (m MockMetadata) Exists() bool {
	return true
}

func (m MockMetadata) Size() int64 {
	return noopFileSize
}

func getMockStorage() storage.DataStore {
	mockStorage := mocks.GetMockStorageClient()
	mockStorage.ComposedProtobufStore.(*mocks.TestDataStore).HeadCb =
		func(ctx context.Context, reference storage.DataReference) (storage.Metadata, error) {
			return MockMetadata{}, nil
		}
	return *mockStorage
}

func TestNoopRemoteURLGet(t *testing.T) {
	noopRemoteURL := NewNoopRemoteURL(getMockStorage())
	urlBlob, err := noopRemoteURL.Get(context.Background(), "uri")
	assert.Nil(t, err)
	assert.NotEmpty(t, urlBlob)
	assert.Equal(t, "uri", urlBlob.Url)
	assert.Equal(t, noopFileSize, urlBlob.Bytes)
}
