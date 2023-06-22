package storage

import (
	"buf.build/gen/go/unionai-oss/unionidl/grpc/go/objectstore/objectstorev1grpc"
	objectstorev1 "buf.build/gen/go/unionai-oss/unionidl/protocolbuffers/go/objectstore"
	"context"
	"fmt"
	"github.com/flyteorg/flytestdlib/ioutils"
	v1 "github.com/unionai-oss/unionidl/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net/url"
)

type UnionMetaStore struct {
	copyImpl
	client    objectstorev1grpc.ObjectStoreServiceClient
	publicUrl string
}

func (u UnionMetaStore) GetBaseContainerFQN(_ context.Context) DataReference {
	return DataReference("unionmeta://")
}

func (u UnionMetaStore) CreateSignedURL(ctx context.Context, reference DataReference, properties SignedURLProperties) (SignedURLResponse, error) {
	if len(u.publicUrl) == 0 {
		return SignedURLResponse{}, fmt.Errorf("public url is not configured")
	}

	signedURL, err := url.Parse(fmt.Sprintf(u.publicUrl, reference.String()))
	if err != nil {
		return SignedURLResponse{}, fmt.Errorf("failed to parse public url. Error: %w", err)
	}

	return SignedURLResponse{
		URL: *signedURL,
	}, nil
}

func (u UnionMetaStore) Head(ctx context.Context, reference DataReference) (Metadata, error) {
	resp, err := u.client.Head(ctx, &objectstorev1.HeadRequest{
		Key: reference.String(),
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return MemoryMetadata{
				exists: false,
			}, nil
		}

		return nil, fmt.Errorf("failed to head object. Error: %w", err)
	}

	return MemoryMetadata{
		exists: true,
		size:   int64(resp.SizeBytes),
		etag:   resp.Etag,
	}, nil
}

func (u UnionMetaStore) ReadRaw(ctx context.Context, reference DataReference) (io.ReadCloser, error) {
	resp, err := u.client.Get(ctx, &objectstorev1.GetRequest{
		Key: reference.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object. Error: %w", err)
	}

	return ioutils.NewBytesReadCloser(resp.Object.Contents), nil
}

func (u UnionMetaStore) WriteRaw(ctx context.Context, reference DataReference, size int64, opts Options, raw io.Reader) error {
	if size <= 0 {
		return fmt.Errorf("invalid size: %d", size)
	}

	m := &objectstorev1.Metadata{
		Tag: make(map[string]string, len(opts.Metadata)),
	}
	for k, v := range opts.Metadata {
		m.Tag[k] = fmt.Sprintf("%v", v)
	}

	data, err := io.ReadAll(raw)
	if err != nil {
		return fmt.Errorf("failed to read raw data. Error: %w", err)
	}

	resp, err := u.client.Put(ctx, &objectstorev1.PutRequest{
		Key:      reference.String(),
		Metadata: m,
		Object: &objectstorev1.Object{
			Contents: data,
		},
	})

	if int64(resp.SizeBytes) != size {
		return fmt.Errorf("size mismatch. Expected: %d, got: %d", size, resp.SizeBytes)
	}

	return nil
}

func (u UnionMetaStore) Delete(ctx context.Context, reference DataReference) error {
	//TODO implement me
	panic("implement me")
}

func NewUnionMetaStore(ctx context.Context, cfg *Config, metrics *dataStoreMetrics) (RawStore, error) {
	client, err := v1.NewClient(ctx, cfg.UnionMeta.Connection)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client. Error: %w", err)
	}

	store := &UnionMetaStore{
		client: client,
	}

	store.copyImpl = newCopyImpl(store, metrics.copyMetrics)
	return store, nil
}
