package cmd

import (
	"regexp"
	pr "treenity/src/protocol"
)

var CheckIdRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

func IsValidId(id string) bool {
	return CheckIdRegex.MatchString(id) && len(id) <= int(pr.MaxIdCharacters)
}
