package commands

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/ryanmoran/inspector/flags"
)

type Grep struct {
	productParser productParser
	Options       struct {
		Match string `short:"m"  long:"match"  description:"text to match"`
	}
}

func NewGrep(productParser productParser) Grep {
	return Grep{
		productParser: productParser,
	}
}

func (g Grep) Execute(args []string) error {
	_, err := flags.Parse(&g.Options, args)
	if err != nil {
		return err
	}

	if g.Options.Match == "" {
		return errors.New("-match is a required flag")
	}

	matchRegexp, err := regexp.Compile(g.Options.Match)
	if err != nil {
		panic(err)
	}

	product, err := g.productParser.Parse()
	if err != nil {
		panic(err)
	}

	for _, metadataJob := range product.Metadata.Jobs {
		for _, metadataJobProperty := range metadataJob.Properties() {
			if matchRegexp.MatchString(metadataJobProperty.Name) ||
				matchRegexp.MatchString(fmt.Sprintf("%+v", metadataJobProperty.Value)) {
				fmt.Printf("metadata ref: %s: %s: %+v\n", metadataJob.Name, metadataJobProperty.Name, metadataJobProperty.Value)
			}
		}
	}

	for _, release := range product.Releases {
		for _, releaseJob := range release.Jobs {
			for _, releaseJobProperty := range releaseJob.Properties {
				if matchRegexp.MatchString(releaseJobProperty.Name) ||
					matchRegexp.MatchString(fmt.Sprintf("%+v", releaseJobProperty.Default)) {
					fmt.Printf("release ref: %s: %s: %s: %+v\n", release.Name, releaseJob.Name, releaseJobProperty.Name, releaseJobProperty.Default)
				}
			}
		}
	}

	return nil
}

func (g Grep) Usage() Usage {
	return Usage{}
}
