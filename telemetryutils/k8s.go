package telemetryutils

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	K8S_SERVICE_NAME = "k8s-client"
	PACKAGE = "controller-runtime.pkg.client"
)

type K8sCacheWrapper struct {
	cache.Cache
}

func WrapK8sCache(c cache.Cache) cache.Cache {
	return &K8sCacheWrapper{c}
}

func (c *K8sCacheWrapper) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Cache/Get", PACKAGE))
	defer span.End()
	return c.Cache.Get(ctx, key, obj)
}

func (c *K8sCacheWrapper) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Cache/List", PACKAGE))
	defer span.End()
	return c.Cache.List(ctx, list, opts...)
}

type K8sClientWrapper struct {
	client.Client

	statusWriter *K8sStatusWriterWrapper
}

func WrapK8sClient(c client.Client) client.Client {
	return &K8sClientWrapper{c, &K8sStatusWriterWrapper{c.Status()}}
}

func (c *K8sClientWrapper) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/Get", PACKAGE))
	defer span.End()
	return c.Client.Get(ctx, key, obj)
}

func (c *K8sClientWrapper) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/List", PACKAGE))
	defer span.End()
	return c.Client.List(ctx, list, opts...)
}

func (c *K8sClientWrapper) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/Create", PACKAGE))
	defer span.End()
	return c.Client.Create(ctx, obj, opts...)
}

func (c *K8sClientWrapper) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/Delete", PACKAGE))
	defer span.End()
	return c.Client.Delete(ctx, obj, opts...)
}

func (c *K8sClientWrapper) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/Update", PACKAGE))
	defer span.End()
	return c.Client.Update(ctx, obj, opts...)
}

func (c *K8sClientWrapper) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/Patch", PACKAGE))
	defer span.End()
	return c.Client.Patch(ctx, obj, patch, opts...)
}

func (c *K8sClientWrapper) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.Client/DeleteAllOf", PACKAGE))
	defer span.End()
	return c.Client.DeleteAllOf(ctx, obj, opts...)
}

func (c *K8sClientWrapper) Status() client.StatusWriter {
	return c.statusWriter
}

type K8sStatusWriterWrapper struct {
	client.StatusWriter
}

func (s *K8sStatusWriterWrapper) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.StatusWriter/Update", PACKAGE))
	defer span.End()
	return s.StatusWriter.Update(ctx, obj, opts...)
}

func (s *K8sStatusWriterWrapper) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	ctx, span := NewSpan(ctx, K8S_SERVICE_NAME, fmt.Sprintf("%s.StatusWriter/Patch", PACKAGE))
	defer span.End()
	return s.StatusWriter.Patch(ctx, obj, patch, opts...)
}
