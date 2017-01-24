package commands

import (
	"fmt"
	"io"
	"sort"

	"github.com/ryanmoran/inspector/tiles"
)

type productParser interface {
	Parse() (tiles.Product, error)
}

type Deadweight struct {
	productParser productParser
	stdout        io.Writer
}

func NewDeadweight(productParser productParser, stdout io.Writer) Deadweight {
	return Deadweight{
		productParser: productParser,
		stdout:        stdout,
	}
}

func (d Deadweight) Execute(args []string) error {
	product, err := d.productParser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(d.stdout, "\n\nThe following job manifest properties are not being used by the included release templates:")
	for _, job := range product.Metadata.Jobs {
		unusedManifestProperties := job.UnusedManifestProperties(product.Releases)
		if len(unusedManifestProperties) > 0 {
			fmt.Fprintf(d.stdout, "Job: %s\n", job.Name)
			sort.Sort(unusedManifestProperties)
			for _, property := range unusedManifestProperties {
				fmt.Fprintf(d.stdout, "  - %s", property.Name)
				if property.ReferencesParsedManifest {
					fmt.Fprint(d.stdout, " (references parsed manifest)")
				}
				fmt.Fprint(d.stdout, "\n")
			}
		}
	}

	fmt.Fprintln(d.stdout, "\n\nThe following job templates are not being used:")
	jobTemplateUsageCounts := map[string]map[string]int{}
	for _, release := range product.Releases {
		jobTemplateUsageCounts[release.Name] = map[string]int{}
		for _, job := range release.Jobs {
			jobTemplateUsageCounts[release.Name][job.Name] = 0
		}
	}
	for _, job := range product.Metadata.Jobs {
		for _, template := range job.Templates {
			jobTemplateUsageCounts[template.Release][template.Name]++
		}
	}
	for _, release := range jobTemplateUsageCounts {
		for jobName, count := range release {
			if count != 0 {
				delete(release, jobName)
			}
		}
	}
	for releaseName, release := range jobTemplateUsageCounts {
		if len(release) > 0 {
			fmt.Fprintf(d.stdout, "Release: %s\n", releaseName)
			var jobs []string
			for jobName, _ := range release {
				jobs = append(jobs, jobName)
			}

			sort.Strings(jobs)

			for _, job := range jobs {
				fmt.Fprintf(d.stdout, "  - %s\n", job)
			}
		}
	}

	return nil
}

func (d Deadweight) Usage() Usage {
	return Usage{
		Description:      "something something dead",
		ShortDescription: "something dead",
	}
}
