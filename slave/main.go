package main

import (
	"os"
	"bufio"
)

func main() {
	slave := NewSlave("127.0.0.1:9800")
	slave.Start()
	
	stdReader := bufio.NewReader(os.Stdin)
	stdReader.ReadLine()
}