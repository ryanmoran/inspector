package tiles_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/inspector/tiles"
)

var _ = Describe("Product", func() {
	Describe("UnusedReleaseJobs", func() {
		Context("when job refernces a releae that is not in the tile", func() {
			It("returns an error", func() {
				p := tiles.Product{
					Metadata: tiles.Metadata{
						Jobs: []tiles.MetadataJob{
							{
								Name: "some-job",
								Templates: []tiles.MetadataJobTemplate{
									{
										Name:    "some-template",
										Release: "some-release",
									},
								},
							},
						},
					},
					Releases: []tiles.Release{
						{
							Name: "some-other-release",
						},
					},
				}

				_, err := p.UnusedReleaseJobs()
				Expect(err).To(MatchError(`"some-release" is not in the tile (referenced by template "some-template" in job "some-job")`))
			})
		})
	})
})
