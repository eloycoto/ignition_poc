package ignition_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/eloycoto/ignition_poc/pkg/ignition"
)

var (
	sampleConfig = `
{
  "ignition": { "version": "3.0.0" },
  "storage": {
    "files": [{
      "path": "/opt/someconfig",
      "mode": 420,
      "contents": { "source": "data:,example%20file%0A" },
      "overwrite": true
    }]
  },
  "systemd": {
    "units": [{
      "name": "eloy.service",
      "enabled": true,
      "contents": "[Service]\nType=oneshot\nExecStart=/usr/bin/echo Hello World\n\n[Install]\nWantedBy=multi-user.target"
    }]
  }
}
`
)

var _ = Describe("Ignition", func() {

	It("Test", func() {
		Expect(true).To(BeTrue())

		cfg, err := ignition.ParseConfig(sampleConfig)
		Expect(cfg).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		// Checking that files are not in place

		Expect("/opt/someconfig").ShouldNot(BeAnExistingFile())

		command := exec.Command("systemctl", "is-active eloy.service")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		session.Wait()
		fmt.Println(session)
		Expect(session.ExitCode()).Should(Equal(1))

		command = exec.Command("systemctl", "is-enabled eloy.service")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		session.Wait()
		Expect(session.ExitCode()).Should(Equal(1))

		err = ignition.RunConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		Expect("/opt/someconfig").Should(BeAnExistingFile())

		command = exec.Command("systemctl", "is-active eloy.service")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		session.Wait()
		fmt.Println(session)
		Expect(session.ExitCode()).Should(Equal(0))

		command = exec.Command("systemctl", "is-enabled eloy.service")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		session.Wait()
		Expect(session.ExitCode()).Should(Equal(0))
	})

})
