package tools

const (
	ILLEGAL                = "ILLEGAL"
	EOF                    = "EOF"
	LEFT_BRA               = "LEFT_BRA"    // [
	RIGHT_BRA              = "RIGHT_BRA"   // ]
	COLUMN                 = "COLUMN"      // :
	SEMI_COLUMN            = "SEMI_COLUMN" // ;
	EQUAL                  = "EQUAL"       // =
	IDENTIFIER             = "IDENTIFIER"  // identifier = value
	LITERAL                = "LITERAL"
	QUOTE                  = "QUOTE"        // '
	DOUBLE_QUOTE           = "DOUBLE_QUOTE" // "
	PROFILE_NAME           = "PROFILE_NAME"
	INHERITED_PROFILE_NAME = "INHERITED_PROFILE_NAME"
	DEFAULT_SWITCH         = "DEFAULT_SWITCH"
)

type TokenType string

type Token struct {
	token_type TokenType
	value      string
	line       int16
}

type ProfileVariable struct {
	name  string
	value string
}

type Profile struct {
	name           string
	inherit_from   string
	default_switch bool
	variables      []ProfileVariable
}
