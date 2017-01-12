package commands_test

import (
	"bytes"

	"github.com/ryanmoran/inspector/commands"
	"github.com/ryanmoran/inspector/commands/fakes"
	"github.com/ryanmoran/inspector/tiles"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PkgDep", func() {
	Describe("Execute", func() {
		It("finds packages that have a given dependency", func() {
			productParser := &fakes.ProductParser{}
			productParser.ParseCall.Returns.Product = tiles.Product{
				Releases: []tiles.Release{
					{
						Name: "some-release",
						CompiledPackages: []tiles.ReleasePackage{
							{
								Name: "some-package-1",
								Dependencies: []string{
									"otherdep-1",
									"somedep-1",
									"otherdep-2",
								},
							},
							{
								Name: "some-package-2",
								Dependencies: []string{
									"otherdep-2",
									"otherdep-1",
								},
							},
						},
						Packages: []tiles.ReleasePackage{
							{
								Name: "some-package-3",
								Dependencies: []string{
									"somedep-3",
									"otherdep-2",
									"another-somedep2",
								},
							},
						},
					},
				},
			}

			stdout := bytes.NewBuffer([]byte{})
			command := commands.NewPkgDep(productParser, stdout)

			err := command.Execute([]string{"-match", "somedep"})
			Expect(err).NotTo(HaveOccurred())

			Expect(stdout.String()).To(ContainSubstring(`Release: some-release
  - some-package-1 [somedep-1]
  - some-package-3 [somedep-3 another-somedep2]`))
		})

		Describe("error cases", func() {
			Context("when no dependency argument is specified", func() {
				It("returns an error", func() {

					command := commands.NewPkgDep(nil, nil)

					err := command.Execute([]string{})
					Expect(err).To(MatchError("-match is a required flag"))
				})
			})
		})
	})

	Describe("Usage", func() {
		It("returns a descriptive usage", func() {
			command := commands.NewPkgDep(nil, nil)
			Expect(command.Usage()).To(Equal(commands.Usage{
				Description:      "something something dep",
				ShortDescription: "something dep",
			}))
		})
	})
})
