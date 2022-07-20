// Copyright Contributors to the Open Cluster Management project

//go:build functional
// +build functional

package functional_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/stolostron/applier/pkg/templateprocessor"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stolostron/applier/pkg/applier"
	"gopkg.in/yaml.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Applier", func() {
	Context("Without Finalizer and no force", func() {
		It("Apply create/update resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample"), nil, clientHub, nil, nil, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.CreateOrUpdateInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrSample).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).Should(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).Should(BeNil())

			Consistently(func() error {
				has, _, err := hasCRDs(clientAPIExt, []string{"samples.functional-test.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).Should(BeNil())
		})

		It("Apply delete resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample"), nil, clientHub, nil, nil, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.DeleteInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrSample).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				has, _, err := hasCRDs(clientAPIExt, []string{"samples.functional-test.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).ShouldNot(BeNil())
		})
	})

	Context("With Finalizer and force", func() {
		It("Apply create/update resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample_with_finalizers"), nil, clientHub, nil, nil, nil)
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.CreateOrUpdateInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrSample).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).Should(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).Should(BeNil())

			Consistently(func() error {
				has, _, err := hasCRDs(clientAPIExt, []string{"samples.functional-test.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).Should(BeNil())
		})

		It("Apply delete resources", func() {
			applier, err := applier.NewApplier(templateprocessor.NewYamlFileReader("resources/sample_with_finalizers"),
				nil,
				clientHub,
				nil,
				nil,
				&applier.Options{
					ForceDelete: true,
				})
			Expect(err).Should(BeNil())

			b, err := ioutil.ReadFile(filepath.Clean("resources/sample/values.yaml"))
			Expect(err).Should(BeNil())

			valuesc := &Values{}
			err = yaml.Unmarshal(b, valuesc)
			Expect(err).Should(BeNil())

			values := Values{}
			values["Values"] = *valuesc

			Expect(applier.DeleteInPath("template", nil, true, values)).Should(BeNil())

			Consistently(func() error {
				ns := clientHubDynamic.Resource(gvrSample).Namespace("default")
				_, err = ns.Get(context.TODO(), "sample-name", metav1.GetOptions{})
				return err
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				secret := corev1.Secret{}
				return clientHub.Get(context.TODO(),
					types.NamespacedName{Name: "mysecret", Namespace: "default"},
					&secret)
			},
			).ShouldNot(BeNil())

			Consistently(func() error {
				has, _, err := hasCRDs(clientAPIExt, []string{"samples.functional-test.open-cluster-management.io"})
				if !has {
					return fmt.Errorf("CRD not not found")
				}
				return err
			},
			).ShouldNot(BeNil())
		})
	})

})

func hasCRDs(client clientset.Interface, expectedCRDs []string) (has bool, missingCRDs []string, err error) {
	missingCRDs = make([]string, 0)
	has = true
	clientAPIExtensionV1 := client.ApiextensionsV1()
	for _, crd := range expectedCRDs {
		klog.V(1).Infof("Check if %s exists", crd)
		_, errGet := clientAPIExtensionV1.CustomResourceDefinitions().Get(context.TODO(), crd, metav1.GetOptions{})
		if errGet != nil {
			if errors.IsNotFound(errGet) {
				missingCRDs = append(missingCRDs, crd)
				has = false
			} else {
				klog.V(1).Infof("Error while retrieving crd %s: %s", crd, errGet.Error())
				return false, missingCRDs, errGet
			}
		}
	}
	return has, missingCRDs, err
}
