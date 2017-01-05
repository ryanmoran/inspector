package tiles_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTiles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tiles")
}

var pathToProduct string

var _ = BeforeSuite(func() {
	var err error
	pathToProduct, err = generateProduct()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(pathToProduct)
	Expect(err).NotTo(HaveOccurred())
})

func generateProduct() (string, error) {
	metadataFile := bytes.NewBufferString(`---
job_types:
- name: some-job
  manifest: |
    property:
      first: one
`)

	releaseManifestFile := bytes.NewBufferString(`---
name: some-release
`)

	releaseJobManifestFile := bytes.NewBufferString(`---
name: some-job
properties:
  some-property: {}
`)

	releaseJobFile := bytes.NewBuffer([]byte{})
	gzipReleaseJobFile := gzip.NewWriter(releaseJobFile)
	tarGzipReleaseJobFile := tar.NewWriter(gzipReleaseJobFile)
	err := tarGzipReleaseJobFile.WriteHeader(&tar.Header{
		Name: "job.MF",
		Size: int64(releaseJobManifestFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseJobFile, releaseJobManifestFile)
	if err != nil {
		return "", err
	}

	err = tarGzipReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	err = gzipReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	releaseFile := bytes.NewBuffer([]byte{})
	gzipReleaseFile := gzip.NewWriter(releaseFile)
	tarGzipReleaseFile := tar.NewWriter(gzipReleaseFile)
	err = tarGzipReleaseFile.WriteHeader(&tar.Header{
		Name: "release.MF",
		Size: int64(releaseManifestFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseFile, releaseManifestFile)
	if err != nil {
		return "", err
	}

	err = tarGzipReleaseFile.WriteHeader(&tar.Header{
		Name: "jobs/some_job.tgz",
		Size: int64(releaseJobFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseFile, releaseJobFile)
	if err != nil {
		return "", err
	}

	err = tarGzipReleaseFile.Close()
	if err != nil {
		return "", nil
	}

	err = gzipReleaseFile.Close()
	if err != nil {
		return "", nil
	}

	productFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}

	zipFile := zip.NewWriter(productFile)

	zipMetadataFile, err := zipFile.Create("metadata/banana.yml")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(zipMetadataFile, metadataFile)
	if err != nil {
		return "", err
	}

	zipReleaseFile, err := zipFile.Create("compiled_releases/some-release-1.2.3.tgz")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(zipReleaseFile, releaseFile)
	if err != nil {
		return "", err
	}

	err = zipFile.Close()
	if err != nil {
		return "", err
	}

	err = productFile.Close()
	if err != nil {
		return "", err
	}

	return productFile.Name(), nil
}
