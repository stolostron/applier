// Copyright Contributors to the Open Cluster Management project

// +build functional

package functional_test

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Values map[string]interface{}

var (
	clientHub        client.Client
	clientHubDynamic dynamic.Interface
	clientAPIExt     clientset.Interface

	gvrSample schema.GroupVersionResource
)

func init() {
	klog.SetOutput(GinkgoWriter)
	klog.InitFlags(nil)
}

var _ = BeforeSuite(func() {
	By("Setup Hub client")
	gvrSample = schema.GroupVersionResource{Group: "functional-test.open-cluster-management.io", Version: "v1", Resource: "samples"}

	var kubeconfig string

	kubeconfig = os.Getenv("KUBECONFIG")
	var apiconfig *api.Config
	if kubeconfig == "" {
		if usr, err := user.Current(); err == nil {
			kubeconfig = filepath.Join(usr.HomeDir, ".kube", "config")
		}
	}

	klog.Infof("Kubeconfig=%s", kubeconfig)
	apiconfig, err := clientcmd.LoadFromFile(kubeconfig)
	Expect(err).To(BeNil())
	config, err := clientcmd.NewDefaultClientConfig(
		*apiconfig,
		&clientcmd.ConfigOverrides{}).ClientConfig()
	Expect(err).To(BeNil())
	clientHub, err = client.New(config, client.Options{})
	Expect(err).To(BeNil())
	clientHubDynamic, err = dynamic.NewForConfig(config)
	Expect(err).To(BeNil())
	clientAPIExt, err = clientset.NewForConfig(config)
	Expect(err).To(BeNil())
})

func TestFunctional(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional Suite")
}
