package core

var countInput int


func AnyCount() Spec {
	return Spec{
		Name:    "count-any",
		Category: []string{"toggle"},
		Inputs:   []Pin{Pin{"x", ANY}},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			countInput = countInput + 1
			out[0] = countInput
			return nil
		},
	}
}