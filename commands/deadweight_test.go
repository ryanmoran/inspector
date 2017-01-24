package commands_test

import (
	"bytes"

	"github.com/ryanmoran/inspector/commands"
	"github.com/ryanmoran/inspector/commands/fakes"
	"github.com/ryanmoran/inspector/tiles"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deadweight", func() {
	Describe("Execute", func() {
		It("finds the manifest properties that are unused and reports them", func() {
			productParser := &fakes.ProductParser{}
			productParser.ParseCall.Returns.Product = tiles.Product{
				Metadata: tiles.Metadata{
					Jobs: []tiles.MetadataJob{
						{
							Name: "some-job",
							Templates: []tiles.MetadataJobTemplate{
								{
									Name:    "some-job-template-1",
									Release: "some-release",
								},
							},
							ParsedManifest: map[interface{}]interface{}{
								"property": map[interface{}]interface{}{
									"first":  "one",
									"second": "two",
									"fourth": "(( .properties.references.parsed_manifest(four) ))",
								},
							},
						},
					},
				},
				Releases: []tiles.Release{
					{
						Name: "some-release",
						Jobs: []tiles.ReleaseJob{
							{
								Name: "some-job-template-1",
								Properties: []tiles.ReleaseJobProperty{
									{Name: "property.first"},
									{Name: "property.third"},
								},
							},
							{
								Name: "some-job-template-2",
							},
							{
								Name: "some-job-template-3",
							},
						},
					},
				},
			}

			stdout := bytes.NewBuffer([]byte{})
			command := commands.NewDeadweight(productParser, stdout)

			err := command.Execute([]string{})
			Expect(err).NotTo(HaveOccurred())

			Expect(stdout.String()).To(ContainSubstring(`Job: some-job
  - property.fourth (references parsed manifest)
  - property.second`))

			Expect(stdout.String()).To(ContainSubstring(`Release: some-release
  - some-job-template-2
  - some-job-template-3`))
		})
	})

	Describe("Usage", func() {
		It("returns a descriptive usage", func() {
			command := commands.NewDeadweight(nil, nil)
			Expect(command.Usage()).To(Equal(commands.Usage{
				Description:      "something something dead",
				ShortDescription: "something dead",
			}))
		})
	})
})
