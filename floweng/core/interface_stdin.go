package core

import (
	"bufio"
	"errors"
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
	quit        chan chan error
	fromScanner chan stdinMsg
}

func NewStdin(eq chan<- BlockState) Source {
	stdin := &Stdin{
		SourceCommon: SourceCommon{eq, make(map[*Block]chan interface{})},
		quit:         make(chan chan error),
		fromScanner:  make(chan stdinMsg),
	}
	return stdin
}

func (stdin *Stdin) AddLink(b *Block, l chan interface{}) {
	stdin.SourceCommon.AddLink(b, l)
}

func (stdin *Stdin) RemoveLink(b *Block, l chan interface{}) {
	stdin.SourceCommon.RemoveLink(b, l)
}

func (stdin *Stdin) Serve() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		txt := scanner.Text()
		stdin.fromScanner <- stdinMsg{txt, nil}
	}
	err := scanner.Err()
	if err != nil {
		for block := range stdin.Links {
			stdin.EventQueue <- block.createState()
		}
		for _, link := range stdin.Links {
			link <- stdinMsg{"", err}
		}
	} else {
		stdin.fromScanner <- stdinMsg{"", errors.New("EOF")}
	}
}

func (stdin *Stdin) Stop() {
}

func (stdin Stdin) ReceiveMessage(i chan Interrupt, recChan <-chan interface{}) (string, Interrupt, error) {
	select {
	case i := <-recChan:
		msg := i.(stdinMsg)
		log.Println("FROM SCANNER", msg.err == nil, msg.err)
		return msg.msg, nil, msg.err
	case f := <-i:
		return "", f, nil
	}
}

func StdinReceive() Spec {
	return Spec{
		Name:    "stdinReceive",
		Outputs: []Pin{Pin{"msg", STRING}},
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
