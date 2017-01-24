package tiles

type ReleasePackage struct {
	Name         string   `yaml:"name"`
	Dependencies []string `yaml:"dependencies"`
}

func (rp ReleasePackage) AllPackages(release Release) []string {
	var packages []string
	for _, dependency := range rp.Dependencies {
		for _, pkg := range append(release.Packages, release.CompiledPackages...) {
			if dependency == pkg.Name {
				packages = append(packages, pkg.Name)
				packages = append(packages, pkg.AllPackages(release)...)
			}
		}
	}

	return packages
}
