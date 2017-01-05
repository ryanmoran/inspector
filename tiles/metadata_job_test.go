package tiles_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ryanmoran/inspector/tiles"
)

var _ = Describe("MetadataJob", func() {
	Describe("UnusedManifestProperties", func() {
		It("returns a list of manifest properties that are not used by a job's templates", func() {
			metadataJob := tiles.MetadataJob{
				Templates: []tiles.MetadataJobTemplate{
					{
						Name:    "second-template",
						Release: "second-release",
					},
					{
						Name:    "fourth-template",
						Release: "fourth-release",
					},
				},
				ParsedManifest: map[interface{}]interface{}{
					"property": map[interface{}]interface{}{
						"first":  "one",
						"second": "two",
						"third":  "three",
						"fourth": "four",
					},
				},
			}

			releases := []tiles.Release{
				{
					Name: "first-release",
					Jobs: []tiles.ReleaseJob{
						{
							Name: "first-template",
							Properties: []tiles.ReleaseJobProperty{
								{Name: "property.not-used"},
							},
						},
					},
				},
				{
					Name: "second-release",
					Jobs: []tiles.ReleaseJob{
						{
							Name: "second-template",
							Properties: []tiles.ReleaseJobProperty{
								{Name: "property.second"},
							},
						},
					},
				},
				{
					Name: "fourth-release",
					Jobs: []tiles.ReleaseJob{
						{
							Name: "fourth-template",
							Properties: []tiles.ReleaseJobProperty{
								{Name: "property.fourth"},
							},
						},
					},
				},
			}

			Expect(metadataJob.UnusedManifestProperties(releases)).To(ConsistOf([]string{
				"property.first",
				"property.third",
			}))
		})
	})
})