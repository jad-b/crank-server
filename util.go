package torque

import "strings"

// SlashJoin performs a strings.Join using '/' as a separator.
func SlashJoin(args ...string) string {
	return strings.Join(args, "/")
}
