package core


func And2() Spec {
	return Spec{
		Name:    "and2",
		Inputs:  []Pin{Pin{"in1", BOOLEAN}, Pin{"in2", BOOLEAN}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			x, ok := in[0].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			y, ok := in[1].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			out[0] = x && y
			return nil
		},
	}
}

func And() Spec {
	return Spec{
		Name:    "and",
		Inputs:  []Pin{Pin{"in1", BOOLEAN}, Pin{"in2", BOOLEAN}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			x, ok := in[0].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			y, ok := in[1].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			out[0] = x && y
			return nil
		},
	}
}

func Or() Spec {
	return Spec{
		Name:    "or",
		Inputs:  []Pin{Pin{"in", BOOLEAN}, Pin{"in", BOOLEAN}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			x, ok := in[0].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			y, ok := in[1].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			out[0] = x || y
			return nil
		},
	}
}


func Not() Spec {
	return Spec{
		Name:    "not",
		Inputs:  []Pin{Pin{"in", BOOLEAN}},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			x, ok := in[0].(bool)
			if !ok {
				out[0] = NewError("need boolean")
				return nil
			}
			out[0] = !x
			return nil
		},
	}
}


// GreaterThan returns true if value[0] > value[1] or false otherwise
func GreaterThan() Spec {
	return Spec{
		Name:     ">",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}},
		Outputs:  []Pin{Pin{"x>y", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("GreaterThan requires float on x")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("GreaterThan requires float on y")
				return nil
			}
			out[0] = d1 > d2
			return nil
		},
	}
}

// LessThan returns true if value[0] < value[1] or false otherwise
func LessThan() Spec {
	return Spec{
		Name:     "<",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}},
		Outputs:  []Pin{Pin{"x<y", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("LessThan requires floats")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("LessThan requires floats")
				return nil
			}
			out[0] = d1 < d2
			return nil
		},
	}
}

// EqualTo returns true if value[0] == value[1] or false otherwise
func EqualTo() Spec {
	return Spec{
		Name:     "==",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", ANY}, Pin{"y", ANY}},
		Outputs:  []Pin{Pin{"x==y", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			out[0] = in[0] == in[1]
			return nil
		},
	}
}

// NotEqualTo returns true if value[0] != value[1] or false otherwise
func NotEqualTo() Spec {
	return Spec{
		Name:     "!=",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", ANY}, Pin{"y", ANY}},
		Outputs:  []Pin{Pin{"x!=y", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			out[0] = in[0] != in[1]
			return nil
		},
	}
}

