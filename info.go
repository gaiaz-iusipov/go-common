package app

import (
	"log/slog"
	"runtime"
	"time"
)

var (
	name      string
	version   string
	vcsRef    string
	buildDate string
	buildUser string
	startDate = time.Now().UTC()
)

func Name() string {
	return name
}

func Version() string {
	return version
}

func VCSRef() string {
	return vcsRef
}

func BuildDate() string {
	return buildDate
}

func BuildUser() string {
	return buildUser
}

func StartDate() time.Time {
	return startDate
}

type Info struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	VCSRef    string    `json:"vcs_ref"`
	BuildDate string    `json:"build_date"`
	BuildUser string    `json:"build_user"`
	StartDate time.Time `json:"start_date"`
	GoVersion string    `json:"go_version"`
	GoOS      string    `json:"go_os"`
	GoArch    string    `json:"go_arch"`
}

// GetInfo returns build and runtime information.
func GetInfo() Info {
	return Info{
		Name:      Name(),
		Version:   Version(),
		VCSRef:    VCSRef(),
		BuildDate: BuildDate(),
		BuildUser: BuildUser(),
		StartDate: StartDate(),
		GoVersion: runtime.Version(),
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
	}
}

func LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("name", Name()),
		slog.String("version", Version()),
		slog.String("vcs_ref", VCSRef()),
		slog.String("build_date", BuildDate()),
		slog.String("build_user", BuildUser()),
		slog.String("go_version", runtime.Version()),
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
	)
}
