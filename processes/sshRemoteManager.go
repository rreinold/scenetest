package processes

//
// This guy manages starting and stopping novi and edge processes on the
// localhost. Thus, no sshing, etc.
//

import (
	//"log"
	//"os"
	//pif "github.com/clearblade/scenetest/processes/interfaces"
	"os"
	"os/exec"
	//"strings"
	//"syscall"
)

type (
	sshRemoteFactory struct{}
	sshRemoteManager struct {
		cmd  *exec.Cmd
		name string
		args []string
	}
)

func init() {
	AddProcessManagerFactory("sshRemote", &sshRemoteFactory{})
}

func (srm *sshRemoteFactory) MatchAddress(address string) bool {
	return false // TODO for now
}

func (srm *sshRemoteFactory) GetManager(name string, args []string) ProcessManager {
	return &sshRemoteManager{name: name, args: args}
}

func (srm *sshRemoteManager) Start() {
}

func (srm *sshRemoteManager) Stop() {
	// ssh in and execute a "kill <pid>" command
	if srm.cmd == nil {
		return
	}
}

func (srm *sshRemoteManager) GetPid() int {
	return 0
}

func (srm *sshRemoteManager) SetPid(pid int) {
	if srm.cmd == nil {
		srm.cmd = exec.Command(srm.name)
	}
	srm.cmd.Process, _ = os.FindProcess(pid)
}
