package config

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func UserAgent() string {
	return "idpproxy/" + Version + " (commit=" + Commit + ")"
}
