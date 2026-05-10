package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
)

const ReleaseLayerMediaType = "application/vnd.sloth.contracts.release.v1+json"

type Contract struct {
	Name          string         `json:"name"`
	Label         string         `json:"label,omitempty"`
	Version       string         `json:"version,omitempty"`
	SchemaVersion string         `json:"schemaVersion,omitempty"`
	ContentHash   string         `json:"contentHash,omitempty"`
	Payload       map[string]any `json:"payload,omitempty"`
}

type Release struct {
	Version       string     `json:"version"`
	SchemaVersion string     `json:"schemaVersion,omitempty"`
	Contracts     []Contract `json:"contracts"`
}

type OCIClientOptions struct {
	Host                  string
	Repository            string
	AuthorizationToken    string
	UseAuthorizationToken bool
	PlainHTTP             bool
}

type OCIClient struct {
	repo *remote.Repository
}

func NewOCIClient(opts OCIClientOptions) (*OCIClient, error) {
	host := strings.TrimSpace(opts.Host)
	if host == "" {
		return nil, fmt.Errorf("registry host is required")
	}
	repository := strings.TrimSpace(opts.Repository)
	if repository == "" {
		return nil, fmt.Errorf("registry repository is required")
	}

	repo, err := remote.NewRepository(fmt.Sprintf("%s/%s", host, repository))
	if err != nil {
		return nil, fmt.Errorf("create remote repository: %w", err)
	}
	repo.PlainHTTP = shouldUsePlainHTTP(host, opts.PlainHTTP)

	if opts.UseAuthorizationToken && strings.TrimSpace(opts.AuthorizationToken) != "" {
		repo.Client = &auth.Client{
			Credential: auth.StaticCredential(host, auth.Credential{
				AccessToken: strings.TrimSpace(opts.AuthorizationToken),
			}),
		}
	}

	return &OCIClient{repo: repo}, nil
}

func shouldUsePlainHTTP(host string, explicit bool) bool {
	if explicit {
		return true
	}

	hostname := host
	if h, _, err := net.SplitHostPort(host); err == nil {
		hostname = h
	}
	hostname = strings.TrimPrefix(hostname, "[")
	hostname = strings.TrimSuffix(hostname, "]")

	return hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1"
}

func (c *OCIClient) ListVersions(ctx context.Context) ([]string, error) {
	if c == nil || c.repo == nil {
		return nil, fmt.Errorf("oci client is not initialized")
	}

	versions := make([]string, 0)
	err := c.repo.Tags(ctx, "", func(tags []string) error {
		versions = append(versions, tags...)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list OCI tags: %w", err)
	}

	return sortVersions(versions), nil
}

func (c *OCIClient) ResolveVersion(ctx context.Context, requested string) (string, error) {
	trimmed := strings.TrimSpace(requested)
	if trimmed != "" && trimmed != "latest" {
		return trimmed, nil
	}

	versions, err := c.ListVersions(ctx)
	if err != nil {
		return "", err
	}
	if len(versions) == 0 {
		return "", fmt.Errorf("registry has no release versions")
	}

	return versions[len(versions)-1], nil
}

func (c *OCIClient) PullRelease(ctx context.Context, requestedVersion string) (*Release, string, error) {
	if c == nil || c.repo == nil {
		return nil, "", fmt.Errorf("oci client is not initialized")
	}

	resolvedVersion, err := c.ResolveVersion(ctx, requestedVersion)
	if err != nil {
		return nil, "", err
	}

	_, manifestBytes, err := oras.FetchBytes(ctx, c.repo, resolvedVersion, oras.DefaultFetchBytesOptions)
	if err != nil {
		return nil, "", fmt.Errorf("fetch OCI manifest for %q: %w", resolvedVersion, err)
	}

	manifest := &ocispec.Manifest{}
	if err := json.Unmarshal(manifestBytes, manifest); err != nil {
		return nil, "", fmt.Errorf("parse OCI manifest for %q: %w", resolvedVersion, err)
	}

	layer, err := pickReleaseLayer(manifest.Layers)
	if err != nil {
		return nil, "", err
	}

	releaseBytes, err := content.FetchAll(ctx, c.repo, layer)
	if err != nil {
		return nil, "", fmt.Errorf("fetch OCI release payload for %q: %w", resolvedVersion, err)
	}

	release := &Release{}
	if err := json.Unmarshal(releaseBytes, release); err != nil {
		return nil, "", fmt.Errorf("parse OCI release payload for %q: %w", resolvedVersion, err)
	}

	if strings.TrimSpace(release.Version) == "" {
		release.Version = resolvedVersion
	}

	return release, resolvedVersion, nil
}

func pickReleaseLayer(layers []ocispec.Descriptor) (ocispec.Descriptor, error) {
	for _, layer := range layers {
		if layer.MediaType == ReleaseLayerMediaType {
			return layer, nil
		}
	}
	return ocispec.Descriptor{}, fmt.Errorf("manifest does not contain %q layer", ReleaseLayerMediaType)
}

func sortVersions(tags []string) []string {
	if len(tags) == 0 {
		return tags
	}

	type parsedVersion struct {
		raw string
		sem *semver.Version
	}
	parsed := make([]parsedVersion, 0, len(tags))
	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			parsed = append(parsed, parsedVersion{raw: tag})
			continue
		}
		parsed = append(parsed, parsedVersion{raw: tag, sem: v})
	}

	sort.Slice(parsed, func(i, j int) bool {
		left := parsed[i]
		right := parsed[j]
		if left.sem != nil && right.sem != nil {
			return left.sem.LessThan(right.sem)
		}
		if left.sem != nil {
			return false
		}
		if right.sem != nil {
			return true
		}
		return left.raw < right.raw
	})

	out := make([]string, 0, len(parsed))
	for _, item := range parsed {
		out = append(out, item.raw)
	}
	return out
}
