package rungo

import (
	"os"
	"syscall"
)

type NamespacesConfig struct {
	Pid    bool
	Mnt    string
	Net    bool
	Ipc    bool
	Uts    bool
	User   bool
	Cgroup bool
}

func (n *NamespacesConfig) Get() int {
	flags := 0

	if n.Pid {
		flags |= syscall.CLONE_NEWPID
	}
	if _, err := os.Stat(n.Mnt); !os.IsNotExist(err) {
		flags |= syscall.CLONE_NEWNS
	}
	if n.Net {
		flags |= syscall.CLONE_NEWNET
	}
	if n.Ipc {
		flags |= syscall.CLONE_NEWIPC
	}
	if n.Uts {
		flags |= syscall.CLONE_NEWUTS
	}
	if n.User {
		flags |= syscall.CLONE_NEWUSER
	}
	if n.Cgroup {
		flags |= syscall.CLONE_NEWCGROUP
	}

	return flags
}

type RungoConfig struct {
	ProcessPath      string
	Args             []string
	Hostname         string
	NamespacesConfig NamespacesConfig
}
