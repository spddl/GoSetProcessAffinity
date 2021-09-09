package main

import (
	"log"
	"strconv"
	"strings"
	"unsafe"

	"os/exec"
	"os/user"

	"golang.org/x/sys/windows"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Panic(err)
	}

	if user.Uid != "S-1-5-18" { // https://docs.microsoft.com/de-de/windows/security/identity-protection/access-control/security-identifier
		var pidString = []string{pidFlag}
		if exeFlag != "" {
			cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", exeFlag) // https://github.com/api0cradle/UltimateAppLockerByPassList/blob/master/DLL-Execution.md#urldll---fileprotocolhandler
			err = cmd.Start()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			pidString = []string{strconv.Itoa(cmd.Process.Pid)}
		}

		if procFlag != "" {
			proc, err := processes()
			if err != nil {
				log.Panic(err)
			}
			pidString = findProcessIDByName(proc, procFlag)
		}

		t := Task{name: "GoSetProcessAffinity-" + RandStringBytes(10)}
		defer func() {
			for {
				err := t.deleteTask()
				if err == nil {
					break
				}
			}
		}()

		t.createTask(strings.Join(pidString, ","))
		t.runTask()

	} else { // SystemKontext

		if len(pid) == 0 {
			log.Panic("nothing to do")
		}

		for _, pidVal := range pid {
			var processId = uint32(pidVal)
			log.Println("processId", processId)

			pHndl, err := windows.OpenProcess(ProcessSetIinformation, false, processId)
			defer windows.CloseHandle(pHndl)
			if err != nil {
				log.Panic(err)
			}
			if pHndl == 0 {
				log.Panic("no handle")
			}

			info(pHndl)

			if boost {
				SetProcessPriorityBoost(pHndl, false) // If the parameter is FALSE, dynamic boosting is enabled.
			}

			if len(cpu) != 0 {
				var cpuMask uint64
				for _, cpuVal := range cpu {
					cpuMask += CPUint[cpuVal]
				}
				err := SetProcessAffinityMask(pHndl, cpuMask)
				if err != nil {
					log.Println("setProcessAffinityMask: ", err)
				}
			}

			if priorityClass != -1 {
				err := windows.SetPriorityClass(pHndl, PRIORITY_CLASS[priorityClass])
				if err != nil {
					log.Println("SetPriorityClass: ", err)
				}
			}

			if ioPriority != -1 {
				var IoPriorityByte = uint32(ioPriority)
				ntStatus := NtSetInformationProcess(pHndl, ProcessIoPriority, &IoPriorityByte, 4)
				log.Println("NtStatus (ProcessIoPriority)", ntStatus)
			}

			if memoryPriority != -1 {
				var Memorypriorityuint32 = uint32(memoryPriority)
				ntStatus := NtSetInformationProcess(pHndl, ProcessPagePriority, &Memorypriorityuint32, 4)
				log.Println("NtStatus (ProcessPagePriority)", ntStatus)
			}

			info(pHndl)
		}
	}

}

func info(handle windows.Handle) {
	if handle == 0 {
		log.Println("The handle is invalid.", handle)
		return
	}

	log.Println()
	// find the CPU priority
	cpu, err := windows.GetPriorityClass(windows.Handle(handle))
	if err != nil {
		log.Println(err)
	}
	var cpulist string
	for b, n := range PRIORITY_CLASS_Map {
		if Has(b, cpu) {
			cpulist = n
			break
		}
	}

	// find the memory priority
	var memory uint32
	size := uint32(unsafe.Sizeof(memory))
	err = NtQueryInformationProcess(handle, ProcessPagePriority, windows.Pointer(unsafe.Pointer(&memory)), size, nil)
	if err != nil {
		log.Printf("NtQueryInformationProcess fails with %v\n", err)
	}

	// find the IO priority
	var io uint32
	size = uint32(unsafe.Sizeof(io))
	err = NtQueryInformationProcess(handle, ProcessIoPriority, windows.Pointer(unsafe.Pointer(&io)), size, nil)
	if err != nil {
		log.Printf("NtQueryInformationProcess fails with %v\n", err)
	}

	log.Printf("CPU priority:    %s\n", cpulist)
	log.Printf("Memory priority: %d (default: %d)\n", memory, 5)
	log.Printf("IO priority:     %d (default: %d)\n", io, 2)
	log.Println()
}
