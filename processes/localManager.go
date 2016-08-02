package processes

//
// This guy manages starting and stopping novi and edge processes on the
// localhost. Thus, no sshing, etc.
//

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type (
	localManagerFactory struct{}
	localManager        struct {
		cmd  *exec.Cmd
		name string
		args []string
	}
)

func init() {
	AddProcessManagerFactory("local", &localManagerFactory{})
}

func (l *localManagerFactory) MatchAddress(address string) bool {
	return strings.Contains(address, "127.0.0.") || strings.Contains(address, "localhost")
}

func (l *localManagerFactory) GetManager(name string, args []string) ProcessManager {
	return &localManager{name: name, args: args}
}

func (l *localManager) Start() {
	path := getPathFor(l.name)
	sOut, sErr := getOutputFiles(l.name)
	l.cmd = startCommand(path, l.args, sOut, sErr)
}

func startCommand(path string, args []string, sOut, sErr *os.File) *exec.Cmd {
	cmd := exec.Command(path, args...)
	cmd.Stdout = sOut
	cmd.Stderr = sErr

	//  Detach child process so when we exit,
	//  novi keeps running. That's all this flag is doing
	//  It is a Unix Thing :-)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not start command %s: %s\n", path, err)
	}
	return cmd
}

func (l *localManager) Stop() {
	if l.cmd == nil {
		return
	}
	if err := l.cmd.Process.Kill(); err != nil {
		log.Printf("Warning: Could not kill %s: %s\n", l.name, err)
	}
}

func (l *localManager) GetPid() int {
	return int(l.cmd.Process.Pid)
}

func (l *localManager) SetPid(pid int) {
	if l.cmd == nil {
		l.cmd = exec.Command(getPathFor(l.name))
	}
	l.cmd.Process, _ = os.FindProcess(pid) // Docs say on unix FindProcess always succeeds. !?
}

func getPathFor(commandName string) string {
	path, err := exec.LookPath(commandName)
	if err != nil {
		log.Fatalf("Could not find command '%s' in your search path: %s\n", commandName, err)
		os.Exit(1)
	}
	return path
}

func getOutputFiles(fileBaseName string) (sOut *os.File, sErr *os.File) {
	sOut, err := os.OpenFile(fileBaseName+".out", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not open command output file for '%s': %s\n", fileBaseName, err)
	}
	sErr, err = os.OpenFile(fileBaseName+".err", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Could not open command error output file for '%s': %s\n", fileBaseName, err)
	}
	return
}
