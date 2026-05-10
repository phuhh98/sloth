package registry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func TestPullReleaseAndResolveLatestVersion(t *testing.T) {
	t.Parallel()

	releaseJSON := []byte(`{"version":"0.0.2","schemaVersion":"0.0.2","contracts":[{"name":"hero-banner","label":"Hero Banner","version":"0.0.2","schemaVersion":"0.0.2","contentHash":"abc"}]}`)
	layerDesc := descriptorFor(ReleaseLayerMediaType, releaseJSON)
	manifestBytes := manifestJSON(layerDesc)
	manifestDesc := descriptorFor(ocispec.MediaTypeImageManifest, manifestBytes)

	server := newOCIRegistryServer(t, map[string][]string{
		"org/contracts": {"0.0.1", "0.0.2"},
	}, map[string][]byte{
		"org/contracts:0.0.2": manifestBytes,
		"org/contracts@" + layerDesc.Digest.String(): releaseJSON,
		"org/contracts@" + manifestDesc.Digest.String(): manifestBytes,
	})
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")
	client, err := NewOCIClient(OCIClientOptions{
		Host:       host,
		Repository: "org/contracts",
		PlainHTTP:  true,
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	release, resolved, err := client.PullRelease(context.Background(), "latest")
	if err != nil {
		t.Fatalf("pull release: %v", err)
	}
	if resolved != "0.0.2" {
		t.Fatalf("expected resolved 0.0.2, got %q", resolved)
	}
	if release.Version != "0.0.2" {
		t.Fatalf("expected release version 0.0.2, got %q", release.Version)
	}
	if len(release.Contracts) != 1 || release.Contracts[0].Name != "hero-banner" {
		t.Fatalf("unexpected contracts payload: %+v", release.Contracts)
	}
}

func TestListVersionsSortsSemverThenLexical(t *testing.T) {
	t.Parallel()

	server := newOCIRegistryServer(t, map[string][]string{
		"org/contracts": {"alpha", "0.0.10", "0.0.2"},
	}, nil)
	defer server.Close()

	host := strings.TrimPrefix(server.URL, "http://")
	client, err := NewOCIClient(OCIClientOptions{Host: host, Repository: "org/contracts", PlainHTTP: true})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	versions, err := client.ListVersions(context.Background())
	if err != nil {
		t.Fatalf("list versions: %v", err)
	}
	got := strings.Join(versions, ",")
	want := "alpha,0.0.2,0.0.10"
	if got != want {
		t.Fatalf("expected versions %q, got %q", want, got)
	}
}

func TestShouldUsePlainHTTP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		host     string
		explicit bool
		want     bool
	}{
		{name: "explicit true", host: "ghcr.io", explicit: true, want: true},
		{name: "localhost", host: "localhost:5000", explicit: false, want: true},
		{name: "loopback ipv4", host: "127.0.0.1:5000", explicit: false, want: true},
		{name: "loopback ipv6", host: "[::1]:5000", explicit: false, want: true},
		{name: "remote host", host: "ghcr.io", explicit: false, want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := shouldUsePlainHTTP(tc.host, tc.explicit)
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func descriptorFor(mediaType string, payload []byte) ocispec.Descriptor {
	d := digest.FromBytes(payload)
	return ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    d,
		Size:      int64(len(payload)),
	}
}

func manifestJSON(layer ocispec.Descriptor) []byte {
	manifest := ocispec.Manifest{
		Versioned: specs.Versioned{SchemaVersion: 2},
		MediaType: ocispec.MediaTypeImageManifest,
		Config:    descriptorFor(ocispec.MediaTypeImageConfig, []byte(`{}`)),
		Layers:    []ocispec.Descriptor{layer},
	}
	raw, _ := json.Marshal(manifest)
	return raw
}

func newOCIRegistryServer(t *testing.T, tagsByRepo map[string][]string, payloadByKey map[string][]byte) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/v2/" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if !strings.HasPrefix(path, "/v2/") {
			http.NotFound(w, r)
			return
		}

		rest := strings.TrimPrefix(path, "/v2/")
		if strings.HasSuffix(rest, "/tags/list") {
			repo := strings.TrimSuffix(rest, "/tags/list")
			repo = strings.TrimSuffix(repo, "/")
			tags := tagsByRepo[repo]
			_ = json.NewEncoder(w).Encode(map[string]any{"name": repo, "tags": tags})
			return
		}

		if strings.Contains(rest, "/manifests/") {
			parts := strings.SplitN(rest, "/manifests/", 2)
			repo := parts[0]
			ref := parts[1]
			key := repo + ":" + ref
			if strings.HasPrefix(ref, "sha256:") {
				key = repo + "@" + ref
			}
			payload, ok := payloadByKey[key]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
			_, _ = w.Write(payload)
			return
		}

		if strings.Contains(rest, "/blobs/") {
			parts := strings.SplitN(rest, "/blobs/", 2)
			repo := parts[0]
			d := parts[1]
			key := repo + "@" + d
			payload, ok := payloadByKey[key]
			if !ok {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write(payload)
			return
		}

		http.NotFound(w, r)
	}))
}
