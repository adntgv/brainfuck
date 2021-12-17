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
	stack      []int64
	pointer    uint64
	operations map[op]func() error
	debug      bool
	input      bufio.Reader
	ip         uint64
	memory     []byte
	printer    Printer
}

func NewProcessor() *Processor {
	processor := &Processor{
		stack:   make([]int64, 1),
		debug:   true,
		memory:  make([]byte, 0),
		printer: NewPrinter(),
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

	if p.debug {
		//p.Print(b)
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
	if shouldStore {
		p.ip++
	}

	return nil
}

func (p *Processor) start() error {
	return nil
}

func (p *Processor) clone() *Processor {
	clone := NewProcessor()
	clone.stack = p.stack
	clone.pointer = p.pointer
	clone.debug = p.debug
	clone.operations = p.operations
	clone.printer = p.printer

	return clone
}

func (p *Processor) end() error {
	for p.getPointerVal() != 0 {
		operations := p.getLoopBody()
		for _, b := range operations {
			err := p.process(b, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Processor) getLoopBody() []byte {
	operations := make([]byte, len(p.memory))

	for i := p.ip; p.memory[i] != '['; i-- {
		operations = append(operations, p.memory[i])
	}

	// reverse
	for i, j := 0, len(operations)-1; i < j; i, j = i+1, j-1 {
		operations[i], operations[j] = operations[j], operations[i]
	}

	return operations
}

func (p *Processor) getPointerVal() int64 {
	return p.stack[p.pointer]
}

func (p *Processor) setPointerVal(v int64) {
	p.stack[p.pointer] = v
}
