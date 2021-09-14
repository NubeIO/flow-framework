package core

import "log"

type stcoreError struct {
	S string `json:"core"`
}

func (e *stcoreError) Error() string {
	log.Println(e.S)
	return e.S
}

func NewError(s string) *stcoreError {
	return &stcoreError{
		S: s,
	}
}

// GetLibrary Library is the set of all core block Specs
// over all func's every time we need the library
func GetLibrary() map[string]Spec {
	b := []Spec{
		// mechanism
		First(),
		Delay(),
		Console(), //Console log
		Sink(),
		Latch(),
		Gate(),
		Identity(), //pass through
		Timestamp(),

		//count
		AnyCount(),

		//toggle
		ToggleInt(),
		ToggleBool(),

		//logic
		And(),
		Or(),
		Not(),
		GreaterThan(),
		LessThan(),
		EqualTo(),
		NotEqualTo(),

		// maths
		Addition(),
		Subtraction(),
		Multiplication(),
		Division(),
		Exponentiation(),
		Modulation(),
		Exp(),
		Floor(),
		Ceil(),
		Log10(),
		Log(),
		Sqrt(),
		Sin(),
		Cos(),
		Tan(),

		//assertions
		IsBoolean(),
		IsNumber(),
		IsString(),
		IsArray(),
		IsObject(),
		IsError(),

		// conversion
		ToString(),
		ToNumber(),

		MQTTClientConnect(),
		MQTTClientSend(),

		// random
		UniformRandom(),
		NormalRandom(),
		ZipfRandom(),
		PoissonRandom(),
		ExponentialRandom(),
		BernoulliRandom(),

		// string
		InString(),
		HasPrefix(),
		HasSuffix(),
		StringConcat(),
		StringSplit(),

		// key value
		kvGet(),
		kvSet(),
		kvClear(),
		kvDump(),
		kvDelete(),

		// parsers
		ParseJSON(),

		// object
		Set(),
		Get(),
		Keys(),
		Merge(),
		HasField(),

		// array
		Append(),
		Tail(),
		Head(),
		Init(),
		Last(),
		Len(),
		InArray(),

		// primitive value
		ValueGet(),
		ValueSet(),

		// list
		listGet(),
		listSet(),
		listShift(),
		listAppend(),
		listPop(),
		listDump(),

		// priority queue
		pqPush(),
		pqPop(),
		pqPeek(),
		pqLen(),
		pqClear(),

		// network IO
		HTTPRequest(),


		//string functions
		StringConcat(),
		StringSplit(),

		// websocket
		wsClientConnect(),
		wsClientReceive(),
		wsClientSend(),

		// stdin
		StdinReceive(),
	}

	library := make(map[string]Spec)

	for _, s := range b {
		library[s.Name] = s
	}

	return library
}
