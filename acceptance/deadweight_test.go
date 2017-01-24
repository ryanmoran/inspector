package acceptance

import (
	"os/exec"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("deadweight", func() {
	It("prints out a list of metadata manifest properties that are no longer used by their respective release jobs", func() {
		command := exec.Command(pathToMain, "--path", pathToTile, "deadweight")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(`The following job manifest properties are not being used by the included release templates:
Job: diego_brain
  - capi.nsync.cc.staging_upload_password
  - capi.nsync.cc.staging_upload_user
  - capi.stager.cc.staging_upload_password
  - capi.stager.cc.staging_upload_user
  - nats.machines
  - nats.password
  - nats.port
  - nats.user
  - parsed.manifest (references parsed manifest)


The following release jobs are not being used:
Release: diego
  - nsync
  - rep
`))
	})
})
