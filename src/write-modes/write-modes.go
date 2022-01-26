package write_modes

type WriteMode string

const (
	ReadOnly         WriteMode = "read_only"          //Only Read Point Value.
	WriteOnce        WriteMode = "write_once"         //Write the value on COV, don't Read.
	WriteAlways      WriteMode = "write_always"       //Write the value on every poll (poll rate defined by setting).
	WriteThenRead    WriteMode = "write_then_read"    //Write the value on COV, then Read on each poll (poll rate defined by setting).
	WriteAndMaintain WriteMode = "write_and_maintain" //Write the value on COV, then Read on each poll (poll rate defined by setting). If the Read value does not match the Write value, Write the value again.
)
