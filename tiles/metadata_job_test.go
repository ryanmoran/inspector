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
						"first": "one",
						"second": map[interface{}]interface{}{
							"other": "fields",
						},
						"third":    "(( .properties.references.parsed_manifest(three) ))",
						"fourth":   "four",
						"not-used": "bad",
						"default":  "default",
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
								{Name: "property.default", Default: "default"},
							},
						},
					},
				},
			}

			Expect(metadataJob.UnusedManifestProperties(releases)).To(ConsistOf([]tiles.MetadataJobManifestProperty{
				{
					Name:  "property.first",
					Value: "one",
				},
				{
					Name:  "property.third",
					Value: "(( .properties.references.parsed_manifest(three) ))",
					ReferencesParsedManifest: true,
				},
				{
					Name:  "property.not-used",
					Value: "bad",
				},
				{
					Name:           "property.default",
					Value:          "default",
					MirrorsDefault: true,
				},
			}))
		})
	})

	Describe("LinkableProperties", func() {
		It("returns a list of properties that can be referenced by link", func() {
			metadataJob := tiles.MetadataJob{
				Templates: []tiles.MetadataJobTemplate{
					{
						Name:    "some-job",
						Release: "some-release",
					},
				},
				ParsedManifest: map[interface{}]interface{}{
					"property": map[interface{}]interface{}{
						"first": "one",
					},
					"link": map[interface{}]interface{}{
						"first": "banana",
					},
				},
			}

			releases := []tiles.Release{
				{
					Name: "some-release",
					Jobs: []tiles.ReleaseJob{
						{
							Name: "some-job",
							Properties: []tiles.ReleaseJobProperty{
								{Name: "link.first", Link: "some-link", Job: "some-job", Release: "some-release"},
							},
							Provides: []tiles.ReleaseJobProvideLink{
								{
									Name:       "some-link",
									Type:       "some-link-type",
									Properties: []string{"link.first"},
								},
							},
						},
					},
				},
			}

			Expect(metadataJob.LinkableProperties(releases)).To(ConsistOf([]tiles.MetadataJobManifestProperty{
				{
					Name:    "link.first",
					Link:    "some-link",
					Job:     "some-job",
					Release: "some-release",
					Value:   "banana",
				},
			}))
		})
	})
})
