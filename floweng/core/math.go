package core

import "math"

// Addition returns the sum of the addenda
func Addition() Spec {
	return Spec{
		Name:     "+",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}},
		Outputs:  []Pin{Pin{"x+y", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			a1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Addition requires floats")
				return nil
			}
			a2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Addition requires floats")
				return nil
			}
			out[0] = a1 + a2
			return nil
		},
	}
}

// Subtraction returns the difference of the minuend - subtrahend
func Subtraction() Spec {
	return Spec{
		Name:     "-",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}, Pin{"yy", NUMBER}},
		Outputs:  []Pin{Pin{"x-yy", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			minuend, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Subtraction requires floats")
				return nil
			}
			subtrahend, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Subtraction requires floats")
				return nil
			}
			out[0] = minuend - subtrahend
			return nil
		},
	}
}

// Multiplication returns the product of the multiplicanda
func Multiplication() Spec {
	return Spec{
		Name:     "*",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}},
		Outputs:  []Pin{Pin{"x*y", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			m1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Multiplication requires floats")
				return nil
			}
			m2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Multiplication requires floats")
				return nil
			}
			out[0] = m1 * m2
			return nil
		},
	}
}

// Division returns the quotient of the dividend / divisor
func Division() Spec {
	return Spec{
		Name:     "/",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"x", NUMBER}, Pin{"y", NUMBER}},
		Outputs:  []Pin{Pin{"x/y", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Division requires floats")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Division requires floats")
				return nil
			}
			out[0] = d1 / d2
			return nil
		},
	}
}

// Exponentiation returns the base raised to the exponent
func Exponentiation() Spec {
	return Spec{
		Name:     "^",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"base", NUMBER}, Pin{"exponent", NUMBER}},
		Outputs:  []Pin{Pin{"power", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Exponentiation requires floats")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Exponentiation requires floats")
				return nil
			}
			out[0] = math.Pow(d1, d2)
			return nil
		},
	}
}

// Modulation returns the remainder of the dividend mod divisor
func Modulation() Spec {
	return Spec{
		Name:     "mod",
		Category: []string{"maths"},
		Inputs:   []Pin{Pin{"dividend", NUMBER}, Pin{"divisor", NUMBER}},
		Outputs:  []Pin{Pin{"remainder", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			d1, ok := in[0].(float64)
			if !ok {
				out[0] = NewError("Modulation requires floats")
				return nil
			}
			d2, ok := in[1].(float64)
			if !ok {
				out[0] = NewError("Modultion requires floats")
				return nil
			}
			out[0] = math.Mod(d1, d2)
			return nil
		},
	}
}
