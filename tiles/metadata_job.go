package tiles

import (
	"reflect"
	"strings"
)

type MetadataJob struct {
	Name           string                `yaml:"name"`
	Templates      []MetadataJobTemplate `yaml:"templates"`
	Manifest       string                `yaml:"manifest"`
	ParsedManifest map[interface{}]interface{}
}

type MetadataJobManifestProperty struct {
	Name    string
	Value   interface{}
	Link    string
	Job     string
	Release string

	ReferencesParsedManifest bool
	MirrorsDefault           bool
}

type MetadataJobManifestProperties []MetadataJobManifestProperty

func (mjmps MetadataJobManifestProperties) Len() int {
	return len(mjmps)
}

func (mjmps MetadataJobManifestProperties) Less(i, j int) bool {
	return mjmps[i].Name < mjmps[j].Name
}

func (mjmps MetadataJobManifestProperties) Swap(i, j int) {
	mjmps[i], mjmps[j] = mjmps[j], mjmps[i]
}

func (mj MetadataJob) UnusedManifestProperties(releases []Release) MetadataJobManifestProperties {
	var releaseJobProperties ReleaseJobProperties
	for _, template := range mj.Templates {
		for _, release := range releases {
			if template.Release == release.Name {
				for _, releaseJob := range release.Jobs {
					if template.Name == releaseJob.Name {
						releaseJobProperties = append(releaseJobProperties, releaseJob.Properties...)
					}
				}
			}
		}
	}

	var unusedManifestProperties []MetadataJobManifestProperty
	for _, property := range mj.Properties() {
		jobProperty, found := releaseJobProperties.Find(property.Name)
		if !found {
			unusedManifestProperties = append(unusedManifestProperties, property)
		} else {
			if reflect.DeepEqual(property.Value, jobProperty.Default) {
				property.MirrorsDefault = true
				unusedManifestProperties = append(unusedManifestProperties, property)
			}
		}
	}

	return unusedManifestProperties
}

func (mj MetadataJob) Properties() []MetadataJobManifestProperty {
	return propertiesFromManifest(mj.ParsedManifest)
}

func propertiesFromManifest(node map[interface{}]interface{}) []MetadataJobManifestProperty {
	var keys []MetadataJobManifestProperty

	for key, value := range node {
		switch m := value.(type) {
		case map[interface{}]interface{}:
			for _, k := range propertiesFromManifest(m) {
				keys = append(keys, MetadataJobManifestProperty{
					Name: strings.Join([]string{key.(string), k.Name}, "."),
					ReferencesParsedManifest: k.ReferencesParsedManifest,
					Value: k.Value,
				})
			}
		case string:
			keys = append(keys, MetadataJobManifestProperty{
				Name: key.(string),
				ReferencesParsedManifest: strings.Contains(m, ".parsed_manifest("),
				Value: m,
			})
		default:
			keys = append(keys, MetadataJobManifestProperty{
				Name:  key.(string),
				Value: m,
			})
		}
	}

	return keys
}

func (mj MetadataJob) LinkableProperties(releases []Release) MetadataJobManifestProperties {
	var releaseJobProperties ReleaseJobProperties
	for _, template := range mj.Templates {
		for _, release := range releases {
			if template.Release == release.Name {
				for _, releaseJob := range release.Jobs {
					if template.Name == releaseJob.Name {
						releaseJobProperties = append(releaseJobProperties, releaseJob.Properties...)
					}
				}
			}
		}
	}

	var linkableProperties []MetadataJobManifestProperty
	for _, property := range propertiesFromManifest(mj.ParsedManifest) {
		jobProperty, found := releaseJobProperties.Find(property.Name)
		if found && jobProperty.Link != "" {
			property.Link = jobProperty.Link
			property.Job = jobProperty.Job
			property.Release = jobProperty.Release
			linkableProperties = append(linkableProperties, property)
		}
	}

	return linkableProperties
}
