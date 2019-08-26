package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "0.0.7"
	Commit  = ""
)

type versionInfo struct {
	QOSVersion string `json:"version"`
	GitCommit  string `json:"commit"`
	GoVersion  string `json:"go"`
}

func (v versionInfo) String() string {
	return fmt.Sprintf(`qos: %s
git commit: %s
%s`, v.QOSVersion, v.GitCommit, v.GoVersion)
}

func newVersionInfo() versionInfo {
	return versionInfo{
		Version,
		Commit,
		fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)}
}
