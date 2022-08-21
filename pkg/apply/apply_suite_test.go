// Copyright Red Hat

package apply

import (
	"context"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/stolostron/applier/test/unit/resources/scenario"
)

var testEnv *envtest.Environment
var restConfig *rest.Config
var kubeClient kubernetes.Interface
var apiExtensionsClient apiextensionsclient.Interface
var dynamicClient dynamic.Interface
var GvrSCR schema.GroupVersionResource = schema.GroupVersionResource{Group: "example.com", Version: "v1", Resource: "samplecustomresources"}

func TestApply(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TemplateFunction Suite")
}

var _ = BeforeSuite(func() {
	By("bootstrapping test environment")

	// start a kube-apiserver
	testEnv = &envtest.Environment{
		ErrorIfCRDPathMissing: true,
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "test", "unit", "resources", "scenario", "config", "crd", "crd.yaml"),
		},
	}

	cfg, err := testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	kubeClient, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	apiExtensionsClient, err = apiextensionsclient.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	dynamicClient, err = dynamic.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	restConfig = cfg
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("setOwnerRef", func() {
	It("Add OwnerRef to core item", func() {
		var nsOwner *corev1.Namespace
		By("Creating ns owner", func() {
			nsOwner = &corev1.Namespace{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Namespace",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-ns-owner-1",
				},
			}
			var err error
			nsOwner, err = kubeClient.CoreV1().Namespaces().Create(context.TODO(), nsOwner, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).WithOwner(nsOwner, false, false, scheme.Scheme).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyDirectly(reader, values, false, "", "ownerref/ns.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			ns, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), "my-ownerns", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(ns.GetOwnerReferences())).To(Equal(1))
			Expect(ns.OwnerReferences[0].APIVersion).To(Equal(nsOwner.APIVersion))
			Expect(ns.OwnerReferences[0].Kind).To(Equal(nsOwner.Kind))
			Expect(ns.OwnerReferences[0].Name).To(Equal(nsOwner.Name))
			Expect(ns.OwnerReferences[0].UID).To(Equal(nsOwner.UID))
			Expect(ns.OwnerReferences[0].Controller).To(BeNil())
			Expect(ns.OwnerReferences[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to CRD item", func() {
		var sampleOwner *unstructured.Unstructured
		By("Creating sample owner", func() {
			object := make(map[string]interface{})
			object["metadata"] = metav1.ObjectMeta{
				Name: "my-sampleowner-owner-1",
			}
			sampleOwner = &unstructured.Unstructured{
				Object: object,
			}
			sampleOwner.SetAPIVersion("example.com/v1")
			sampleOwner.SetKind("SampleCustomResource")
			var err error
			sampleOwner, err = dynamicClient.Resource(GvrSCR).Create(context.TODO(), sampleOwner, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).WithOwner(sampleOwner, false, false, scheme.Scheme).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-owner",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/sampleowner.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			sample, err := dynamicClient.Resource(GvrSCR).Get(context.TODO(), "my-owner", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(sample.GetOwnerReferences())).To(Equal(1))
			Expect(sample.GetOwnerReferences()[0].APIVersion).To(Equal(sampleOwner.GetAPIVersion()))
			Expect(sample.GetOwnerReferences()[0].Kind).To(Equal(sampleOwner.GetKind()))
			Expect(sample.GetOwnerReferences()[0].Name).To(Equal(sampleOwner.GetName()))
			Expect(sample.GetOwnerReferences()[0].UID).To(Equal(sampleOwner.GetUID()))
			Expect(sample.GetOwnerReferences()[0].Controller).To(BeNil())
			Expect(sample.GetOwnerReferences()[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to Deployment item", func() {
		var deployment *appsv1.Deployment
		By("Creating cluster owner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
			deployment, err = kubeClient.AppsV1().Deployments("my-ownerns").Get(context.TODO(), "my-deployment-owner", metav1.GetOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).WithOwner(deployment, false, false, scheme.Scheme).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			dep, err := kubeClient.AppsV1().Deployments("my-ownerns").Get(context.TODO(), "my-deployment", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(dep.GetOwnerReferences())).To(Equal(1))
			Expect(dep.OwnerReferences[0].APIVersion).To(Equal(deployment.APIVersion))
			Expect(dep.OwnerReferences[0].Kind).To(Equal(deployment.Kind))
			Expect(dep.OwnerReferences[0].Name).To(Equal(deployment.Name))
			Expect(dep.OwnerReferences[0].UID).To(Equal(deployment.UID))
			Expect(dep.OwnerReferences[0].Controller).To(BeNil())
			Expect(dep.OwnerReferences[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to Deployment item with controller and blockDeletion", func() {
		var deployment *appsv1.Deployment
		By("Creating cluster owner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner-controller",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
			deployment, err = kubeClient.AppsV1().Deployments("my-ownerns").Get(context.TODO(), "my-deployment-owner-controller", metav1.GetOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).WithOwner(deployment, true, true, scheme.Scheme).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-controller",
				Namespace: "my-ownerns",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			dep, err := kubeClient.AppsV1().Deployments("my-ownerns").Get(context.TODO(), "my-deployment-controller", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(dep.GetOwnerReferences())).To(Equal(1))
			Expect(dep.OwnerReferences[0].APIVersion).To(Equal(deployment.APIVersion))
			Expect(dep.OwnerReferences[0].Kind).To(Equal(deployment.Kind))
			Expect(dep.OwnerReferences[0].Name).To(Equal(deployment.Name))
			Expect(dep.OwnerReferences[0].UID).To(Equal(deployment.UID))
			Expect(*dep.OwnerReferences[0].Controller).To(BeTrue())
			Expect(*dep.OwnerReferences[0].BlockOwnerDeletion).To(BeTrue())
		})
	})
})

var _ = Describe("setOwnerRef applier level", func() {
	It("Add OwnerRef to core item", func() {
		var nsOwner *corev1.Namespace
		By("Creating ns owner", func() {
			nsOwner = &corev1.Namespace{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Namespace",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-ns-owner-o",
				},
			}
			var err error
			nsOwner, err = kubeClient.CoreV1().Namespaces().Create(context.TODO(), nsOwner, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).Build()
			applier = applier.WithOwner(nsOwner, false, false, scheme.Scheme)

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyDirectly(reader, values, false, "", "ownerref/ns.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			ns, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), "my-ownerns-o", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(ns.GetOwnerReferences())).To(Equal(1))
			Expect(ns.OwnerReferences[0].APIVersion).To(Equal(nsOwner.APIVersion))
			Expect(ns.OwnerReferences[0].Kind).To(Equal(nsOwner.Kind))
			Expect(ns.OwnerReferences[0].Name).To(Equal(nsOwner.Name))
			Expect(ns.OwnerReferences[0].UID).To(Equal(nsOwner.UID))
			Expect(ns.OwnerReferences[0].Controller).To(BeNil())
			Expect(ns.OwnerReferences[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to CRD item", func() {
		var sampleOwner *unstructured.Unstructured
		By("Creating sample owner", func() {
			object := make(map[string]interface{})
			object["metadata"] = metav1.ObjectMeta{
				Name: "my-sampleowner-owner-1-o",
			}
			sampleOwner = &unstructured.Unstructured{
				Object: object,
			}
			sampleOwner.SetAPIVersion("example.com/v1")
			sampleOwner.SetKind("SampleCustomResource")
			var err error
			sampleOwner, err = dynamicClient.Resource(GvrSCR).Create(context.TODO(), sampleOwner, metav1.CreateOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			applier = applier.WithOwner(sampleOwner, false, false, scheme.Scheme)

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-owner",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/sampleowner.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			sample, err := dynamicClient.Resource(GvrSCR).Get(context.TODO(), "my-owner", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(sample.GetOwnerReferences())).To(Equal(1))
			Expect(sample.GetOwnerReferences()[0].APIVersion).To(Equal(sampleOwner.GetAPIVersion()))
			Expect(sample.GetOwnerReferences()[0].Kind).To(Equal(sampleOwner.GetKind()))
			Expect(sample.GetOwnerReferences()[0].Name).To(Equal(sampleOwner.GetName()))
			Expect(sample.GetOwnerReferences()[0].UID).To(Equal(sampleOwner.GetUID()))
			Expect(sample.GetOwnerReferences()[0].Controller).To(BeNil())
			Expect(sample.GetOwnerReferences()[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to Deployment item", func() {
		var deployment *appsv1.Deployment
		By("Creating cluster owner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
			deployment, err = kubeClient.AppsV1().Deployments("my-ownerns-o").Get(context.TODO(), "my-deployment-owner", metav1.GetOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			applier = applier.WithOwner(deployment, false, false, scheme.Scheme)

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			dep, err := kubeClient.AppsV1().Deployments("my-ownerns-o").Get(context.TODO(), "my-deployment", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(dep.GetOwnerReferences())).To(Equal(1))
			Expect(dep.OwnerReferences[0].APIVersion).To(Equal(deployment.APIVersion))
			Expect(dep.OwnerReferences[0].Kind).To(Equal(deployment.Kind))
			Expect(dep.OwnerReferences[0].Name).To(Equal(deployment.Name))
			Expect(dep.OwnerReferences[0].UID).To(Equal(deployment.UID))
			Expect(dep.OwnerReferences[0].Controller).To(BeNil())
			Expect(dep.OwnerReferences[0].BlockOwnerDeletion).To(BeNil())
		})
	})
	It("Add OwnerRef to Deployment item with controller and blockDeletion", func() {
		var deployment *appsv1.Deployment
		By("Creating cluster owner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-owner-controller",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
			deployment, err = kubeClient.AppsV1().Deployments("my-ownerns-o").Get(context.TODO(), "my-deployment-owner-controller", metav1.GetOptions{})
			Expect(err).To(BeNil())
		})
		By("setReferenceOwner", func() {
			reader := scenario.GetScenarioResourcesReader()
			applierBuilder := NewApplierBuilder()
			applier := applierBuilder.
				WithClient(kubeClient, apiExtensionsClient, dynamicClient).
				WithTemplateFuncMap(FuncMap()).
				Build()

			applier = applier.WithOwner(deployment, true, true, scheme.Scheme)

			values := struct {
				Name      string
				Namespace string
			}{
				Name:      "my-deployment-controller",
				Namespace: "my-ownerns-o",
			}
			_, err := applier.ApplyCustomResources(reader, values, false, "", "ownerref/deployment.yaml")
			Expect(err).To(BeNil())
		})
		By("Checking Ownerref", func() {
			dep, err := kubeClient.AppsV1().Deployments("my-ownerns-o").Get(context.TODO(), "my-deployment-controller", metav1.GetOptions{})
			Expect(err).To(BeNil())
			Expect(len(dep.GetOwnerReferences())).To(Equal(1))
			Expect(dep.OwnerReferences[0].APIVersion).To(Equal(deployment.APIVersion))
			Expect(dep.OwnerReferences[0].Kind).To(Equal(deployment.Kind))
			Expect(dep.OwnerReferences[0].Name).To(Equal(deployment.Name))
			Expect(dep.OwnerReferences[0].UID).To(Equal(deployment.UID))
			Expect(*dep.OwnerReferences[0].Controller).To(BeTrue())
			Expect(*dep.OwnerReferences[0].BlockOwnerDeletion).To(BeTrue())
		})
	})
})

var _ = Describe("apply resources files", func() {
	It("Create resources", func() {
		reader := scenario.GetScenarioResourcesReader()
		applierBuilder := NewApplierBuilder()
		applier := applierBuilder.
			WithClient(kubeClient, apiExtensionsClient, dynamicClient).
			Build()
		files := []string{"multicontent/clusterrole.yaml",
			"multicontent/clusterrolebinding.yaml",
			"multicontent/file1.yaml",
			"multicontent/sample.yaml",
		}
		values := struct {
			Multicontent map[string]string
		}{
			Multicontent: map[string]string{
				"ServiceAccount": "compute-operator",
				"Namespace":      "compute-config",
			},
		}
		results, err := applier.Apply(reader, values, false, "", files...)
		Expect(err).To(BeNil())
		Expect(len(results)).To(Equal(5))
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "compute-config", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("compute-config").Get(context.TODO(), "compute-operator", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = dynamicClient.Resource(GvrSCR).Get(context.TODO(), "my-sample", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("apply resources directory", func() {
	It("Create resources", func() {
		reader := scenario.GetScenarioResourcesReader()
		applierBuilder := NewApplierBuilder()
		applier := applierBuilder.
			WithClient(kubeClient, apiExtensionsClient, dynamicClient).
			Build()
		files := []string{
			"multicontent",
			"ownerref",
		}
		values := struct {
			Name         string
			Namespace    string
			Multicontent map[string]string
		}{
			Name:      "my-owner",
			Namespace: "my-ownerns-dir",
			Multicontent: map[string]string{
				"ServiceAccount": "compute-operator",
				"Namespace":      "compute-config-dir",
			},
		}
		results, err := applier.Apply(reader, values, false, "", files...)
		Expect(err).To(BeNil())
		Expect(len(results)).To(Equal(8))
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "compute-config-dir", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("compute-config-dir").Get(context.TODO(), "compute-operator", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = dynamicClient.Resource(GvrSCR).Get(context.TODO(), "my-sample", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = dynamicClient.Resource(GvrSCR).Get(context.TODO(), "my-owner", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "my-ownerns-dir", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.AppsV1().Deployments("my-ownerns-dir").Get(context.TODO(), "my-owner", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("apply resources files and directory", func() {
	It("Create resources", func() {
		reader := scenario.GetScenarioResourcesReader()
		applierBuilder := NewApplierBuilder()
		applier := applierBuilder.
			WithClient(kubeClient, apiExtensionsClient, dynamicClient).
			Build()
		files := []string{
			"multicontent",
			"ownerref/ns.yaml",
		}
		values := struct {
			Name         string
			Namespace    string
			Multicontent map[string]string
		}{
			Name:      "my-owner",
			Namespace: "my-ownerns-mix",
			Multicontent: map[string]string{
				"ServiceAccount": "compute-operator",
				"Namespace":      "compute-config-mix",
			},
		}
		results, err := applier.Apply(reader, values, false, "", files...)
		Expect(err).To(BeNil())
		Expect(len(results)).To(Equal(6))
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "compute-config-mix", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("compute-config-mix").Get(context.TODO(), "compute-operator", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "my-ownerns-mix", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("applydirectly resources files", func() {
	It("Create resources", func() {
		reader := scenario.GetScenarioResourcesReader()
		applierBuilder := NewApplierBuilder()
		applier := applierBuilder.
			WithClient(kubeClient, apiExtensionsClient, dynamicClient).
			Build()
		files := []string{"multicontent/clusterrole.yaml",
			"multicontent/clusterrolebinding.yaml",
			"multicontent/file1.yaml",
		}
		values := struct {
			Multicontent map[string]string
		}{
			Multicontent: map[string]string{
				"ServiceAccount": "compute-operator",
				"Namespace":      "compute-config",
			},
		}
		results, err := applier.ApplyDirectly(reader, values, false, "", files...)
		Expect(err).To(BeNil())
		Expect(len(results)).To(Equal(4))
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "compute-config", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("compute-config").Get(context.TODO(), "compute-operator", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})

var _ = Describe("applydirectly resources files with restConfig", func() {
	It("Create resources", func() {
		reader := scenario.GetScenarioResourcesReader()
		applierBuilder := NewApplierBuilder()
		applier := applierBuilder.
			WithRestConfig(restConfig).
			Build()
		files := []string{"multicontent/clusterrole.yaml",
			"multicontent/clusterrolebinding.yaml",
			"multicontent/file1.yaml",
		}
		values := struct {
			Multicontent map[string]string
		}{
			Multicontent: map[string]string{
				"ServiceAccount": "compute-operator-rc",
				"Namespace":      "compute-config-rc",
			},
		}
		results, err := applier.ApplyDirectly(reader, values, false, "", files...)
		Expect(err).To(BeNil())
		Expect(len(results)).To(Equal(4))
		_, err = kubeClient.RbacV1().ClusterRoles().Get(context.TODO(), "cluster-role", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), "clusterrole-binding", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().Namespaces().Get(context.TODO(), "compute-config-rc", metav1.GetOptions{})
		Expect(err).To(BeNil())
		_, err = kubeClient.CoreV1().ServiceAccounts("compute-config-rc").Get(context.TODO(), "compute-operator-rc", metav1.GetOptions{})
		Expect(err).To(BeNil())
	})
})
