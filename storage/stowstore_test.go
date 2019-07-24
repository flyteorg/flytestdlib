package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lyft/flytestdlib/promutils"
)

func Test_newStowRawStore(t *testing.T) {
	type args struct {
		cfg          *Config
		metricsScope promutils.Scope
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"google", args{&Config{}, promutils.NewTestScope()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newStowRawStore(tt.args.cfg, tt.args.metricsScope)
			if (err != nil) != tt.wantErr {
				t.Errorf("newStowRawStore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got, "Expected rawstore, found nil!")
		})
	}
}
