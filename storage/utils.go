package storage

import (
	"context"
	"os"
	"regexp"

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

type SubexpName = string
type MatchedString = string

// MatchRegex returns all matches for the sub-expressions within the regex.
func MatchRegex(reg *regexp.Regexp, input string) map[SubexpName]MatchedString {
	names := reg.SubexpNames()
	res := reg.FindAllStringSubmatch(input, -1)
	if len(res) == 0 {
		return nil
	}

	dict := make(map[string]string, len(names))
	// Start from 1 since names[0] is always empty per docs on reg.SubexpNames()
	for i := 1; i < len(res[0]); i++ {
		dict[names[i]] = res[0][i]
	}

	return dict
}

// IsNotFound gets a value indicating whether the underlying error is a Not Found error.
func IsNotFound(err error) bool {
	if root := errors.Cause(err); os.IsNotExist(root) {
		return true
	}

	if stdErrs.IsCausedByError(err, stow.ErrNotFound) {
		return true
	}

	return false
}

// IsExists gets a value indicating whether the underlying error is "already exists" error.
func IsExists(err error) bool {
	if root := errors.Cause(err); os.IsExist(root) {
		return true
	}

	return false
}

// IsExceedsLimit gets a value indicating whether the root cause of error is a "limit exceeded" error.
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
