package tiles_test

import (
	"bytes"

	"github.com/ryanmoran/inspector/tiles"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		It("parses the tile contents", func() {
			stdout := bytes.NewBuffer([]byte{})
			parser := tiles.NewParser(pathToProduct, stdout)
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

			Expect(stdout.String()).To(ContainSubstring("parsing product metadata"))
			Expect(stdout.String()).To(ContainSubstring("parsing release: some-release"))
			Expect(stdout.String()).To(ContainSubstring("  - parsing job: some-job"))
		})
	})
})
