package commands

import (
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/ryanmoran/inspector/flags"
)

type pkgDepMatch struct {
	Release  string
	Packages []pkgDepMatchPackage
}

type pkgDepMatchPackage struct {
	Name         string
	Dependencies []string
}

type PkgDep struct {
	productParser productParser
	stdout        io.Writer
	Options       struct {
		Match string `short:"m"  long:"match"  description:"text to match in package dependency"`
	}
}

func NewPkgDep(productParser productParser, stdout io.Writer) PkgDep {
	return PkgDep{
		productParser: productParser,
		stdout:        stdout,
	}
}

func (pd PkgDep) Execute(args []string) error {
	_, err := flags.Parse(&pd.Options, args)
	if err != nil {
		return err
	}

	if pd.Options.Match == "" {
		return errors.New("-match is a required flag")
	}

	matchRegexp, err := regexp.Compile(pd.Options.Match)
	if err != nil {
		panic(err)
	}

	product, err := pd.productParser.Parse()
	if err != nil {
		panic(err)
	}

	var matches []pkgDepMatch
	for _, release := range product.Releases {
		var packages []pkgDepMatchPackage
		for _, pkg := range release.CompiledPackages {
			var dependencies []string
			for _, dependency := range pkg.Dependencies {
				if matchRegexp.MatchString(dependency) {
					dependencies = append(dependencies, dependency)
				}
			}

			if len(dependencies) > 0 {
				packages = append(packages, pkgDepMatchPackage{
					Name:         pkg.Name,
					Dependencies: dependencies,
				})
			}
		}

		for _, pkg := range release.Packages {
			var dependencies []string
			for _, dependency := range pkg.Dependencies {
				if matchRegexp.MatchString(dependency) {
					dependencies = append(dependencies, dependency)
				}
			}

			if len(dependencies) > 0 {
				packages = append(packages, pkgDepMatchPackage{
					Name:         pkg.Name,
					Dependencies: dependencies,
				})
			}
		}

		if len(packages) > 0 {
			matches = append(matches, pkgDepMatch{
				Release:  release.Name,
				Packages: packages,
			})
		}
	}

	fmt.Fprintf(pd.stdout, "\n\nThe following packages have a dependency that matches %q:\n", pd.Options.Match)
	for _, m := range matches {
		fmt.Fprintf(pd.stdout, "Release: %s\n", m.Release)
		for _, pkg := range m.Packages {
			fmt.Fprintf(pd.stdout, "  - %s %v\n", pkg.Name, pkg.Dependencies)
		}
	}

	return nil
}

func (pd PkgDep) Usage() Usage {
	return Usage{
		Description:      "something something dep",
		ShortDescription: "something dep",
	}
}
