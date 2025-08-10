package config

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func UserAgent() string {
	return UserAgentProduct + "/" + Version + " (commit=" + Commit + ")"
}
