package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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

func main() {
	var tokens []Token = parseIniFile("my.ini")
	var profiles []Profile = buildProfiles(tokens)

	for _, profile := range profiles {
		fmt.Printf("%s %s %t\n", profile.name, profile.inherit_from, profile.default_switch)
	}

}

func buildProfiles(tokens []Token) []Profile {
	var profiles []Profile

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.token_type == PROFILE_NAME {
			profile := Profile{name: token.value, default_switch: false}
			if tokens[i+1].token_type == INHERITED_PROFILE_NAME {
				profile.inherit_from = tokens[i+1].value
			}
			if tokens[i+2].token_type == DEFAULT_SWITCH {
				profile.default_switch = true
			}
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j].token_type == PROFILE_NAME {
					break
				}
				variable := ProfileVariable{}
				if tokens[j].token_type == IDENTIFIER {
					variable.name = tokens[j].value
					if tokens[j+1].token_type == LITERAL {
						variable.value = tokens[j+1].value
						profile.variables = append(profile.variables, variable)
						profiles = append(profiles, profile)
					}
				}
			}
		}
	}

	return profiles
}

func parseIniFile(fileName string) []Token {
	var tokens []Token
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var line_no int16 = 1

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")

		if len(parts) == 2 {
			value := strings.TrimSpace(parts[0])
			token := Token{token_type: IDENTIFIER, value: value, line: line_no}
			tokens = append(tokens, token)

			value = strings.TrimSpace(parts[1])
			token = Token{token_type: LITERAL, value: value, line: line_no}
			tokens = append(tokens, token)

			continue
		}

		for pos := 0; pos < len(line); pos++ {

			ch := line[pos]

			if ch == '[' {
				token := Token{token_type: LEFT_BRA, value: string(ch), line: line_no}
				tokens = append(tokens, token)
				pos += parseRestOfProfile(line[1:], line_no, &tokens)
				continue
			}

			if ch == ']' {
				token := Token{token_type: RIGHT_BRA, value: string(ch), line: line_no}
				tokens = append(tokens, token)
				continue
			}

			if ch == '=' {
				token := Token{token_type: EQUAL, value: string(ch), line: line_no}
				tokens = append(tokens, token)
				continue
			}

			if ch == ':' {
				token := Token{token_type: COLUMN, value: string(ch), line: line_no}
				tokens = append(tokens, token)
				continue
			}
		}
		line_no += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return tokens
}

func parseRestOfProfile(line string, line_no int16, tokens *[]Token) int {
	parts := strings.Split(line, ":")

	if len(parts) == 1 {
		value := strings.TrimSuffix(parts[0], "]")
		value = strings.TrimSpace(value)
		token := Token{token_type: PROFILE_NAME, value: value, line: line_no}
		*tokens = append(*tokens, token)
		return len(value)
	}

	if len(parts) == 2 {
		value1 := strings.TrimSpace(parts[0])
		token := Token{token_type: PROFILE_NAME, value: value1, line: line_no}
		*tokens = append(*tokens, token)

		value2 := strings.TrimSuffix(parts[1], "]")
		value2 = strings.TrimSpace(value2)
		token = Token{token_type: INHERITED_PROFILE_NAME, value: value2, line: line_no}
		*tokens = append(*tokens, token)

		return len(value1) + len(value2) + 1
	}

	if len(parts) == 3 {
		value1 := parts[0]
		value1 = strings.TrimSpace(value1)
		token := Token{token_type: PROFILE_NAME, value: value1, line: line_no}
		*tokens = append(*tokens, token)

		value2 := parts[1]
		value2 = strings.TrimSpace(value2)
		token = Token{token_type: INHERITED_PROFILE_NAME, value: value2, line: line_no}
		*tokens = append(*tokens, token)

		value3 := strings.TrimSuffix(parts[2], "]")
		value3 = strings.TrimSpace(value3)
		token = Token{token_type: DEFAULT_SWITCH, value: value3, line: line_no}
		*tokens = append(*tokens, token)

		return len(value1) + len(value2) + len(value3) + 2
	}

	return 0
}
