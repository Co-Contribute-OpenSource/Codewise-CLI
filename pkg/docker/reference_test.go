package docker

import "testing"

func TestParseImageReference(t *testing.T) {
	tests := []struct {
		input, repository, tag, digest string
	}{
		{"nginx", "nginx", "", ""},
		{"nginx:1.27", "nginx", "1.27", ""},
		{"localhost:5000/team/api:v2", "localhost:5000/team/api", "v2", ""},
		{"ghcr.io/team/api@sha256:abcdef", "ghcr.io/team/api", "", "sha256:abcdef"},
		{"localhost:5000/api@sha256:abcdef", "localhost:5000/api", "", "sha256:abcdef"},
	}
	for _, tt := range tests {
		got, err := ParseImageReference(tt.input)
		if err != nil {
			t.Fatalf("ParseImageReference(%q): %v", tt.input, err)
		}
		if got.Repository != tt.repository || got.Tag != tt.tag || got.Digest != tt.digest {
			t.Errorf("ParseImageReference(%q) = %#v", tt.input, got)
		}
	}
}

func TestParseImageReferenceRejectsInvalidValues(t *testing.T) {
	for _, input := range []string{"", "bad image", ":latest", "image:", "image@sha256"} {
		if _, err := ParseImageReference(input); err == nil {
			t.Errorf("ParseImageReference(%q) expected error", input)
		}
	}
}
