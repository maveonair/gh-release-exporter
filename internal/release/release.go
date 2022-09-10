package release

// Release represents a release
type Release struct {
	LastKnownVersion string `toml:"last_known_version"`
	GitHubRepo       string `toml:"github_repo"`
}
