package tiles

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	yaml "gopkg.in/yaml.v2"
)

var (
	metadataRegexp           = regexp.MustCompile(`metadata\/.*\.yml$`)
	releaseRegexp            = regexp.MustCompile(`releases\/.*\.tgz$`)
	releaseManifestRegexp    = regexp.MustCompile(`release\.MF$`)
	releaseJobRegexp         = regexp.MustCompile(`jobs\/.*\.tgz$`)
	releaseJobManifestRegexp = regexp.MustCompile(`job\.MF$`)
)

type Parser struct {
	path string
}

func NewParser(path string) Parser {
	return Parser{
		path: path,
	}
}

func (p Parser) Parse() (Product, error) {
	var productManifest Metadata
	var productReleases []Release

	productFile, err := os.Open(p.path)
	if err != nil {
		panic(err)
	}

	productFileInfo, err := productFile.Stat()
	if err != nil {
		panic(err)
	}

	productFileZip, err := zip.NewReader(productFile, productFileInfo.Size())
	if err != nil {
		panic(err)
	}

	for _, f := range productFileZip.File {
		if metadataRegexp.MatchString(f.FileHeader.Name) {
			metadataFile, err := f.Open()
			if err != nil {
				panic(err)
			}

			metadataFileContents, err := ioutil.ReadAll(metadataFile)
			if err != nil {
				panic(err)
			}

			err = yaml.Unmarshal(metadataFileContents, &productManifest)
			if err != nil {
				panic(err)
			}

			for i, job := range productManifest.Jobs {
				var parsedManifest map[interface{}]interface{}
				err = yaml.Unmarshal([]byte(job.Manifest), &parsedManifest)
				if err != nil {
					panic(err)
				}
				productManifest.Jobs[i].ParsedManifest = parsedManifest
			}
		}

		if releaseRegexp.MatchString(f.FileHeader.Name) {
			var releaseManifest Release

			releaseFile, err := f.Open()
			if err != nil {
				panic(err)
			}

			gzipReleaseFile, err := gzip.NewReader(releaseFile)
			if err != nil {
				panic(err)
			}

			tarGzipReleaseFile := tar.NewReader(gzipReleaseFile)

			header, err := tarGzipReleaseFile.Next()
			for err == nil {
				if releaseManifestRegexp.MatchString(header.Name) {
					releaseManifestContents, err := ioutil.ReadAll(tarGzipReleaseFile)
					if err != nil {
						panic(err)
					}

					err = yaml.Unmarshal(releaseManifestContents, &releaseManifest)
					if err != nil {
						panic(err)
					}
				}

				if releaseJobRegexp.MatchString(header.Name) {
					gzipReleaseJobFile, err := gzip.NewReader(tarGzipReleaseFile)
					if err != nil {
						panic(err)
					}

					tarGzipReleaseJobFile := tar.NewReader(gzipReleaseJobFile)

					header, err := tarGzipReleaseJobFile.Next()
					for err == nil {
						if releaseJobManifestRegexp.MatchString(header.Name) {
							releaseJobManifestContents, err := ioutil.ReadAll(tarGzipReleaseJobFile)
							if err != nil {
								panic(err)
							}

							var releaseJobManifest struct {
								Name       string              `yaml:"name"`
								Properties map[string]struct{} `yaml:"properties"`
							}
							err = yaml.Unmarshal(releaseJobManifestContents, &releaseJobManifest)
							if err != nil {
								panic(err)
							}

							var releaseJobProperties []ReleaseJobProperty
							for propertyName, _ := range releaseJobManifest.Properties {
								releaseJobProperties = append(releaseJobProperties, ReleaseJobProperty{
									Name: propertyName,
								})
							}

							releaseManifest.Jobs = append(releaseManifest.Jobs, ReleaseJob{
								Name:       releaseJobManifest.Name,
								Properties: releaseJobProperties,
							})
						}

						header, err = tarGzipReleaseJobFile.Next()
					}
					if err != io.EOF {
						panic(err)
					}
				}

				header, err = tarGzipReleaseFile.Next()
			}
			if err != io.EOF {
				panic(err)
			}

			productReleases = append(productReleases, releaseManifest)
		}
	}

	return Product{
		Metadata: productManifest,
		Releases: productReleases,
	}, nil
}