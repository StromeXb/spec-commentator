package utils

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

const (
	Base      = 10
	BitSize32 = 32
	BitSize64 = 64
)

// Snake copy from entc/gen/func.go
func Snake(s string) string {
	var (
		j int
		b strings.Builder
	)
	for i := 0; i < len(s); i++ {
		r := rune(s[i])
		// Put '_' if it is not a start or end of a word, current letter is uppercase,
		// and previous is lowercase (cases like: "UserInfo"), or next letter is also
		// a lowercase and previous letter is not "_".
		if i > 0 && i < len(s)-1 && unicode.IsUpper(r) {
			if unicode.IsLower(rune(s[i-1])) ||
				j != i-1 && unicode.IsLower(rune(s[i+1])) && unicode.IsLetter(rune(s[i-1])) {
				j = i
				b.WriteString("_")
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func UUIDEqualed(firstID, secondID string) bool {
	apiUserUUID, err := uuid.Parse(firstID)
	if err != nil {
		return false
	}
	headerUserUUID, err := uuid.Parse(secondID)
	if err != nil {
		return false
	}
	return apiUserUUID == headerUserUUID
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

type ServerEnvironment string

const (
	DEV   ServerEnvironment = "dev"
	STAGE ServerEnvironment = "stage"
	PROD  ServerEnvironment = "prod"
)
