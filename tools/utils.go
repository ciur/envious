package tools

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func UseDefaultProfile(profiles []Profile) {
	found := FindDefaultProfile(profiles)
	if found == nil {
		fmt.Printf("Default profile not found\n")
		return
	}

	inherit_from := (*found).inherit_from
	if len(inherit_from) > 0 {
		parent := FindProfile(profiles, inherit_from)

		if parent == nil {
			fmt.Printf("Parent profile %s not found\n", inherit_from)
			return
		}

		for _, v := range (*parent).variables {
			fmt.Printf("export %s=%s\n", v.name, v.value)
		}
	}

	for _, v := range (*found).variables {
		fmt.Printf("export %s=%s\n", v.name, v.value)
	}
}

func UseProfile(profiles []Profile, name string) {
	found := FindProfile(profiles, name)
	if found == nil {
		fmt.Printf("Profile %s not found\n", name)
		return
	}

	inherit_from := (*found).inherit_from
	if len(inherit_from) > 0 {
		parent := FindProfile(profiles, inherit_from)

		if parent == nil {
			fmt.Printf("Parent profile %s not found\n", inherit_from)
			return
		}

		for _, v := range (*parent).variables {
			fmt.Printf("export %s=%s\n", v.name, v.value)
		}
	}

	for _, v := range (*found).variables {
		fmt.Printf("export %s=%s\n", v.name, v.value)
	}
}

func FindDefaultProfile(profiles []Profile) *Profile {
	for _, profile := range profiles {
		if profile.default_switch {
			return &profile
		}
	}

	return nil
}

func FindProfile(profiles []Profile, name string) *Profile {
	for _, profile := range profiles {
		if profile.name == name {
			return &profile
		}
	}

	return nil
}

func listProfileVariables(profile Profile) {
	for _, v := range profile.variables {
		fmt.Printf("%s = %s\n", v.name, v.value)
	}
}

func ListProfiles(profiles []Profile, detailed *bool) {
	for _, profile := range profiles {
		if profile.default_switch {
			fmt.Printf("%s*\n", profile.name)
		} else {
			fmt.Printf("%s\n", profile.name)
		}

		if *detailed {
			inherit_from := profile.inherit_from
			found := FindProfile(profiles, inherit_from)
			if found != nil {
				listProfileVariables(*found)
			}
			listProfileVariables(profile)
		}
	}
}

func BuildProfiles(tokens []Token) []Profile {
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

func ParseIniFile(fileName string) []Token {
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
