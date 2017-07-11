package commands

import (
	"fmt"
	"io"

	"github.com/ryanmoran/inspector/tiles"
)

type Links struct {
	productParser productParser
	stdout        io.Writer
}

func NewLinks(productParser productParser, stdout io.Writer) Links {
	return Links{
		productParser: productParser,
		stdout:        stdout,
	}
}

func (l Links) Execute(args []string) error {
	product, err := l.productParser.Parse()
	if err != nil {
		panic(err)
	}

	var templates []tiles.MetadataJobTemplate
	for _, job := range product.Metadata.Jobs {
		for _, template := range job.Templates {
			templates = append(templates, template)
		}
	}

	for _, release := range product.Releases {
		for _, job := range release.Jobs {
			for _, template := range templates {
				if template.Name == job.Name && template.Release == release.Name {
					for _, link := range job.Provides {
						fmt.Fprintf(l.stdout, "provides: %s %s %s (%s)\n", release.Name, job.Name, link.Name, link.Type)
					}
					for _, link := range job.Consumes {
						fmt.Fprintf(l.stdout, "consumes: %s %s %s (%s)\n", release.Name, job.Name, link.Name, link.Type)
					}
				}
			}
		}
	}

	return nil
}

func (l Links) Usage() Usage {
	return Usage{
		Description:      "something something links",
		ShortDescription: "something links",
	}
}
