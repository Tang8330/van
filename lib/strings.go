package lib

import "strings"

func SpaceToDash(s string) string {
	return strings.Replace(s, " ", "-", -1)
}
