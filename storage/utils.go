package storage

import (
	"context"
	"os"

	stdErrs "github.com/flyteorg/flytestdlib/errors"
	"github.com/flyteorg/flytestdlib/promutils/labeled"

	"github.com/graymeta/stow"
	"github.com/pkg/errors"
)

var (
	ErrExceedsLimit       stdErrs.ErrorCode = "LIMIT_EXCEEDED"
	ErrFailedToWriteCache stdErrs.ErrorCode = "CACHE_WRITE_FAILED"
)

const (
	genericFailureTypeLabel = "Generic"
)

// Gets a value indicating whether the underlying error is a Not Found error.
func IsNotFound(err error) bool {
	if root := errors.Cause(err); os.IsNotExist(root) {
		return true
	}

	if stdErrs.IsCausedByError(err, stow.ErrNotFound) {
		return true
	}

	return false
}

// Gets a value indicating whether the underlying error is "already exists" error.
func IsExists(err error) bool {
	if root := errors.Cause(err); os.IsExist(root) {
		return true
	}

	return false
}

// Gets a value indicating whether the root cause of error is a "limit exceeded" error.
func IsExceedsLimit(err error) bool {
	return stdErrs.IsCausedBy(err, ErrExceedsLimit)
}

func IsFailedWriteToCache(err error) bool {
	return stdErrs.IsCausedBy(err, ErrFailedToWriteCache)
}

func MapStrings(mapper func(string) string, strings ...string) []string {
	if strings == nil {
		return []string{}
	}

	for i, str := range strings {
		strings[i] = mapper(str)
	}

	return strings
}

func incFailureCounterForError(ctx context.Context, counter labeled.Counter, err error) {
	errCode, found := stdErrs.GetErrorCode(err)
	if found {
		counter.Inc(context.WithValue(ctx, FailureTypeLabel, errCode))
	} else {
		counter.Inc(context.WithValue(ctx, FailureTypeLabel, genericFailureTypeLabel))
	}
}
