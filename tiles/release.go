package tiles

type Release struct {
	Name string `yaml:"name"`
	Jobs []ReleaseJob
}
