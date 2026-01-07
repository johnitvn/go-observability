package observability

// Global build metadata (injected via LDFlags)
var (
	ServiceName = ""
	Version     = "dev"
	BuildTime   = "unknown"
)

// Accessor helpers for build metadata.
// Using functions avoids potential typecheck issues for tools
// that evaluate package symbols differently.
func GetServiceName() string { return ServiceName }
func GetVersion() string { return Version }
func GetBuildTime() string { return BuildTime }
