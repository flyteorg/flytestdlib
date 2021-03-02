package k8s

import (
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core"
	"github.com/lyft/flytepropeller/pkg/controller/executors"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type kubeClientObj struct {
	client client.Client
	cache  cache.Cache
}

func (k *kubeClientObj) GetClient() client.Client {
	return k.client
}

func (k *kubeClientObj) GetCache() cache.Cache {
	return k.cache
}

func newKubeClientObj(c client.Client, cache cache.Cache) core.KubeClient {
	return &kubeClientObj{client: c, cache: cache}
}

func GetKubeClient(config ClusterConfig) (core.KubeClient, error) {
	kubeConf, err := NewK8sClusterConfig(config.Endpoint, config.Auth)
	if err != nil {
		return nil, err
	}

	mapper, err := apiutil.NewDynamicRESTMapper(kubeConf)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(kubeConf, cache.Options{Mapper: mapper})
	if err != nil {
		return nil, err
	}

	c, err := client.New(kubeConf, client.Options{Mapper: mapper})
	if err != nil {
		return nil, err
	}

	fallbackClient := executors.NewFallbackClient(&client.DelegatingClient{
		Reader: &client.DelegatingReader{
			CacheReader:  cache,
			ClientReader: c,
		},
		Writer:       c,
		StatusClient: c,
	}, c)

	return newKubeClientObj(fallbackClient, cache), nil
}
