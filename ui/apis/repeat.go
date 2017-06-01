package apis

//String used to recognize commands directed to this module
const REPEATCHANNEL = "repeat"

//Enum for possible user actions
type RepeaterAction int

const (
	//Creates a new repeater entry
	CREATE RepeaterAction = iota
	//Performs the request
	GO
	//Retrieves an history item
	GET
)

var repeaterActions = [...]string{
	"create",
	"go",
	"get",
}

func (a RepeaterAction) String() string {
	return repeaterActions[a]
}
