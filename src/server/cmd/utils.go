package cmd

import (
	"regexp"
)

var CheckIdRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]{1,32}$`)

func IsValidId(id string) bool {
	return CheckIdRegex.MatchString(id)
}
