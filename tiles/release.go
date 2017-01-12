package tiles

type Release struct {
	Name             string `yaml:"name"`
	Jobs             []ReleaseJob
	CompiledPackages []ReleasePackage `yaml:"compiled_packages"`
	Packages         []ReleasePackage `yaml:"packages"`
}
