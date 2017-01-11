package acceptance

import (
	"os/exec"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg-dep", func() {
	It("prints out a list of packages that have a dependency on the given dependency", func() {
		command := exec.Command(pathToMain, "--path", pathToTile, "pkg-dep", "-match", "golang")

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(string(session.Out.Contents())).To(ContainSubstring(`The following packages have a dependency that matches "golang":
Release: diego
  - acceptance-tests
  - auctioneer
`))
	})
})
