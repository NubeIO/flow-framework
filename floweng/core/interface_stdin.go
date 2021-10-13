package core

import (
	"bufio"
	"log"
	"os"
)

func StdinInterface() SourceSpec {
	return SourceSpec{
		Name: "stdin",
		Type: STDIN,
		New:  NewStdin,
	}
}

func (stdin Stdin) GetType() SourceType {
	return STDIN
}

type stdinMsg struct {
	msg string
	err error
}

type Stdin struct {
	SourceCommon
	quit chan chan error
}

func NewStdin(eq chan<- BlockState) Source {
	stdin := &Stdin{
		SourceCommon: SourceCommon{eq, make(map[*Block]chan interface{})},
		quit:         make(chan chan error),
	}
	return stdin
}

func (stdin *Stdin) AddLink(b *Block, l chan interface{}) {
	stdin.SourceCommon.AddLink(b, l)
}

func (stdin *Stdin) RemoveLink(b *Block) {
	stdin.SourceCommon.RemoveLink(b)
}

func (stdin *Stdin) Serve() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		for block := range stdin.Links {
			stdin.EventQueue <- block.createState()
		}
		for _, link := range stdin.Links {
			link <- stdinMsg{txt, nil}
		}
	}
	err := scanner.Err()
	if err != nil {
		log.Println("STDIN ERROR:", err)
	} else {
		log.Println("STDIN EOF")
	}
}

func (stdin *Stdin) Stop() {
}

func (stdin Stdin) ReceiveMessage(i chan Interrupt, recChan <-chan interface{}) (string, Interrupt, error) {
	select {
	case i := <-recChan:
		msg := i.(stdinMsg)
		return msg.msg, nil, msg.err
	case f := <-i:
		return "", f, nil
	}
}

func StdinReceive() Spec {
	return Spec{
		Name:    "stdinReceive",
		Outputs: []Pin{{"msg", STRING}},
		Source:  STDIN,
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			stdin := s.(*Stdin)
			msg, f, err := stdin.ReceiveMessage(i, stdin.Links[block])
			if err != nil {
				out[0] = NewError("EOF")
				return nil
			}
			if f != nil {
				return f
			}
			out[0] = msg
			return nil
		},
	}
}
