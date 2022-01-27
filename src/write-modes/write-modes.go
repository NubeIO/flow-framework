package write_modes

type WriteMode string

//Write Modes Summary:
//  - Write Modes are used to manage when the point should be written/read from a protocol client.
//  - The Write Mode will dictate if a point polling is repeated

//Questions:
// - at what level should we specify the fast, normal, and slow poll rates?  Plugin? Network? Device?  I'm thinking Device level
// - Should write values should be given a higher priority in the poll queue
// - How do I get FF Points by UUID?
// - Are FF Points shared by multiple plugins? NO?
//     - Can I store a Timer as a new property in FF Points?
//     - Can I store a PollRate and PollPriority in FF Points?
//     - Can I store a PollRate and PollPriority in FF Points?

const (
	ReadOnly         WriteMode = "read_only"          //Only Read Point Value.
	WriteOnce        WriteMode = "write_once"         //Write the value on COV, don't Read.
	WriteAlways      WriteMode = "write_always"       //Write the value on every poll (poll rate defined by setting).
	WriteThenRead    WriteMode = "write_then_read"    //Write the value on COV, then Read on each poll (poll rate defined by setting).
	WriteAndMaintain WriteMode = "write_and_maintain" //Write the value on COV, then Read on each poll (poll rate defined by setting). If the Read value does not match the Write value, Write the value again.
)
