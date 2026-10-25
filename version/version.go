// Package version exposes the application version.
package version

// Version is set at build time when -ldflags "-X .../version.Version=..." is
// used; the default identifies the Web UI fork.
var Version = "0.02"
