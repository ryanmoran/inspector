package tiles

import "strings"

type MetadataJob struct {
	Name           string                `yaml:"name"`
	Templates      []MetadataJobTemplate `yaml:"templates"`
	Manifest       string                `yaml:"manifest"`
	ParsedManifest map[interface{}]interface{}
}

type MetadataJobManifestProperty struct {
	Name                     string
	ReferencesParsedManifest bool
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
	for _, property := range propertiesFromManifest(mj.ParsedManifest) {
		if !releaseJobProperties.Contains(property.Name) {
			unusedManifestProperties = append(unusedManifestProperties, property)
		}
	}

	return unusedManifestProperties
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
				})
			}
		case string:
			keys = append(keys, MetadataJobManifestProperty{
				Name: key.(string),
				ReferencesParsedManifest: strings.Contains(m, ".parsed_manifest("),
			})
		default:
			keys = append(keys, MetadataJobManifestProperty{
				Name: key.(string),
			})
		}
	}

	return keys
}
