package acceptance

import (
	"os/exec"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("broken-release-refs", func() {
	It("prints a list of releases that are referenced but not in the tile", func() {
		command := exec.Command(pathToMain, "--path", pathToBrokenReleaseRefsTile, "broken-release-refs")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(1))
		Expect(string(session.Err.Contents())).To(ContainSubstring(`The following releases are referenced but not in the tile:
new-diego (referenced by template "auctioneer" in "diego_brain" job)
new-diego (referenced by template "some-job" in "diego_brain" job)
capi-release (referenced by template "some-other-job" in "diego_brain" job)
capi-release (referenced by template "some-cc-job" in "cloud_controller" job)
`))
	})
})
