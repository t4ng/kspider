package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	master, _ := NewMaster("0.0.0.0:9800")
	master.Start()

	stdReader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := stdReader.ReadLine()
		cmd := string(line)
		if cmd == "exit" {
			break
		} else if strings.HasPrefix(cmd, "http") {
			master.taskMgr.CreateTask(cmd)
			continue
		}

		master.Send(line)
	}
}
