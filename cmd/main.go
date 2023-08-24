//go:build linux
// +build linux

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/amitschendel/rungo/pkg/rungo"
)

type options struct {
	run bool
	ns  bool
}

func main() {
	config := rungo.RungoConfig{
		ProcessPath: "/bin/sh",
	}
	options := new(options)

	flag.Usage = func() {
		help()
	}

	flag.StringVar(&config.NamespacesConfig.Mnt, "mnt", "", "")
	flag.BoolVar(&config.NamespacesConfig.Uts, "uts", false, "")
	flag.StringVar(&config.Hostname, "hostname", "", "")
	flag.BoolVar(&config.NamespacesConfig.Ipc, "ipc", false, "")
	flag.BoolVar(&config.NamespacesConfig.Net, "net", false, "")
	flag.BoolVar(&config.NamespacesConfig.Pid, "pid", false, "")
	flag.BoolVar(&config.NamespacesConfig.User, "uid", false, "")
	flag.BoolVar(&options.run, "run", false, "")
	flag.BoolVar(&options.ns, "ns", false, "")
	flag.Parse()

	rungo := rungo.Rungo{Config: &config}

	switch os.Args[1] {
	case "-run":
		rungo.Init()
	case "-ns":
		rungo.Run()
	default:
		help()
		fmt.Println()
		log.Fatal("Wrong arguments passed.")
	}
}

func help() {
	fmt.Println("Usage: ./Rungo -run -uid [-mnt=/path/rootfs] [-uts [-hostname=new_hostname]] [-ipc] [-net] [-pid]")
	fmt.Println("  -mnt='/path/rootfs'        Enable Mount namespace")
	fmt.Println("  -uts                       Enable UTS namespace")
	fmt.Println("  -hostname='new_hostname'   Set a custom hostname into the container")
	fmt.Println("  -ipc                       Enable IPC namespace")
	fmt.Println("  -net                       Enable Network namespace")
	fmt.Println("  -pid                       Enable PID namespace")
	fmt.Println("  -uid                       Enable User namespace")
}
