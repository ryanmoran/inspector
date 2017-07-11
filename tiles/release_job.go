package tiles

type ReleaseJob struct {
	Name       string
	Properties []ReleaseJobProperty
	Packages   []string
	Provides   []ReleaseJobProvideLink
	Consumes   []ReleaseJobConsumeLink
}

func (rj ReleaseJob) AllPackages(release Release) []string {
	uniquePackages := map[string]struct{}{}

	var packages []ReleasePackage
	for _, packageName := range rj.Packages {
		uniquePackages[packageName] = struct{}{}

		for _, releasePackage := range append(release.Packages, release.CompiledPackages...) {
			if releasePackage.Name == packageName {
				packages = append(packages, releasePackage)
			}
		}
	}

	for _, releasePackage := range packages {
		for _, packageName := range releasePackage.AllPackages(release) {
			uniquePackages[packageName] = struct{}{}
		}
	}

	var packageNames []string

	for name, _ := range uniquePackages {
		packageNames = append(packageNames, name)
	}

	return packageNames
}
