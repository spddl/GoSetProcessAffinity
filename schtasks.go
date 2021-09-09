package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	name string
}

func (t *Task) createTask(Pid string) {
	appLocation, err := os.Executable()
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("schtasks", "/create", "/tn", t.name, "/tr", appLocation+" "+strings.Join(os.Args[1:], " ")+" -pid "+Pid, "/sc", "ONCE", "/sd", "01/01/2337", "/st", "00:00", "/ru", "SYSTEM", "/rl", "HIGHEST", "/f").CombinedOutput()
	if err != nil {
		log.Printf("createTask() err %s\n", out)
	}
}

func (t *Task) runTask() {
	out, err := exec.Command("schtasks", "/run", "/tn", t.name).CombinedOutput()
	if err != nil {
		log.Printf("runTask() err %+v %s\n", err, out)
	}
}

func (t *Task) deleteTask() error {
	cmd := exec.Command("schtasks", "/delete", "/tn", t.name, "/f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("deleteTask() err %+v %s\n", err, out)
	}
	return nil
}
