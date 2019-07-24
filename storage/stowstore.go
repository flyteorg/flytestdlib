package storage

import (
	"fmt"

	"github.com/lyft/flytestdlib/promutils"

	"github.com/pkg/errors"

	"github.com/graymeta/stow"
	"github.com/graymeta/stow/azure"
	"github.com/graymeta/stow/google"
	"github.com/graymeta/stow/oracle"
	"github.com/graymeta/stow/s3"
	"github.com/graymeta/stow/swift"
)

var fQNFn = map[string]func(string) DataReference{
	s3.Kind: func(bucket string) DataReference {
		return DataReference(fmt.Sprintf("s3://%s", bucket))
	},
	google.Kind: func(bucket string) DataReference {
		return DataReference(fmt.Sprintf("gs://%s", bucket))
	},
	oracle.Kind: func(bucket string) DataReference {
		return DataReference(fmt.Sprintf("os://%s", bucket))
	},
	swift.Kind: func(bucket string) DataReference {
		return DataReference(fmt.Sprintf("sw://%s", bucket))
	},
	azure.Kind: func(bucket string) DataReference {
		return DataReference(fmt.Sprintf("afs://%s", bucket))
	},
}

func newStowRawStore(cfg *Config, metricsScope promutils.Scope) (RawStore, error) {
	if cfg.InitContainer == "" {
		return nil, fmt.Errorf("initContainer is required")
	}

	if cfg.Stow != nil {
		fn, ok := fQNFn[cfg.Stow.Kind]
		if !ok {
			return nil, errors.Errorf("unsupported stow.kind [%s], add support in flytestdlib?", cfg.Stow.Kind)
		}
		loc, err := stow.Dial(cfg.Stow.Kind, cfg.Stow.Config)
		if err != nil {
			return emptyStore, fmt.Errorf("unable to configure the storage for s3. Error: %v", err)
		}

		c, err := loc.Container(cfg.InitContainer)
		if err != nil {
			if IsNotFound(err) {
				c, err := loc.CreateContainer(cfg.InitContainer)
				if err != nil {
					return emptyStore, fmt.Errorf("unable to initialize container [%v]. Error: %v", cfg.InitContainer, err)
				}
				return NewStowRawStore(fn(c.Name()), c, metricsScope)
			}
			return emptyStore, err
		}

		return NewStowRawStore(fn(c.Name()), c, metricsScope)
	}
	return nil, errors.Errorf("stow configuration section is missing")
}
