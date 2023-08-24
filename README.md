# rungo

rungo is a simple, lightweight CLI tool for creating and running containers - AKA a container runtime.
This project was created for experimenting with containers internals and better grasp the concepts of containerization.

## Features
- [x] Mount namespace
- [x] UTS namespace
- [x] IPC namespace
- [x] Network namespace
- [x] PID namespace
- [x] User namespace
- [ ] Cgroups

## Installation
```bash
go build -o rungo cmd/main.go
```

## Usage
```
Usage: ./rungo -run -uid [-mnt=/path/rootfs] [-uts [-hostname=new_hostname]] [-ipc] [-net] [-pid] [-command command]
  -mnt='/path/rootfs'           Enable Mount namespace
  -uts                          Enable UTS namespace
  -hostname='new_hostname'      Set a custom hostname into the container
  -ipc                          Enable IPC namespace
  -net                          Enable Network namespace
  -pid                          Enable PID namespace
  -uid                          Enable User namespace
  -command='command'            Command to run "/bin/sh" by default
```

## Examples
```bash
# Run a container with a custom hostname
./rungo -run -uid -uts -hostname=container1
```
```bash
# Run a container with a custom hostname and a custom command
./rungo -run -uid -uts -hostname=container1 -command=/bin/ls
```
```bash
# Run a fully isolated container
./rungo -run -uid -uts -ipc -net -pid -mnt=/path/rootfs
```
The rootfs directory can be created with the following command:
```bash
mkdir rootfs
```
Then, you can populate the rootfs directory with the following command:
```bash
docker export $(docker create busybox) | tar -C rootfs -xvf -
```
## References
- [Namespaces in operation, part 1: namespaces overview](https://lwn.net/Articles/531114/)
- [containers-from-scratch-with-golang](https://medium.com/@ssttehrani/containers-from-scratch-with-golang-5276576f9909)
