package commands

import (
	"fmt"
	"strings"
)

type BrokenReleaseRefs struct {
	productParser productParser
}

func NewBrokenReleaseRefs(productParser productParser) BrokenReleaseRefs {
	return BrokenReleaseRefs{
		productParser: productParser,
	}
}

func (b BrokenReleaseRefs) Execute(args []string) error {
	product, err := b.productParser.Parse()
	if err != nil {
		panic(err)
	}

	releases := map[string]bool{}
	for _, release := range product.Releases {
		releases[release.Name] = true
	}

	errMsgs := []string{}
	for _, job := range product.Metadata.Jobs {
		for _, template := range job.Templates {
			if _, ok := releases[template.Release]; !ok {
				errMsgs = append(errMsgs, fmt.Sprintf("%s (referenced by template %q in %q job)", template.Release, template.Name, job.Name))
			}
		}
	}

	if len(errMsgs) > 0 {
		return fmt.Errorf("The following releases are referenced but not in the tile:\n%s\n", strings.Join(errMsgs, "\n"))
	}

	return nil
}

func (b BrokenReleaseRefs) Usage() Usage {
	return Usage{
		Description:      "prints releases that are referenced by job templates but not in tile",
		ShortDescription: "prints missing releases",
	}
}
