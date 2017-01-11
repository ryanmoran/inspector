package tiles

type ReleasePackage struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
}
