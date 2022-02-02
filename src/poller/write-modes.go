package poller

type WriteMode string

const (
	ReadOnce          WriteMode = "read_once"            //Only Read Point Value Once.
	ReadOnly          WriteMode = "read_only"            //Only Read Point Value (poll rate defined by setting).
	WriteOnce         WriteMode = "write_once"           //Write the value on COV, don't Read.
	WriteOnceReadOnce WriteMode = "write_once_read_once" //Write the value on COV, Read Once.
	WriteAlways       WriteMode = "write_always"         //Write the value on every poll (poll rate defined by setting).
	WriteOnceThenRead WriteMode = "write_once_then_read" //Write the value on COV, then Read on each poll (poll rate defined by setting).
	WriteAndMaintain  WriteMode = "write_and_maintain"   //Write the value on COV, then Read on each poll (poll rate defined by setting). If the Read value does not match the Write value, Write the value again.
)
