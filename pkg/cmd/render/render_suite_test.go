// Copyright Red Hat

package render

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"github.com/stolostron/applier/pkg/helpers"
	"github.com/stolostron/applier/test/unit/resources/scenario"
)

const (
	// The directory for test environment assets
	TestEnvDir string = ".testenv"
	// the test environment kubeconfig file
	TestEnvKubeconfigFile string = TestEnvDir + "/testenv.kubeconfig"
)

var applierFlags *genericclioptionsapplier.ApplierFlags
var streams genericclioptions.IOStreams
var tempFile *os.File

func TestApply(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TemplateFunction Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")
	var err error
	tempFile, err = ioutil.TempFile("", "render_test")
	Expect(err).ToNot(HaveOccurred())

	_, applierFlags, streams = helpers.NewRootCmd()

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")

})

var _ = Describe("render resources files", func() {
	It("Render resources", func() {
		cmd := NewCmd(applierFlags, streams)
		// root.AddCommand(cmd)
		cmd.SetArgs([]string{
			// "render",
			"--path", "../../../test/unit/resources/scenario/multicontent/clusterrole.yaml",
			"--path", "../../../test/unit/resources/scenario/multicontent/clusterrolebinding.yaml",
			"--path", "../../../test/unit/resources/scenario/multicontent/file1.yaml",
			"--path", "../../../test/unit/resources/scenario/multicontent/sample.yaml",
			"--values", "../../../test/unit/resources/scenario/values.yaml",
			"--output-file", tempFile.Name(),
		})
		err := cmd.Execute()
		Expect(err).To(BeNil())
		got, err := ioutil.ReadFile(tempFile.Name())
		Expect(err).To(BeNil())
		expect, err := scenario.GetScenarioResourcesReader().Asset("render/results/multicontent.yaml")
		Expect(err).To(BeNil())
		Expect(string(got)).To(Equal(string(expect)))
	})
})

var _ = Describe("render resources files non-sort", func() {
	It("Render resources non-sort", func() {
		cmd := NewCmd(applierFlags, streams)
		// root.AddCommand(cmd)
		cmd.SetArgs([]string{
			"render",
			"--path", "../../../test/unit/resources/scenario/multicontent/file1.yaml",
			"--values", "../../../test/unit/resources/scenario/values.yaml",
			"--sort-on-kind=false",
			"--output-file", tempFile.Name(),
		})
		err := cmd.Execute()
		Expect(err).To(BeNil())
		got, err := ioutil.ReadFile(tempFile.Name())
		Expect(err).To(BeNil())
		expect, err := scenario.GetScenarioResourcesReader().Asset("render/results/non-sort-multicontent.yaml")
		Expect(err).To(BeNil())
		Expect(string(got)).To(Equal(string(expect)))
	})
})
