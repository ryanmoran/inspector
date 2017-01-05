package tiles_test

import (
	"github.com/ryanmoran/inspector/tiles"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		It("parses the tile contents", func() {
			parser := tiles.NewParser(pathToProduct)
			product, err := parser.Parse()
			Expect(err).NotTo(HaveOccurred())

			Expect(product).To(Equal(tiles.Product{
				Metadata: tiles.Metadata{
					Jobs: []tiles.MetadataJob{
						{
							Name:     "some-job",
							Manifest: "property:\n  first: one\n",
							ParsedManifest: map[interface{}]interface{}{
								"property": map[interface{}]interface{}{
									"first": "one",
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
								Name: "some-job",
								Properties: []tiles.ReleaseJobProperty{
									{Name: "some-property"},
								},
							},
						},
					},
				},
			}))
		})
	})
})
