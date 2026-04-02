package buildinfo

import "saiao/internal/api"

const Service = "saiao"

var (
	Version = "dev"
	Commit  = "unknown"
)

func Current() api.BuildInfo {
	return api.BuildInfo{
		Service: Service,
		Version: Version,
		Commit:  Commit,
	}
}
