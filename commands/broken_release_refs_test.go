package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/inspector/commands"
	"github.com/ryanmoran/inspector/commands/fakes"
	"github.com/ryanmoran/inspector/tiles"
)

var _ = Describe("Broken-Release-Refs", func() {
	Describe("Execute", func() {
		It("finds references to releases that are not in the tile", func() {
			productParser := &fakes.ProductParser{}
			productParser.ParseCall.Returns.Product = tiles.Product{
				Metadata: tiles.Metadata{
					Jobs: []tiles.MetadataJob{
						{
							Name: "some_diego_job",
							Templates: []tiles.MetadataJobTemplate{
								{
									Name:    "some-diego-job-template",
									Release: "some-other-diego-release",
								},
								{
									Name:    "some-capi-job-template-1",
									Release: "some-other-capi-release",
								},
							},
						},
						{
							Name: "some_capi_job",
							Templates: []tiles.MetadataJobTemplate{
								{
									Name:    "some-capi-job-template-2",
									Release: "some-other-capi-release",
								},
							},
						},
					},
				},
			}

			command := commands.NewBrokenReleaseRefs(productParser)

			err := command.Execute([]string{})

			Expect(productParser.ParseCall.CallCount).To(Equal(1))
			Expect(err).To(MatchError(`The following releases are referenced but not in the tile:
some-other-diego-release (referenced by template "some-diego-job-template" in "some_diego_job" job)
some-other-capi-release (referenced by template "some-capi-job-template-1" in "some_diego_job" job)
some-other-capi-release (referenced by template "some-capi-job-template-2" in "some_capi_job" job)
`))
		})
	})

	Describe("Usage", func() {
		It("returns a descriptive usage", func() {
			command := commands.NewBrokenReleaseRefs(nil)
			Expect(command.Usage()).To(Equal(commands.Usage{
				Description:      "prints releases that are referenced by job templates but not in tile",
				ShortDescription: "prints missing releases",
			}))
		})
	})
})
