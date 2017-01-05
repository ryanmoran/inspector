package tiles

import "strings"

type MetadataJob struct {
	Name           string                `yaml:"name"`
	Templates      []MetadataJobTemplate `yaml:"templates"`
	Manifest       string                `yaml:"manifest"`
	ParsedManifest map[interface{}]interface{}
}

func (mj MetadataJob) UnusedManifestProperties(releases []Release) []string {
	var releaseJobProperties ReleaseJobProperties
	for _, release := range releases {
		for _, releaseJob := range release.Jobs {
			releaseJobProperties = append(releaseJobProperties, releaseJob.Properties...)
		}
	}

	var unusedManifestProperties []string
	for _, propertyName := range keysFromManifest(mj.ParsedManifest) {
		if !releaseJobProperties.Contains(propertyName) {
			unusedManifestProperties = append(unusedManifestProperties, propertyName)
		}
	}

	return unusedManifestProperties
}

func keysFromManifest(node map[interface{}]interface{}) []string {
	var keys []string
	for key, value := range node {
		if m, ok := value.(map[interface{}]interface{}); ok {
			for _, k := range keysFromManifest(m) {
				keys = append(keys, strings.Join([]string{key.(string), k}, "."))
			}
		} else {
			keys = append(keys, key.(string))
		}
	}

	return keys
}
