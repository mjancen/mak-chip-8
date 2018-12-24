package main

import (
	"fmt"
	"os"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	fmt.Println("Hello")
	if len(os.Args) < 2 {
		fmt.Printf("Provide file name of ROM as argument.\n")
		os.Exit(1)
	}
	filename := os.Args[1]

	cpu := newCPU()
	err := cpu.loadROM(filename)
	if err != nil {
		fmt.Printf("Error loading ROM file: %v\n", err)
		os.Exit(1)
	}

	var opcode uint16

	for {
		fmt.Printf("PC: %x ", cpu.PC)
		opcode = cpu.fetchOp()
		// fmt.Print(opcode)
		cpu.decodeAndExec(opcode)
		if cpu.PC > maxMem - 2 {
			fmt.Println("End loop.")
			break
		}
	}
}
