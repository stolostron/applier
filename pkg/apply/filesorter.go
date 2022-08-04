// Copyright Contributors to the Open Cluster Management project
package apply

import (
	"sort"

	"github.com/ghodss/yaml"
	"github.com/stolostron/applier/pkg/asset"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

//KindsOrder ...
type KindsOrder []string

//DefaultKindsOrder the default order
var DefaultCreateUpdateKindsOrder KindsOrder = []string{
	"Namespace",
	"NetworkPolicy",
	"ResourceQuota",
	"LimitRange",
	"PodSecurityPolicy",
	"PodDisruptionBudget",
	"ServiceAccount",
	"Secret",
	"SecretList",
	"ConfigMap",
	"StorageClass",
	"PersistentVolume",
	"PersistentVolumeClaim",
	"CustomResourceDefinition",
	"ClusterRole",
	"ClusterRoleList",
	"ClusterRoleBinding",
	"ClusterRoleBindingList",
	"Role",
	"RoleList",
	"RoleBinding",
	"RoleBindingList",
	"Service",
	"DaemonSet",
	"Pod",
	"ReplicationController",
	"ReplicaSet",
	"Deployment",
	"HorizontalPodAutoscaler",
	"StatefulSet",
	"Job",
	"CronJob",
	"Ingress",
	"APIService",
}

type fileInfo struct {
	fileName  string
	kind      string
	name      string
	namespace string
}

func (a *Applier) Sort(reader asset.ScenarioReader,
	values interface{},
	headerFile string,
	files ...string) ([]string, error) {
	filesInfo := make([]fileInfo, 0)
	for _, name := range files {
		b, err := a.MustTemplateAsset(reader, values, headerFile, name)
		if err != nil {
			return nil, err
		}
		unstructuredObj := &unstructured.Unstructured{}
		j, err := yaml.YAMLToJSON(b)
		if err != nil {
			return nil, err
		}

		err = unstructuredObj.UnmarshalJSON(j)
		if err != nil {
			return nil, err
		}
		filesInfo = append(filesInfo,
			fileInfo{
				fileName:  name,
				kind:      unstructuredObj.GetKind(),
				name:      unstructuredObj.GetName(),
				namespace: unstructuredObj.GetNamespace(),
			})
	}

	a.sortFiles(filesInfo)

	files = make([]string, len(filesInfo))
	for i, fileInfo := range filesInfo {
		files[i] = fileInfo.fileName
	}
	return files, nil
}

//sortUnstructuredForApply sorts a list on unstructured
func (a *Applier) sortFiles(filesInfo []fileInfo) {
	sort.Slice(filesInfo[:], func(i, j int) bool {
		return a.less(filesInfo[i], filesInfo[j])
	})
}

func (a *Applier) less(fileInfo1, fileInfo2 fileInfo) bool {
	if a.weight(fileInfo1) == a.weight(fileInfo2) {
		if fileInfo1.namespace == fileInfo2.namespace {
			return fileInfo1.name < fileInfo2.name
		}
		return fileInfo1.namespace < fileInfo2.namespace
	}
	return a.weight(fileInfo1) < a.weight(fileInfo2)
}

func (a *Applier) weight(fileInfo fileInfo) int {
	defaultWeight := len(a.kindOrder)
	for i, k := range a.kindOrder {
		if k == fileInfo.kind {
			return i
		}
	}
	return defaultWeight
}
