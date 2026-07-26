package docker

import (
	"fmt"
	"strings"
)

// ImageReference is a parsed OCI-style image reference. Tag and Digest are
// mutually optional; a digest remains attached to RepositoryAndDigest.
type ImageReference struct {
	Repository string
	Tag        string
	Digest     string
}

func (r ImageReference) RepositoryAndDigest() string {
	if r.Digest != "" {
		return r.Repository + "@" + r.Digest
	}
	return r.Repository
}

// ParseImageReference handles registry ports, tags, and sha256 digests without
// mistaking the registry port for an image tag.
func ParseImageReference(value string) (ImageReference, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return ImageReference{}, fmt.Errorf("image reference is empty")
	}
	if strings.ContainsAny(value, " \t\r\n") {
		return ImageReference{}, fmt.Errorf("image reference must not contain whitespace")
	}

	var digest string
	if before, after, found := strings.Cut(value, "@"); found {
		if before == "" || after == "" || !strings.Contains(after, ":") {
			return ImageReference{}, fmt.Errorf("invalid image digest %q", value)
		}
		value, digest = before, after
	}

	lastSlash := strings.LastIndex(value, "/")
	lastColon := strings.LastIndex(value, ":")
	repository, tag := value, ""
	if lastColon > lastSlash {
		repository, tag = value[:lastColon], value[lastColon+1:]
	}
	if repository == "" || (lastColon > lastSlash && tag == "") {
		return ImageReference{}, fmt.Errorf("invalid image reference %q", value)
	}
	return ImageReference{Repository: repository, Tag: tag, Digest: digest}, nil
}
