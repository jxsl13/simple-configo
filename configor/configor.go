package configor

const (
	ConfigorStructTagName = "configo"

	ConfigorStructTagOptionMandatory = "mandatory"
	ConfigorStructTagOptionDefault   = "default"

	ConfigorStructTagOptionPseudo      = "pseudo"
	ConfigorStructTagOptionDescription = "description"
)

const (
	KeyCaseEqual keyCase = iota + 1
	KeyCaseEnv
	KeyCaseSnake
)

// keyCase is the type of the enum of the different ways to expect a key to be cased.
type keyCase uint

var (
	// DefaultStructTag is the default struct tag in case the 'configo' struct tag is misisng.
	// comma delimited Name,Option,Option,...
	// Leave the Name part empty in order to use the keys based on the struct field name with the
	// below expected Case, e.g. with KeyCaseEnv your struct field Key is expected to be KEY in the map
	DefaultStructTag = `configo:","`

	// DefaultKeyCase is the default case that the struct field name
	DefaultKeyCase = KeyCaseEnv

	DefaultListDelimiter     = ","
	DefaultPairDelimiter     = ";"
	DefaultKeyValueDelimiter = "=>"
)
