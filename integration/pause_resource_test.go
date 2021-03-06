package integration_test

import (
	"fmt"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Fly CLI", func() {
	Describe("Pause Resource", func() {
		var (
			flyCmd     *exec.Cmd
			reqsBefore int
		)

		Context("when the resource flag is provided", func() {
			pipelineName := "pipeline"
			resourceName := "resource-name-potato"
			fullResourceName := fmt.Sprintf("%s/%s", pipelineName, resourceName)

			BeforeEach(func() {
				flyCmd = exec.Command(flyPath, "-t", targetName, "pause-resource", "-r", fullResourceName)
				reqsBefore = len(atcServer.ReceivedRequests())
			})

			Context("when a resource is paused using the API", func() {
				BeforeEach(func() {
					apiPath := fmt.Sprintf("/api/v1/teams/main/pipelines/%s/resources/%s/pause", pipelineName, resourceName)
					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("PUT", apiPath),
							ghttp.RespondWith(http.StatusOK, nil),
						),
					)
				})

				It("successfully pauses the resource", func() {
					Expect(func() {
						sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())

						<-sess.Exited
						Expect(sess.ExitCode()).To(Equal(0))
						Eventually(sess).Should(gbytes.Say(fmt.Sprintf("paused '%s'\n", resourceName)))
					}).To(Change(func() int {
						return len(atcServer.ReceivedRequests())
					}).By(2))
				})
			})

			Context("when a resource is paused using the API", func() {
				BeforeEach(func() {
					apiPath := fmt.Sprintf("/api/v1/teams/main/pipelines/%s/resources/%s/pause", pipelineName, resourceName)
					atcServer.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("PUT", apiPath),
							ghttp.RespondWith(http.StatusInternalServerError, nil),
						),
					)
				})

				It("exits 1 and outputs an error", func() {
					Expect(func() {
						sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())

						Eventually(sess.Err).Should(gbytes.Say(`error`))

						<-sess.Exited
						Expect(sess.ExitCode()).To(Equal(1))
					}).To(Change(func() int {
						return len(atcServer.ReceivedRequests())
					}).By(2))
				})
			})
		})

		Context("when the resource flag is not provided", func() {
			BeforeEach(func() {
				flyCmd = exec.Command(flyPath, "-t", targetName, "pause-resource")
			})

			It("exits 1 and outputs an error", func() {
				sess, err := gexec.Start(flyCmd, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(sess.Err).Should(gbytes.Say(`error`))

				<-sess.Exited
				Expect(sess.ExitCode()).To(Equal(1))
			})
		})
	})
})
