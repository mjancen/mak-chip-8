package main

import (
	"encoding/hex"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	maxMem = 4096
)

// CPU is the type that contains all of the CPU components
type CPU struct {
	PC     uint16   // program counter
	I      uint16   // address register
	SP     uint16   // stack pointer
	stack  []uint16 // stack
	V      []uint8  // general purpose registers
	DT     uint8    // delay timer
	ST     uint8    // sound timer
	mem []uint8
	gfx    []uint8 // graphics memory
	key    []uint8 // keypad
}

// NewCPU initializes the sizes of the slices in the CPU struct
func newCPU() *CPU {
	cpu := CPU{
		PC:     0x200,
		stack:  make([]uint16, 16),
		V:      make([]uint8, 16),
		mem: make([]uint8, 4096),
		gfx:    make([]uint8, 2048),
		key:    make([]uint8, 16)}
	return &cpu
}

func (cpu CPU) viewMem() string {
	return hex.Dump(cpu.mem)
}

func (cpu CPU) viewReg() string {
	return hex.Dump(cpu.V)
}

func (cpu CPU) loadROM(filename string) error {
	fmt.Printf("Loading ROM file %s...\n", filename)
	romFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer romFile.Close()
	n, err := romFile.Read(cpu.mem[0x200:])
	if err != nil {
		return err
	}
	log.Debugln(n, "bytes read from", filename)
	log.Debugln("Memory:")
	log.Debugln(hex.Dump(cpu.mem))
	return nil
}

func (cpu CPU) fetchOp() uint16 {
	opcode := uint16(cpu.mem[cpu.PC]) << 8 | uint16(cpu.mem[cpu.PC + 1])
	return opcode
}

func (cpu *CPU) decodeAndExec(opcode uint16) {
	switch {
	case opcode & 0x1000 != 0:
		fmt.Printf("jump to address %x\n", opcode & 0xFFF)
		cpu.PC = opcode & 0xFFF
	case opcode & 0x2000 != 0:
		fmt.Printf("Call function at %x\n", opcode & 0xFFF)
		cpu.stack[cpu.SP] = cpu.PC
		cpu.SP++
		if int(cpu.SP) == len(cpu.stack) {
			fmt.Printf("STACK OVERFLOW")
			os.Exit(1)
		}
		cpu.PC = opcode & 0xFFF
	case opcode & 0x00EE != 0:
		cpu.PC = cpu.stack[cpu.SP]
		cpu.SP--
	default:
		fmt.Printf("Unknown opcode %x\n", opcode)
		cpu.PC += 2
	}
}
