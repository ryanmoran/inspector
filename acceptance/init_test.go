package acceptance

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	pathToMain string
	pathToTile string
)

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "acceptance")
}

var _ = BeforeSuite(func() {
	var err error
	pathToMain, err = gexec.Build("github.com/ryanmoran/inspector")
	Expect(err).NotTo(HaveOccurred())

	pathToTile, err = buildTile()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(pathToTile)
	Expect(err).NotTo(HaveOccurred())

	gexec.CleanupBuildArtifacts()
})

func buildTile() (string, error) {
	metadataFile := bytes.NewBufferString(`---
job_types:
- name: diego_brain
  manifest: |
    capi:
      nsync:
        cc:
          staging_upload_user: staging-upload-user
          staging_upload_password: staging-upload-password
      stager:
        cc:
          staging_upload_user: staging-upload-user
          staging_upload_password: staging-upload-password
    nats:
      port: 1234
      machines: ["1.2.3.4"]
      user: nats-user
      password: nats-password
    parsed:
      manifest: (( .properties.references.parsed_manifest(example) ))
`)

	releaseManifestFile := bytes.NewBufferString(`---
name: diego
`)

	releaseJobManifestFile := bytes.NewBufferString(`---
name: nsync
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
