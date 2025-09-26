package platform

import (
	"runtime"
)

type Platform struct {
	OS string
}

func New() *Platform {
	return &Platform{
		OS: runtime.GOOS,
	}
}

func (p *Platform) IsWindows() bool {
	return p.OS == "windows"
}

func (p *Platform) IsUnix() bool {
	return p.OS == "darwin" || p.OS == "linux"
}

func (p *Platform) GetListPortsCommand() string {
	if p.IsWindows() {
		return "netstat -ano"
	}
	return "lsof -i -P -n"
}

func (p *Platform) GetKillCommand() string {
	if p.IsWindows() {
		return "taskkill"
	}
	return "kill"
}
