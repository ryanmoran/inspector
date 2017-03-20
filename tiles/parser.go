package tiles

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
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
	path   string
	stdout io.Writer
}

func NewParser(path string, stdout io.Writer) Parser {
	return Parser{
		path:   path,
		stdout: stdout,
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
			fmt.Fprintln(p.stdout, "parsing product metadata")

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

					fmt.Fprintf(p.stdout, "parsing release: %s\n", releaseManifest.Name)
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
								Name       string `yaml:"name"`
								Properties map[string]struct {
									Default interface{} `yaml:"default"`
								} `yaml:"properties"`
								Packages []string `yaml:"packages"`
								Provides []struct {
									Name       string   `yaml:"name"`
									Properties []string `yaml:"properties"`
								} `yaml:"provides"`
							}
							err = yaml.Unmarshal(releaseJobManifestContents, &releaseJobManifest)
							if err != nil {
								panic(err)
							}

							fmt.Fprintf(p.stdout, "  - parsing job: %s\n", releaseJobManifest.Name)

							var releaseJobProperties []ReleaseJobProperty
							for propertyName, propertySpec := range releaseJobManifest.Properties {
								var linkName string
								for _, link := range releaseJobManifest.Provides {
									for _, linkProperty := range link.Properties {
										if propertyName == linkProperty {
											linkName = link.Name
										}
									}
								}

								releaseJobProperties = append(releaseJobProperties, ReleaseJobProperty{
									Name:    propertyName,
									Default: propertySpec.Default,
									Link:    linkName,
									Job:     releaseJobManifest.Name,
									Release: releaseManifest.Name,
								})
							}

							releaseManifest.Jobs = append(releaseManifest.Jobs, ReleaseJob{
								Name:       releaseJobManifest.Name,
								Properties: releaseJobProperties,
								Packages:   releaseJobManifest.Packages,
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

	product := Product{
		Metadata: productManifest,
		Releases: productReleases,
	}

	return product, nil
}
