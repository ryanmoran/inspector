package tiles

type Release struct {
	Name     string `yaml:"name"`
	Jobs     []ReleaseJob
	Packages []ReleasePackage `yaml:"compiled_packages"`
}
