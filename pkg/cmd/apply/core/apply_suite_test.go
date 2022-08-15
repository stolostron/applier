// Copyright Red Hat

package core

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	genericclioptionsapplier "github.com/stolostron/applier/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/stolostron/applier/pkg/helpers"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

const (
	// The directory for test environment assets
	TestEnvDir string = ".testenv"
	// the test environment kubeconfig file
	TestEnvKubeconfigFile string = TestEnvDir + "/testenv.kubeconfig"
)

var testEnv *envtest.Environment
var restConfig *rest.Config
var kubeClient kubernetes.Interface
var apiExtensionsClient apiextensionsclient.Interface
var dynamicClient dynamic.Interface
var GvrSCR schema.GroupVersionResource = schema.GroupVersionResource{Group: "example.com", Version: "v1", Resource: "samplecustomresources"}
var root *cobra.Command
var applierFlags *genericclioptionsapplier.ApplierFlags
var streams genericclioptions.IOStreams

func TestApply(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TemplateFunction Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")

	Expect(os.MkdirAll(TestEnvDir, 0700)).To(BeNil())
	// start a kube-apiserver
	testEnv = &envtest.Environment{
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "test", "unit", "resources", "scenario", "config", "crd", "crd.yaml"),
		},
	}

	if testEnv.UseExistingCluster == nil {
		boolFalse := false
		testEnv.UseExistingCluster = &boolFalse
	}
	var err error
	var hubKubeconfig *rest.Config
	if *testEnv.UseExistingCluster {
		_, hubKubeconfig, err = PersistAndGetRestConfig(*testEnv.UseExistingCluster)
		Expect(err).ToNot(HaveOccurred())
		testEnv.Config = hubKubeconfig
	} else {
		Expect(os.Setenv("KUBECONFIG", "")).To(BeNil())
	}

	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	// Save the testenv kubeconfig
	if !*testEnv.UseExistingCluster {
		_, _, err = PersistAndGetRestConfig(*testEnv.UseExistingCluster)
		Expect(err).ToNot(HaveOccurred())
	}

	kubeClient, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	apiExtensionsClient, err = apiextensionsclient.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	dynamicClient, err = dynamic.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	restConfig = cfg
	PersistAndGetRestConfig(false)

	root, applierFlags, streams = helpers.NewRootCmd()

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func PersistAndGetRestConfig(useExistingCluster bool) (string, *rest.Config, error) {
	var err error
	buf := new(strings.Builder)
	if useExistingCluster {
		cmd := exec.Command("kubectl", "config", "view", "--raw")
		cmd.Stdout = buf
		cmd.Stderr = buf
		err = cmd.Run()
	} else {
		adminInfo := envtest.User{Name: "admin", Groups: []string{"system:masters"}}
		authenticatedUser, err := testEnv.AddUser(adminInfo, testEnv.Config)
		Expect(err).To(BeNil())
		kubectl, err := authenticatedUser.Kubectl()
		Expect(err).To(BeNil())
		var out io.Reader
		out, _, err = kubectl.Run("config", "view", "--raw")
		Expect(err).To(BeNil())
		_, err = io.Copy(buf, out)
		Expect(err).To(BeNil())
	}
	if err != nil {
		return "", nil, err
	}
	if err := ioutil.WriteFile(TestEnvKubeconfigFile, []byte(buf.String()), 0600); err != nil {
		return "", nil, err
	}

	hubKubconfigData, err := ioutil.ReadFile(TestEnvKubeconfigFile)
	if err != nil {
		return "", nil, err
	}
	hubKubeconfig, err := clientcmd.RESTConfigFromKubeConfig(hubKubconfigData)
	if err != nil {
		return "", nil, err
	}
	return buf.String(), hubKubeconfig, err
}

var _ = Describe("apply resources files", func() {
	It("Create resources", func() {
		cmd := NewCmd(applierFlags, streams)
		root.AddCommand(cmd)
		root.SetArgs([]string{
			"core-resources",
			"--path", "../../../../test/unit/resources/scenario/multicontent/clusterrole.yaml",
			"--path", "../../../../test/unit/resources/scenario/multicontent/clusterrolebinding.yaml",
			"--path", "../../../../test/unit/resources/scenario/multicontent/file1.yaml",
			"--values", "../../../../test/unit/resources/scenario/values.yaml",
			"--kubeconfig", TestEnvKubeconfigFile,
		})
		err := cmd.Execute()
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "my-ns", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("my-ns").Get(context.TODO(), "my-sa", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})
