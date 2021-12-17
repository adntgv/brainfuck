package main

import (
	"bufio"
	"fmt"
)

type op string

const (
	right op = ">"
	left  op = "<"
	add   op = "+"
	sub   op = "-"
	print op = "."
	read  op = ","
	start op = "["
	end   op = "]"
)

type Processor struct {
	stack       []int64
	pointer     uint64
	operations  map[op]func() error
	debug       bool
	input       bufio.Reader
	depth       uint64
	ip          uint64
	subroutines map[uint64]uint64
	memory      []byte
	printer     Printer
}

func NewProcessor() *Processor {
	processor := &Processor{
		stack:       make([]int64, 1),
		debug:       true,
		memory:      make([]byte, 0),
		subroutines: make(map[uint64]uint64),
		printer:     NewPrinter(),
	}
	processor.operations = map[op]func() error{
		right: processor.right,
		left:  processor.left,
		add:   processor.add,
		sub:   processor.sub,
		print: processor.print,
		read:  processor.read,
		start: processor.start,
		end:   processor.end,
	}

	return processor
}

func (p *Processor) Print(b byte) {
	fmt.Printf("stack: %v\tpointer: %v\tpointer val: %v\n", p.stack, p.pointer, p.getPointerVal())
	fmt.Printf("memory: %s\tip: %v\tinstruction: %c\n\n", p.memory, p.ip, b)
}

func (p *Processor) right() error {
	if len(p.stack)-1 <= int(p.pointer) {
		p.stack = append(p.stack, 0)
	}
	p.pointer++
	return nil
}

func (p *Processor) left() error {
	if p.pointer == 0 {
		return fmt.Errorf("reached left limit")
	}
	p.pointer--
	return nil
}

func (p *Processor) add() error {
	p.stack[p.pointer] += 1
	return nil
}

func (p *Processor) sub() error {
	p.stack[p.pointer] -= 1
	return nil
}

func (p *Processor) print() error {
	p.printer.Printf("%c", p.getPointerVal())
	return nil
}

func (p *Processor) read() error {
	b, err := p.input.ReadByte()
	if err != nil {
		return fmt.Errorf("could not read byte: %v", err)
	}
	p.setPointerVal(int64(b))
	return nil
}

func (p *Processor) process(b byte, shouldStore bool) error {
	if shouldStore {
		p.memory = append(p.memory, b)
	}

	operation, known := p.operations[op(b)]
	if !known {
		//log.Printf("unkown operation '%v'", string(b))
		return nil
	}

	err := operation()
	if err != nil {
		return err
	}

	p.ip++

	return nil
}

func (p *Processor) start() error {
	p.subroutines[p.depth] = p.ip + 1
	p.depth++
	return nil
}

func (p *Processor) end() error {
	p.depth--
	for p.getPointerVal() != 0 {
		p.ip = p.subroutines[p.depth]
		for _, b := range p.getLoopBody(p.depth) {
			err := p.process(b, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) getLoopBody(depth uint64) []byte {
	start := p.getCorrespondingStart(depth)
	end := p.ip - 1
	return p.memory[start:end]
}

func (p *Processor) getCorrespondingStart(depth uint64) uint64 {
	return p.subroutines[depth]
}

func (p *Processor) getPointerVal() int64 {
	return p.stack[p.pointer]
}

func (p *Processor) setPointerVal(v int64) {
	p.stack[p.pointer] = v
}
