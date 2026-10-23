package cmd

import (
	"regexp"
	pr "treenity/src/protocol"
)

var CheckIdRegex = regexp.MustCompile(pr.ValidRegex)

func IsValidId(id string) bool {
	return CheckIdRegex.MatchString(id)
}
