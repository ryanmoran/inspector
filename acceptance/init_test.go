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
  templates:
  - name: auctioneer
    release: diego
  manifest: |
    link-property: banana
    capi:
      nsync:
        cc:
          staging_upload_user: staging-upload-user
          staging_upload_password: staging-upload-password
      stager:
        cc:
          staging_upload_user: staging-upload-user
          staging_upload_password: staging-upload-password
          property_default: default-value
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
compiled_packages:
- name: acceptance-tests
  dependencies:
  - cli
  - golang1.5
packages:
- name: auctioneer-pkg
  dependencies:
  - golang1.6
  - 1.7golang
- name: golang1.6
- name: 1.7golang
- name: golang1.5
jobs:
- name: auctioneer
- name: nsync
- name: rep
`)

	auctioneerReleaseJobManifestFile := bytes.NewBufferString(`---
name: auctioneer
packages:
- auctioneer-pkg
- golang1.6
- 1.7golang
provides:
- name: auctioneer
  type: auctioneer
  properties:
  - link-property
consumes:
- name: auctioneer
  type: auctioneer
  optional: true
properties:
  some-property: {}
  link-property: {}
  capi.stager.cc.property_default:
    default: default-value
`)

	nsyncReleaseJobManifestFile := bytes.NewBufferString(`---
name: nsync
packages:
- nsync-pkg
properties:
  capi.nsync.cc.staging_upload_password: {}
  capi.nsync.cc.staging_upload_user: {}
`)

	repReleaseJobManifestFile := bytes.NewBufferString(`---
name: rep
packages:
- rep-pkg
properties: {}
`)

	// create auctioneer job
	auctioneerReleaseJobFile := bytes.NewBuffer([]byte{})
	gzipAuctioneerReleaseJobFile := gzip.NewWriter(auctioneerReleaseJobFile)
	tarGzipAuctioneerReleaseJobFile := tar.NewWriter(gzipAuctioneerReleaseJobFile)
	err := tarGzipAuctioneerReleaseJobFile.WriteHeader(&tar.Header{
		Name: "job.MF",
		Size: int64(auctioneerReleaseJobManifestFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipAuctioneerReleaseJobFile, auctioneerReleaseJobManifestFile)
	if err != nil {
		return "", err
	}

	err = tarGzipAuctioneerReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	err = gzipAuctioneerReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	// create nsync job
	nsyncReleaseJobFile := bytes.NewBuffer([]byte{})
	gzipNsyncReleaseJobFile := gzip.NewWriter(nsyncReleaseJobFile)
	tarGzipNsyncReleaseJobFile := tar.NewWriter(gzipNsyncReleaseJobFile)
	err = tarGzipNsyncReleaseJobFile.WriteHeader(&tar.Header{
		Name: "job.MF",
		Size: int64(nsyncReleaseJobManifestFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipNsyncReleaseJobFile, nsyncReleaseJobManifestFile)
	if err != nil {
		return "", err
	}

	err = tarGzipNsyncReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	err = gzipNsyncReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	// create rep job
	repReleaseJobFile := bytes.NewBuffer([]byte{})
	gzipRepReleaseJobFile := gzip.NewWriter(repReleaseJobFile)
	tarGzipRepReleaseJobFile := tar.NewWriter(gzipRepReleaseJobFile)
	err = tarGzipRepReleaseJobFile.WriteHeader(&tar.Header{
		Name: "job.MF",
		Size: int64(repReleaseJobManifestFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipRepReleaseJobFile, repReleaseJobManifestFile)
	if err != nil {
		return "", err
	}

	err = tarGzipRepReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	err = gzipRepReleaseJobFile.Close()
	if err != nil {
		return "", err
	}

	// create release
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
		Name: "jobs/auctioneer.tgz",
		Size: int64(auctioneerReleaseJobFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseFile, auctioneerReleaseJobFile)
	if err != nil {
		return "", err
	}

	err = tarGzipReleaseFile.WriteHeader(&tar.Header{
		Name: "jobs/nsync.tgz",
		Size: int64(nsyncReleaseJobFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseFile, nsyncReleaseJobFile)
	if err != nil {
		return "", err
	}

	err = tarGzipReleaseFile.WriteHeader(&tar.Header{
		Name: "jobs/rep.tgz",
		Size: int64(repReleaseJobFile.Len()),
	})
	if err != nil {
		return "", err
	}

	_, err = io.Copy(tarGzipReleaseFile, repReleaseJobFile)
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

	zipReleaseFile, err := zipFile.Create("compiled_releases/diego-1.2.3.tgz")
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
