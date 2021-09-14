package core

var toggleInt int

func ToggleInt() Spec {
	return Spec{
		Name:    "toggle-number",
		Category: []string{"toggle"},
		Inputs:  []Pin{},
		Outputs: []Pin{Pin{"out", NUMBER}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			toggleInt = (toggleInt + 1) % 2
			out[0] = toggleInt
			return nil
		},
	}
}


var toggleBool int
func ToggleBool() Spec {
	return Spec{
		Name:    "toggle-boolean",
		Category: []string{"toggle"},
		Inputs:  []Pin{},
		Outputs: []Pin{Pin{"out", BOOLEAN}},
		Kernel: func(in, out, internal MessageMap, s Source, i chan Interrupt) Interrupt {
			toggleBool = (toggleBool + 1) % 2
			if toggleBool == 1{
				out[0] = true
			} else {
				out[0] = false
			}

			return nil
		},
	}
}

