package storage

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"net/url"

	"github.com/flyteorg/flytestdlib/logger"
)

const (
	separator           = "/"
	storageURLFormatter = "%s://%s/%s"
)

type SignedURLPatternMatcher = *regexp.Regexp

var (
	SignedURLPattern SignedURLPatternMatcher = regexp.MustCompile(`https://((storage\.googleapis\.com/(?P<bucket_gcs>[^/]+))|((?P<bucket_s3>[^\.]+)\.s3\.amazonaws\.com)|(.*\.blob\.core\.windows\.net/(?P<bucket_az>[^/]+)))/(?P<path>[^?]*)`)
)

// URLPathConstructor implements ReferenceConstructor that assumes paths are URL-compatible.
type URLPathConstructor struct {
	scheme string
}

func formatStorageURL(scheme, bucket, path string) DataReference {
	return DataReference(fmt.Sprintf(storageURLFormatter, scheme, bucket, path))
}

func ensureEndingPathSeparator(path DataReference) DataReference {
	if len(path) > 0 && path[len(path)-1] == separator[0] {
		return path
	}

	return path + separator
}

func (c URLPathConstructor) FromSignedURL(_ context.Context, signedURL string) (DataReference, error) {
	if len(c.scheme) == 0 {
		return "", fmt.Errorf("scheme cannot be empty in order to interact with SignedURLs")
	}

	matches := MatchRegex(SignedURLPattern, signedURL)
	if bucket := matches["bucket"]; len(bucket) == 0 {
		return "", fmt.Errorf("failed to parse signedURL [%v]. Resulted in an empty bucket", signedURL)
	} else if path := matches["path"]; len(path) == 0 {
		return "", fmt.Errorf("failed to parse signedURL [%v]. Resulted in an empty path", signedURL)
	} else {
		ref := formatStorageURL(c.scheme, matches["bucket"], matches["path"])
		_, err := url.Parse(ref.String())
		return ref, err
	}
}

func (URLPathConstructor) ConstructReference(ctx context.Context, reference DataReference, nestedKeys ...string) (DataReference, error) {
	u, err := url.Parse(string(ensureEndingPathSeparator(reference)))
	if err != nil {
		logger.Errorf(ctx, "Failed to parse prefix: %v", reference)
		return "", errors.Wrap(err, fmt.Sprintf("Reference is of an invalid format [%v]", reference))
	}

	rel, err := url.Parse(strings.Join(MapStrings(func(s string) string {
		return strings.Trim(s, separator)
	}, nestedKeys...), separator))
	if err != nil {
		logger.Errorf(ctx, "Failed to parse nested keys: %v", reference)
		return "", errors.Wrap(err, fmt.Sprintf("Reference is of an invalid format [%v]", reference))
	}

	u = u.ResolveReference(rel)

	return DataReference(u.String()), nil
}

func NewURLPathConstructor(scheme string) URLPathConstructor {
	return URLPathConstructor{
		scheme: scheme,
	}
}
