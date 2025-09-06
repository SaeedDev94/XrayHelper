package commands

import (
	"XrayHelper/main/builds"
	"XrayHelper/main/common"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"os"
	"path"
	"strconv"
	"time"
)

const tagService = "service"

type ServiceCommand struct{}

func (this *ServiceCommand) Execute(args []string) error {
	if err := builds.LoadConfig(); err != nil {
		return err
	}
	if len(args) == 0 {
		return e.New("not specify operation, available operation [start|stop|restart|status]").WithPrefix(tagService).WithPathObj(*this)
	}
	if len(args) > 1 {
		return e.New("too many arguments").WithPrefix(tagService).WithPathObj(*this)
	}
	log.HandleInfo("service: current core type is " + builds.Config.XrayHelper.CoreType)
	switch args[0] {
	case "start":
		log.HandleInfo("service: starting core")
		if err := startService(); err != nil {
			return err
		}
		log.HandleInfo("service: core is running, pid is " + getServicePid())
	case "stop":
		log.HandleInfo("service: stopping core")
		stopService()
		log.HandleInfo("service: core is stopped")
	case "restart":
		log.HandleInfo("service: restarting core")
		if err := restartService(); err != nil {
			return err
		}
		log.HandleInfo("service: core is running, pid is " + getServicePid())
	case "status":
		pidStr := getServicePid()
		if len(pidStr) > 0 {
			log.HandleInfo("service: core is running, pid is " + pidStr)
		} else {
			log.HandleInfo("service: core is stopped")
		}
	default:
		return e.New("unknown operation " + args[0] + ", available operation [start|stop|restart|status]").WithPrefix(tagService).WithPathObj(*this)
	}
	return nil
}

// newServices get core service
func newServices(serviceLogFile *os.File) (service common.External, err error) {
	return common.NewExternal(0, serviceLogFile, serviceLogFile, builds.Config.XrayHelper.CorePath, "run", "-c", builds.Config.XrayHelper.CoreConfig), nil
}

// getServicePid get core pid from pid file
func getServicePid() string {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "core.pid")); err == nil {
		pidFile, err := os.ReadFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
		if err != nil {
			log.HandleDebug(err)
		}
		return string(pidFile)
	} else {
		log.HandleDebug(err)
	}
	return ""
}

// startService start core service
func startService() error {
	// check current core status
	servicePid := getServicePid()
	if len(servicePid) > 0 {
		return e.New("core is running, pid is " + servicePid).WithPrefix(tagService)
	}
	// get core service log file
	serviceLogFile, err := os.OpenFile(path.Join(builds.Config.XrayHelper.RunDir, "error.log"), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		return e.New("open core log file failed, ", err).WithPrefix(tagService)
	}
	// get core service
	service, err := newServices(serviceLogFile)
	if err != nil {
		return err
	}
	// add core env variable and prepare core config
	service.AppendEnv("XRAY_LOCATION_ASSET=" + builds.Config.XrayHelper.DataDir)
	service.SetUidGid("0", common.CoreGid)
	service.Start()
	if service.Err() != nil {
		return e.New("start core service failed, ", service.Err()).WithPrefix(tagService)
	}
	if err := os.WriteFile(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"), []byte(strconv.Itoa(service.Pid())), 0644); err != nil {
		_ = service.Kill()
		stopService()
		return e.New("write core pid failed, ", err).WithPrefix(tagService)
	}
	return nil
}

// stopService stop core service
func stopService() {
	if _, err := os.Stat(path.Join(builds.Config.XrayHelper.RunDir, "core.pid")); err == nil {
		pidStr := getServicePid()
		if len(pidStr) > 0 {
			pid, _ := strconv.Atoi(pidStr)
			if serviceProcess, err := os.FindProcess(pid); err == nil {
				_ = serviceProcess.Kill()
				_ = os.Remove(path.Join(builds.Config.XrayHelper.RunDir, "core.pid"))
			} else {
				log.HandleDebug(err)
			}
		}
	} else {
		log.HandleDebug(err)
	}
}

// restartService restart core service
func restartService() error {
	stopService()
	time.Sleep(1 * time.Second)
	return startService()
}
