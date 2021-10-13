package core

import (
	"errors"
	"sync"
	"time"
)

// NewBlock creates a new block from a spec
func NewBlock(s Spec) *Block {
	var in []Input
	var out []Output

	for _, v := range s.Inputs {
		in = append(in, Input{
			Name:  v.Name,
			Type:  v.Type,
			Value: nil,
		})
	}

	for _, v := range s.Outputs {
		out = append(out, Output{
			Name:        v.Name,
			Type:        v.Type,
			Connections: make(map[Connection]struct{}),
		})
	}

	return &Block{
		routing: BlockRouting{
			Inputs:        in,
			Outputs:       out,
			InterruptChan: make(chan Interrupt),
		},
		kernel:     s.Kernel,
		sourceType: s.Source,
		Monitor:    make(chan MonitorMessage, 1),
		lastCrank:  time.Now(),
		done:       make(chan struct{}),
	}
}

func (b *Block) createState() BlockState {
	return BlockState{
		b,
		make(MessageMap),
		make(MessageMap),
		make(MessageMap),
		make(Manifest),
		false,
	}
}

func (b *Block) exportInput(id RouteIndex) (*Input, error) {

	if int(id) >= len(b.routing.Inputs) || int(id) < 0 {
		return nil, errors.New("index out of range")
	}

	if b.routing.Inputs[id].Value == nil {
		return &b.routing.Inputs[id], nil
	}

	return &Input{
		Value: &InputValue{
			Data: Copy((*b.routing.Inputs[id].Value).Data),
		},
		Name: b.routing.Inputs[id].Name,
	}, nil

}

// GetInput returns the specified Input
func (b *Block) GetInput(id RouteIndex) (Input, error) {
	b.routing.RLock()
	r, err := b.exportInput(id)
	if err != nil {
		b.routing.RUnlock()
		return Input{}, err
	}
	b.routing.RUnlock()
	return *r, err
}

// GetInputs returns all inputs for a block.
func (b *Block) GetInputs() []Input {
	b.routing.RLock()
	re := make([]Input, len(b.routing.Inputs), len(b.routing.Inputs))
	for i := range b.routing.Inputs {
		r, _ := b.exportInput(RouteIndex(i))
		re[i] = *r
	}
	b.routing.RUnlock()
	return re
}

// SetInput RouteValue sets the route to always be the specified value
func (b *Block) SetInput(id RouteIndex, v *InputValue) error {
	var returnVal error = nil
	if int(id) < 0 || int(id) >= len(b.routing.Inputs) {
		returnVal = errors.New("input out of range")
	}

	b.routing.Inputs[id].Value = v

	return returnVal
}

// GetOutputs Outputs return a list of manifest pairs for the block
func (b *Block) GetOutputs() []Output {
	b.routing.RLock()
	m := make([]Output, len(b.routing.Outputs), len(b.routing.Outputs))
	for id, out := range b.routing.Outputs {
		m[id] = Output{
			Name:        out.Name,
			Type:        out.Type,
			Connections: make(map[Connection]struct{}),
		}
		for k := range out.Connections {
			m[id].Connections[k] = struct{}{}
		}
	}
	b.routing.RUnlock()
	return m
}

func (b *Block) GetSource() Source {
	b.routing.RLock()
	v := b.routing.Source
	b.routing.RUnlock()
	return v
}

// SetSource sets a store for the block. can be set to nil
func (b *Block) SetSource(s Source) error {
	var returnVal error = nil
	if s != nil && s.GetType() != b.sourceType {
		returnVal = errors.New("invalid source type for this block")
	}
	b.routing.Source = s
	return returnVal
}

// Connect connects a Route, specified by ID, to a connection
func (b *Block) Connect(id RouteIndex, c Connection) error {
	var returnVal error = nil
	if int(id) < 0 || int(id) >= len(b.routing.Outputs) {
		returnVal = errors.New("output out of range")
	}

	if _, ok := b.routing.Outputs[id].Connections[c]; ok {
		returnVal = errors.New("this connection already exists on this output")
	}

	b.routing.Outputs[id].Connections[c] = struct{}{}
	return returnVal
}

// Disconnect removes a connection from a Input
func (b *Block) Disconnect(id RouteIndex, c Connection) error {
	var returnVal error = nil
	if int(id) < 0 || int(id) >= len(b.routing.Outputs) {
		returnVal = errors.New("output out of range")
	}

	if _, ok := b.routing.Outputs[id].Connections[c]; !ok {
		returnVal = errors.New("connection does not exist")
	}

	delete(b.routing.Outputs[id].Connections, c)
	return returnVal
}

func (b *BlockState) Serve() []BlockState {
	// defer func() {
	//     b.done <- struct{}{}
	// }()
	// for {

	if b.block.sourceType != NONE && b.block.routing.Source == nil {
		return nil
	}
	// TODO: deal with interrupts
	// var interrupt Interrupt
	// b.routing.RLock()
	_, success := b.receive()
	if !success {
		b.crank()
		return nil
	}

	_, success = b.process()
	if !success {
		b.crank()
		return nil
	}

	_, nextStates := b.broadcast()
	// if interrupt != nil {
	//     return nil
	// }

	b.crank()
	// b.routing.RUnlock()
	// b.routing.Lock()
	// if ok := interrupt(); !ok {
	//     b.routing.Unlock()
	//     return
	// }
	// b.routing.Unlock()
	// }
	return nextStates
}

func (b *BlockState) Reset() {
	b.crank()

	// reset block's state as well. currently, this only applies to a handful of
	// blocks, like GET and first.
	for k := range b.internalValues {
		delete(b.internalValues, k)
	}

	// if there are any messages on the input channels, flush them.
	// note: all blocks that are sending to this block MUST BE IN A
	// STOPPED STATE. if any block routines that possess this block's
	// input channel are in a RUNNING state, this flush will not work
	// because it will simply pull another message into the buffer.
	// for _, input := range b.routing.Inputs {
	//     select {
	//     case <-input.C:
	//     default:
	//     }
	// }
}

// TODO: fix
func (b *Block) Stop() {
	b.routing.InterruptChan <- func() bool {
		return false
	}
	<-b.done
}

// wait and listen for all kernel inputs to be filled.
func (b *BlockState) receive() (Interrupt, bool) {
	for id, input := range b.block.routing.Inputs {
		// b.block.Monitor <- MonitorMessage{
		//     BI_INPUT,
		//     id,
		// }

		if value, ok := b.inputValues[RouteIndex(id)]; ok {
			if input.Value != nil {
				input.Value.Data = Copy(value)
			} else {
				b.block.SetInput(RouteIndex(id), &InputValue{Data: Copy(value)})
			}
		} else if input.Value != nil {
			b.inputValues[RouteIndex(id)] = Copy(input.Value.Data)
			continue
		} else {
			return nil, false
		}

		// TODO: check interrupt
		// select {
		// // case m := <-input.C:
		// //     b.state.inputValues[RouteIndex(id)] = m
		// case f := <-b.routing.InterruptChan:
		//     return f
		// }
	}
	return nil, true
}

// run kernel on inputs, produce outputs
func (b *BlockState) process() (Interrupt, bool) {

	// b.block.Monitor <- MonitorMessage{
	//     BI_KERNEL,
	//     nil,
	// }

	if b.Processed {
		return nil, true
	}

	// block until connected to source if necessary
	if b.block.sourceType != NONE && b.block.routing.Source == nil {
		// select {
		// case f := <-b.block.routing.InterruptChan:
		//     return f
		// }
		return nil, false
	}

	// we should only be able to get here if
	// - we don't need an shared state
	// - we have an external shared state and it has been attached

	// if we have a store, lock it
	var store sync.Locker
	var ok bool
	if store, ok = b.block.routing.Source.(sync.Locker); ok {
		store.Lock()
	}

	// run the kernel
	interrupt := b.block.kernel(b.inputValues,
		b.outputValues,
		b.internalValues,
		b.block.routing.Source,        // TODO: probs have to lock source
		b.block.routing.InterruptChan, // TODO: deal with interrupts
		b.block)

	// unlock the store if necessary
	if store != nil {
		store.Unlock()
	}

	// if an interrupt was received, return it
	if interrupt != nil {
		return interrupt, true
	}

	b.Processed = true
	return nil, true
}

// broadcast the kernel output to all connections on all outputs.
func (b *BlockState) broadcast() (Interrupt, []BlockState) {
	nextStates := make([]BlockState, 0)
	for id, out := range b.block.routing.Outputs {
		// b.block.Monitor <- MonitorMessage{
		//     BI_OUTPUT,
		//     id,
		// }

		// if the output key is not present in the output map, then we
		// don't deliver any message
		_, ok := b.outputValues[RouteIndex(id)]
		if !ok {
			continue
		}
		// if there no connection for this output then wait until there
		// is one. that means we have to wait for an interrupt.
		// if len(out.Connections) == 0 {
		//     select {
		//     case f := <-b.routing.InterruptChan:
		//         return f
		//     }
		// }

		for c := range out.Connections {
			// check to see if we have delivered a message to this
			// connection for this block crank. if we have, then
			// skip this delivery.
			m := ManifestPair{id, c.TargetId}
			if _, ok := b.manifest[m]; ok {
				continue
			}

			newState := c.Target.createState()
			nextStates = append(nextStates, newState)

			newState.inputValues[c.RouteId] = Copy(b.outputValues[RouteIndex(id)])
			b.manifest[m] = struct{}{}
			// select {
			// case c <- b.state.outputValues[RouteIndex(id)]:
			//     // set that we have delivered the message.
			//     b.state.manifest[m] = struct{}{}
			// case f := <-b.routing.InterruptChan:
			//     return f
			// }
		}
	}
	return nil, nextStates
}

// cleanup all block state for this crank of the block
func (b *BlockState) crank() {
	for k := range b.inputValues {
		delete(b.inputValues, k)
	}
	for k := range b.outputValues {
		delete(b.outputValues, k)
	}
	for k := range b.manifest {
		delete(b.manifest, k)
	}
	b.Processed = false
}
