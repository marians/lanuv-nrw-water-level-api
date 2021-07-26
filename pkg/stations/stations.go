package stations

import "strings"

func Normalize(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")

	return name
}
