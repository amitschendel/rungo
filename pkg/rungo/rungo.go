package rungo

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type Rungo struct {
	Config *RungoConfig
}

func (r *Rungo) Run() {
	r.setNamespaces()
	cmd := exec.Command(r.Config.ProcessPath, r.Config.Args...)
	cmd.Env = []string{"PS1=📦 [$(whoami)@$(hostname)] ~$(pwd) ‣ "}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	defer must(r.unsetProcessID())
}

func (r *Rungo) Init() {
	log.Info("Initiating container process!")
	cmd := exec.Command(CMD_PATH, append([]string{"-ns"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	flags := r.Config.NamespacesConfig.Get()

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: uintptr(flags),
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Run()
}

func (r *Rungo) setMntNs() (bool, error) {
	if r.Config.NamespacesConfig.Mnt != "" {
		if _, err := os.Stat(r.Config.NamespacesConfig.Mnt); !os.IsNotExist(err) {
			if err := syscall.Chroot(r.Config.NamespacesConfig.Mnt); err != nil {
				log.Error("Error setting MNT namespace")
				return false, errors.New("error setting MNT namespace")
			}
			if err := syscall.Chdir("/"); err != nil {
				log.Error("Error changing dir")
				return false, errors.New("error changing dir")
			}
		} else {
			return false, errors.New("error setting MNT namespace")
		}
	}
	return true, nil
}

func (r *Rungo) setPidNs() (bool, error) {
	if r.Config.NamespacesConfig.Mnt != "" {
		if r.Config.NamespacesConfig.Pid {
			if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
				log.Error("Error setting PID namespace")
				return false, errors.New("error setting PID namespace")
			}
			log.Info("PID namespace set\n")
			return true, nil
		}
	}

	return false, nil
}

func (r *Rungo) unsetProcessID() (bool, error) {
	if r.Config.NamespacesConfig.Pid {
		if err := syscall.Unmount("proc", 0); err != nil {
			log.Error("Error unsetting PID namespace")
			return false, errors.New("error unsetting PID namespace")
		}
		log.Info("PID namespace unset\n")
		return true, nil
	}
	return false, nil
}

func (r *Rungo) setIpcNs() (bool, error) {
	if r.Config.NamespacesConfig.Ipc {
		log.Info("Setting IPC namespace")
		return true, nil
	}
	return false, nil
}

func (r *Rungo) setNetNs() (bool, error) {
	if r.Config.NamespacesConfig.Net {
		log.Info("Setting NET namespace")
		return true, nil
	}
	return false, nil
}

func (r *Rungo) setUtsNs() (bool, error) {
	var containerHostname string
	if r.Config.NamespacesConfig.Uts {
		if r.Config.Hostname != "" {
			containerHostname = r.Config.Hostname
		} else {
			containerHostname = DEFAULT_HOSTNAME
		}
		if err := syscall.Sethostname([]byte(containerHostname)); err != nil {
			log.Printf("err: %v", err)
			return false, errors.New("error setting UTS namespace")
		}
	}
	return true, nil
}

func (r *Rungo) setUserNs() (bool, error) {
	if r.Config.NamespacesConfig.User {
		log.Info("Setting USER namespace")
		return true, nil
	}
	return false, nil
}

func (r *Rungo) setNamespaces() {
	must(r.setMntNs())
	must(r.setPidNs())
	must(r.setIpcNs())
	must(r.setNetNs())
	must(r.setUtsNs())
	must(r.setUserNs())
}

func must(result bool, err error) {
	if !result && err != nil {
		log.Fatal(err)
	}
}
