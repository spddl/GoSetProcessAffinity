package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// -flag
// -flag=x
// -flag x  // non-boolean flags only

var (
	exeFlag  string
	procFlag string
	pidFlag  string
	pid      []int
	cpuFlag  string
	cpu      []int

	logging        bool
	suspend        bool
	resume         bool
	boost          bool
	priorityClass  int
	ioPriority     int
	memoryPriority int
)

const (
	// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getpriorityclass
	ABOVE_NORMAL_PRIORITY_CLASS uint32 = 0x00008000
	BELOW_NORMAL_PRIORITY_CLASS uint32 = 0x00004000
	HIGH_PRIORITY_CLASS         uint32 = 0x00000080
	IDLE_PRIORITY_CLASS         uint32 = 0x00000040
	NORMAL_PRIORITY_CLASS       uint32 = 0x00000020
	REALTIME_PRIORITY_CLASS     uint32 = 0x00000100
)

var CPUint []uint64
var PRIORITY_CLASS_Map map[uint32]string
var PRIORITY_CLASS map[int]uint32

func init() {
	flag.StringVar(&pidFlag, "pid", "", "pid")            // eine PID die schon läuft
	flag.StringVar(&procFlag, "proc", "", "proc")         // eine Exe die schon läuft
	flag.StringVar(&exeFlag, "exe", "", "a program/game") // eine Exe die erst gestartet wird
	flag.StringVar(&cpuFlag, "cpu", "", "0,1,2,3")
	flag.BoolVar(&logging, "logging", false, "logging")
	flag.BoolVar(&suspend, "suspend", false, "Suspend")
	flag.BoolVar(&resume, "resume", false, "Resume")
	flag.BoolVar(&boost, "boost", false, "https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setprocesspriorityboost")
	flag.IntVar(&priorityClass, "priorityClass", -1, "https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setpriorityclass")
	flag.IntVar(&ioPriority, "ioPriority", -1, "ioPriority")
	flag.IntVar(&memoryPriority, "memoryPriority", -1, "memoryPriority")

	flag.Parse()
	// log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	rand.Seed(time.Now().UnixNano())

	if logging {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)

		file, err := os.OpenFile(filepath.Join(exPath, "debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	}

	if pidFlag == "" && procFlag == "" && exeFlag == "" {
		log.Panic("nothing to do")
	}

	if cpuFlag != "" {
		for _, val := range strings.Split(cpuFlag, ",") {
			i, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			cpu = append(cpu, i)
		}
	}

	if pidFlag != "" {
		for _, val := range strings.Split(pidFlag, ",") {
			i, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			pid = append(pid, i)
		}
	}

	PRIORITY_CLASS_Map = map[uint32]string{
		ABOVE_NORMAL_PRIORITY_CLASS: "ABOVE_NORMAL_PRIORITY_CLASS",
		BELOW_NORMAL_PRIORITY_CLASS: "BELOW_NORMAL_PRIORITY_CLASS",
		HIGH_PRIORITY_CLASS:         "HIGH_PRIORITY_CLASS",
		IDLE_PRIORITY_CLASS:         "IDLE_PRIORITY_CLASS",
		NORMAL_PRIORITY_CLASS:       "NORMAL_PRIORITY_CLASS",
		REALTIME_PRIORITY_CLASS:     "REALTIME_PRIORITY_CLASS",
	}

	var index uint64 = 1
	for i := 0; i < int(runtime.NumCPU()); i++ {
		CPUint = append(CPUint, index)
		index *= 2
	}

	PRIORITY_CLASS = map[int]uint32{
		0: IDLE_PRIORITY_CLASS,
		1: BELOW_NORMAL_PRIORITY_CLASS,
		2: NORMAL_PRIORITY_CLASS,
		3: ABOVE_NORMAL_PRIORITY_CLASS,
		4: HIGH_PRIORITY_CLASS,
		5: REALTIME_PRIORITY_CLASS,
	}

	switch priorityClass {
	case -1:
		// placeholder
	case 0:
		log.Println("priorityClass: Idle")
	case 1:
		log.Println("priorityClass: Below Normal")
	case 2:
		log.Println("priorityClass: Normal")
	case 3:
		log.Println("priorityClass: Above Normal")
	case 4:
		log.Println("priorityClass: High")
	case 5:
		log.Println("priorityClass: Realtime")
	default:
		log.Println("priorityClass Error")
	}

	switch ioPriority {
	case -1:
		// placeholder
	case 0:
		log.Println("IoPriority: Very Low")
	case 1:
		log.Println("IoPriority: Low")
	case 2:
		log.Println("ioPriority: Normal")
	case 3:
		log.Println("ioPriority: High")
	case 4:
		log.Println("ioPriority: Critical")
	default:
		log.Println("IoPriority Error")
	}

	switch memoryPriority {
	case -1:
		// placeholder
	case 1:
		log.Println("memoryPriority: Very Low")
	case 2:
		log.Println("memoryPriority: Low")
	case 3:
		log.Println("memoryPriority: Medium")
	case 4:
		log.Println("memoryPriority: Below Normal")
	case 5:
		log.Println("memoryPriority: Normal")
	case 6:
		log.Println("memoryPriority: UNDEFINED_HIGH") // STATUS_INVALID_PARAMETER
	case 7:
		log.Println("memoryPriority: UNDEFINED_HIGHEST") // STATUS_INVALID_PARAMETER
	default:
		log.Println("memoryPriority Error")
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Has(b, flag uint32) bool { return b&flag != 0 }
