package main

import (
	"os"
	"bufio"
)

func main() {
	master := NewMaster("0.0.0.0:9800")
	master.Start()
	
	stdReader := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := stdReader.ReadLine()
		cmd := string(line)
		if cmd == "exit" {
			break
		}
		
		master.Send(line)
	}
}