package test

import (
	"runtime"
	"testing"
)


func TestSystem(t *testing.T)  {
	os := runtime.GOOS
	arch:=runtime.GOARCH
	version:=runtime.Version()

	t.Logf("Operating System: %s\nArchitecture: %s\nGo Version: %s", os, arch, version)
}
