package vm

import "fmt"

// load r1 addr : Load value at given address into given register
// store r2 addr : store r2 addr: store the value in register at the
// given memory address
// add r1 r2: Set r1 = r1 + r2
// sub r1 r2: Set r1 = r1 - r2
const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
func compute(memory []byte) {

	registers := [3]byte{
		8, // PC
		0, // R1
		0, // R2
	}

	// Keep looping, like a physical computer's clock
loop:
	for {
		programCounter := registers[0]
		op := memory[programCounter]
		switch op {
		case Load:
			registerIdx := memory[programCounter+1]
			dataIdx := memory[programCounter+2]
			registers[registerIdx] = memory[dataIdx]

			registers[0] += 3
		case Store:
			registerIdx := memory[programCounter+1]
			dataIdx := memory[programCounter+2]
			memory[dataIdx] = registers[registerIdx]

			registers[0] += 3
		case Add:
			mutatedRegisterIdx := memory[programCounter+1]
			addendRegisterIdx := memory[programCounter+2]
			registers[mutatedRegisterIdx] += registers[addendRegisterIdx]

			registers[0] += 3
		case Sub:
			mutatedRegisterIdx := memory[programCounter+1]
			minuendRegisterIdx := memory[programCounter+2]
			registers[mutatedRegisterIdx] -= registers[minuendRegisterIdx]

			registers[0] += 3
		case Addi:
			registerIdx := memory[programCounter+1]
			addend := memory[programCounter+2]
			registers[registerIdx] += addend

			registers[0] += 3
		case Subi:
			registerIdx := memory[programCounter+1]
			minuend := memory[programCounter+2]
			registers[registerIdx] -= minuend

			registers[0] += 3
		case Jump:
			registers[0] = memory[programCounter+1]
		case Beqz:
			registerIdx := memory[programCounter+1]
			if registers[registerIdx] == 0 {
				relativeOffset := memory[programCounter+2]
				registers[0] += relativeOffset
			}
			// we have to process these instructions regardless of processing
			// the offset argument
			registers[0] += 3
		case Halt:
			break loop
		default:
			panic(fmt.Errorf("unrecognized operation %x", op))
		}

	}
}
