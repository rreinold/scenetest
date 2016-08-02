package processes

import ()

const ()

type ProcessManagerFactory interface {
	MatchAddress(address string) bool
	GetManager(name string, args []string) ProcessManager
}

type ProcessManager interface {
	Start()
	Stop()
	GetPid() int
	SetPid(pid int)
}

var (
	factories = map[string]ProcessManagerFactory{}
)

func AddProcessManagerFactory(name string, mgr ProcessManagerFactory) {
	factories[name] = mgr
}

func GetProcessManager(address, name string, args []string) ProcessManager {
	for _, factory := range factories {
		if factory.MatchAddress(address) {
			return factory.GetManager(name, args)
		}
	}
	return nil
}

func GetProcessManagerWithPid(address, name string, pid int) ProcessManager {
	for _, factory := range factories {
		if factory.MatchAddress(address) {
			mgr := factory.GetManager(name, []string{})
			mgr.SetPid(pid)
			return mgr
		}
	}
	return nil
}
