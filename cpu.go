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
		SP:		^uint16(0), //SP = uint16(0) - 1 represents empty stack
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
	switch opcode & 0xF000 {
	case 0:
		switch opcode & 0xFF {
		case 0xE0:
			// clear screen
			cpu.PC += 2
		case 0xEE:
			log.Debugf("%x - return from function \n", opcode)
			cpu.PC = cpu.stack[cpu.SP]
			cpu.SP--
		}
	case 0x1000:
		log.Debugf("%x - jump to address %x\n", opcode, opcode & 0xFFF)
		cpu.PC = opcode & 0xFFF
	case 0x2000:
		log.Debugf("%x - call function at %x\n", opcode, opcode & 0xFFF)
		cpu.SP++
	case 0x3000:
		log.Debugf("%x - skip next opcode if Vx == NN \n", opcode)
		if cpu.V[(opcode & 0x0F00) >> 8] == uint8(opcode & 0xFF) {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x4000:
		log.Debugf("%x - skip next opcode if Vx != NN \n", opcode)
		if cpu.V[(opcode & 0x0F00) >> 8] != uint8(opcode & 0xFF) {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x5000:
		log.Debugf("%x - skip next opcode if Vx == Vy \n", opcode)
		if cpu.V[(opcode & 0x0F00) >> 8] == cpu.V[(opcode & 0x00F0) >> 4] {
			cpu.PC += 4
		} else {
			cpu.PC += 2
		}
	case 0x6000:
		log.Debugf("%x - set Vx to NN \n", opcode)
		cpu.V[(opcode & 0x0F00) >> 8] = uint8(opcode & 0xFF)
		cpu.PC += 2
	case 0x7000:
		log.Debugf("%x - add NN to Vx \n", opcode)
		cpu.V[(opcode & 0x0F00) >> 8] += uint8(opcode & 0xFF)
		cpu.PC += 2
	case 0x8000:
		switch opcode & 0x000F {
		case 0:
			log.Debugf("%x - set Vx to value of Vy \n", opcode)
			cpu.V[(opcode & 0x0F00) >> 8] = cpu.V[(opcode & 0x0F0) >> 4]
			cpu.PC += 2
		case 1:
			log.Debugf("%x - set Vx to Vx OR Vy \n", opcode)
			cpu.V[(opcode & 0x0F00) >> 8] = cpu.V[(opcode & 0x0F00) >> 8] |
											cpu.V[(opcode & 0x0F0) >> 4]
			cpu.PC += 2
		case 2:
			log.Debugf("%x - set Vx to Vx OR Vy \n", opcode)
			cpu.V[(opcode & 0x0F00) >> 8] = cpu.V[(opcode & 0x0F00) >> 8] &
											cpu.V[(opcode & 0x0F0) >> 4]
			cpu.PC += 2
		case 3:
			log.Debugf("%x - set Vx to Vx OR Vy \n", opcode)
			cpu.V[(opcode & 0x0F00) >> 8] = cpu.V[(opcode & 0x0F00) >> 8] ^
											cpu.V[(opcode & 0x0F0) >> 4]
			cpu.PC += 2
		}
	default:
		// change later to exit on unknown opcode
		log.Debugf("Unknown opcode %x\n", opcode)
		cpu.PC += 2
	}
}
