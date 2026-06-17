package version

import (
	"runtime"
	"runtime/debug"
)

// Overridden at build time via -ldflags by GoReleaser
// (see ldflags in .goreleaser.yaml).
var (
	version = "unset"
	commit  = "none"
	date    = "unknown"
)

// Info returns build/version metadata for logging at startup. It prefers the
// values injected via -ldflags and falls back to the toolchain's VCS stamps
// (available for plain `go build`) when those weren't set.
func Info() (ver, vcsCommit, buildDate, goVer, arch string) {
	ver = version
	vcsCommit = commit
	buildDate = date
	goVer = runtime.Version()
	arch = runtime.GOOS + "/" + runtime.GOARCH

	if info, ok := debug.ReadBuildInfo(); ok {
		if ver == "unset" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			ver = info.Main.Version
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				if vcsCommit == "none" {
					vcsCommit = s.Value
				}
			case "vcs.time":
				if buildDate == "unknown" {
					buildDate = s.Value
				}
			}
		}
	}

	return
}
