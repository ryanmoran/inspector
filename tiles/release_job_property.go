package tiles

import "strings"

type ReleaseJobProperty struct {
	Name string
}

type ReleaseJobProperties []ReleaseJobProperty

func (rjps ReleaseJobProperties) Contains(name string) bool {
	for _, property := range rjps {
		if strings.HasPrefix(name, property.Name) {
			return true
		}
	}

	return false
}
