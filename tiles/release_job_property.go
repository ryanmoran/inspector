package tiles

type ReleaseJobProperty struct {
	Name string
}

type ReleaseJobProperties []ReleaseJobProperty

func (rjps ReleaseJobProperties) Contains(name string) bool {
	for _, property := range rjps {
		if property.Name == name {
			return true
		}
	}

	return false
}
