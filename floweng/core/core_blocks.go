package core

import (
	"encoding/json"
	"fmt"
	"time"
)

// First emits true when it receives its first message, and false otherwise
func First() Spec {
	return Spec{
		Name:     "first",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}},
		Outputs:  []Pin{{"first", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			_, ok := internal[0]
			if !ok {
				out[0] = true
				internal[0] = struct{}{}
			} else {
				out[0] = false
			}
			return nil
		},
	}
}

// Delay emits the message on pass through after the specified duration
func Delay() Spec {
	return Spec{
		Name:     "delay",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}, {"duration", STRING}},
		Outputs:  []Pin{{"out", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			t, err := time.ParseDuration(in[1].(string))
			if err != nil {
				out[0] = err
				return nil
			}
			timer := time.NewTimer(t)
			select {
			case <-timer.C:
				out[0] = in[0]
				return nil
			case f := <-i:
				return f
			}
		},
	}
}

// Set creates a new message with the specified key and value
func Set() Spec {
	return Spec{
		Name:     "set",
		Category: []string{"object"},
		Inputs:   []Pin{{"key", STRING}, {"value", ANY}},
		Outputs:  []Pin{{"object", OBJECT}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			out[0] = map[string]interface{}{
				in[0].(string): in[1],
			}
			return nil
		},
	}
}

// Get emits the value against the specified key
func Get() Spec {
	return Spec{
		Name:     "get",
		Category: []string{"object"},
		Inputs:   []Pin{{"in", OBJECT}, {"key", STRING}},
		Outputs:  []Pin{{"out", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			obj, ok := in[0].(map[string]interface{})
			if !ok {
				out[0] = NewError("inbound message must be an object")
			}
			key, ok := in[1].(string)
			if !ok {
				out[0] = NewError("key must be a string")
				return nil
			}
			out[0] = obj[key]
			return nil
		},
	}
}

// Console writes the inbound message to stdout
// TODO where should this write exactly?
// TODO there should be a stdout source block and Log should be a compound block with a writer
func Console() Spec {
	return Spec{
		Name:     "console",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"log", ANY}},
		Outputs:  []Pin{},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			o, err := json.Marshal(in[0])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(o))
			return nil
		},
	}
}

// Sink discards the inbound message
func Sink() Spec {
	return Spec{
		Name:     "sink",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}},
		Outputs:  []Pin{},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			return nil
		},
	}
}

// Latch emits the inbound message on the 0th output if ctrl is true,
// and the 1st output if ctrl is false
func Latch() Spec {
	return Spec{
		Name:     "latch",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}, {"ctrl", BOOLEAN}},
		Outputs:  []Pin{{"true", ANY}, {"false", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			controlSignal, ok := in[1].(bool)
			if !ok {
				out[0] = NewError("Latch ctrl requires bool")
				return nil
			}
			if controlSignal {
				out[0] = in[0]
			} else {
				out[1] = in[0]
			}
			return nil
		},
	}
}

// Gate emits the inbound message upon receiving a message on its trigger
func Gate() Spec {
	return Spec{
		Name:     "gate",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}, {"ctrl", ANY}},
		Outputs:  []Pin{{"out", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			out[0] = in[0]
			return nil
		},
	}
}

// Identity emits the inbound message immediately
func Identity() Spec {
	return Spec{
		Name:     "identity",
		Category: []string{"mechanism"},
		Inputs:   []Pin{{"in", ANY}},
		Outputs:  []Pin{{"out", ANY}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			out[0] = in[0]
			return nil
		},
	}
}

// Merge recursively merges two objects, favouring the first input to resolve conflicts
func Merge() Spec {
	return Spec{
		Name:     "merge",
		Category: []string{"object"},
		Inputs: []Pin{
			{"in", OBJECT},
			{"in", OBJECT},
		},
		Outputs: []Pin{{"out", OBJECT}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt, block *Block) Interrupt {
			result := make(map[string]interface{})
			var err error
			in0, ok := in[0].(map[string]interface{})
			if !ok {
				out[0] = NewError("Merge needs map")
				return nil
			}
			result, err = MergeMap(result, in0)
			if err != nil {
				out[0] = err
				return nil
			}
			in1, ok := in[1].(map[string]interface{})
			if !ok {
				out[0] = NewError("Merge needs map")
				return nil
			}
			result, err = MergeMap(result, in1)
			if err != nil {
				out[0] = err
				return nil
			}
			out[0] = result
			return nil
		},
	}
}
