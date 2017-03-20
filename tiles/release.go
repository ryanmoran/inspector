package tiles

type Release struct {
	Name             string           `yaml:"name"`
	Jobs             []ReleaseJob     `yaml:"-"`
	CompiledPackages []ReleasePackage `yaml:"compiled_packages"`
	Packages         []ReleasePackage `yaml:"packages"`
}
