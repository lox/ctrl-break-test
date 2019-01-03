package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	libkernel32                  = syscall.MustLoadDLL("kernel32")
	procGenerateConsoleCtrlEvent = libkernel32.MustFindProc("GenerateConsoleCtrlEvent")
)

const (
	createNewProcessGroupFlag = 0x00000200
)

func main() {
	if os.Getenv(`CREATE_SUB_PROCESS`) == `1` {
		c := make(chan os.Signal, 1)
		signal.Notify(c)

		log.Printf("Running in sub process PID %d", os.Getpid())
		for i := 0; i < 5; i++ {
			select {
			case <-time.After(time.Second):
				log.Printf("Tick")
			case s := <-c:
				log.Printf("Got interrupted: %v", s)
				os.Exit(0)
			}
		}
		os.Exit(1)
	}

	log.Printf("Creating sub process from PID %d", os.Getpid())
	cmd := exec.Command(os.Args[0])
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT | createNewProcessGroupFlag,
	}
	cmd.Env = append(cmd.Env, `CREATE_SUB_PROCESS=1`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Started process as PID %d", cmd.Process.Pid)
	time.Sleep(time.Second)

	// See https://docs.microsoft.com/en-us/windows/console/generateconsolectrlevent
	r1, _, err := procGenerateConsoleCtrlEvent.Call(syscall.CTRL_BREAK_EVENT, uintptr(cmd.Process.Pid))
	if r1 == 0 {
		log.Printf("Error sending CTRL_BREAK_EVENT: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Process exited successfully")
}
