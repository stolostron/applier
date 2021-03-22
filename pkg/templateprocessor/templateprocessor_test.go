// Copyright Contributors to the Open Cluster Management project

package templateprocessor

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var assetsB = []byte(`---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: system:test:{{ .ManagedClusterName }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:test:{{ .ManagedClusterName }}
subjects:
- kind: ServiceAccount
  name: {{ .BootstrapServiceAccountName }}
  namespace: {{ .ManagedClusterNamespace }}
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .BootstrapServiceAccountName }}
  namespace: {{ .ManagedClusterNamespace }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: system:test:{{ .ManagedClusterName }}
rules:
# Allow managed agent to rotate its certificate
- apiGroups: ['certificates.k8s.io']
  resources: ['certificatesigningrequests']
  verbs: ['create', 'get', 'list', 'watch']
# Allow managed agent to get
- apiGroups: ['cluster.open-cluster-management.io']
  resources: ['managedclusters']
  resourceNames: ['{{ .ManagedClusterName }}']
  verbs: ['get']
`)

func TestTemplateResourcesToMapOfUnstructuredWithStringReader(t *testing.T) {
	tp, err := NewTemplateProcessor(NewYamlStringReader(string(assetsB), KubernetesYamlsDelimiter), nil)
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	kindsOrder := []string{
		"ServiceAccount",
		"ClusterRole",
		"ClusterRoleBinding",
	}

	kindsNewOrder := []string{
		"ClusterRole",
		"ClusterRoleBinding",
		"ServiceAccount",
	}
	tpNewOrder, err := NewTemplateProcessor(NewYamlStringReader(string(assetsB), KubernetesYamlsDelimiter), &Options{CreateUpdateKindsOrder: kindsNewOrder})
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	type config struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}
	type args struct {
		path              string
		recursive         bool
		config            config
		templateProcessor *TemplateProcessor
		values            interface{}
	}
	type check struct {
		kinds []string
	}
	tests := []struct {
		name       string
		args       args
		check      check
		wantAssets map[string]*unstructured.Unstructured
		wantErr    bool
	}{
		{
			name: "Parse",
			args: args{
				path:      ".",
				recursive: false,
				config: config{
					ManagedClusterName:          "mymanagedcluster",
					ManagedClusterNamespace:     "mymanagedclusterNS",
					BootstrapServiceAccountName: "mymanagedcluster",
				},
				templateProcessor: tp,
				values:            values,
			},
			check: check{
				kinds: kindsOrder,
			},
			wantErr: false,
		},
		{
			name: "Parse new order",
			args: args{
				path:      ".",
				recursive: false,
				config: config{
					ManagedClusterName:          "mymanagedcluster",
					ManagedClusterNamespace:     "mymanagedclusterNS",
					BootstrapServiceAccountName: "mymanagedcluster",
				},
				templateProcessor: tpNewOrder,
				values:            values,
			},
			check: check{
				kinds: kindsNewOrder,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAssets, err := tt.args.templateProcessor.TemplateResourcesInPathUnstructured(tt.args.path, nil, tt.args.recursive, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetsUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotAssets) != 3 {
				t.Errorf("The number of unstructured asset must be 3 got: %d", len(gotAssets))
				return
			}
			for i := range gotAssets {
				if gotAssets[i].GetKind() != tt.check.kinds[i] {
					t.Errorf("Sort is not correct wanted %s and got: %s", tt.check.kinds[i], gotAssets[i].GetKind())
				}
			}
		})
	}
}

func TestTemplateResourcesToMapOfUnstructuredWithTestReady(t *testing.T) {
	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	kindsOrder := []string{
		"ServiceAccount",
		"ClusterRole",
		"ClusterRoleBinding",
	}

	kindsNewOrder := []string{
		"ClusterRole",
		"ClusterRoleBinding",
		"ServiceAccount",
	}
	tpNewOrder, err := NewTemplateProcessor(NewTestReader(assets), &Options{CreateUpdateKindsOrder: kindsNewOrder})
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	type config struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}
	type args struct {
		path              string
		recursive         bool
		config            config
		templateProcessor *TemplateProcessor
		values            interface{}
	}
	type check struct {
		kinds []string
	}
	tests := []struct {
		name       string
		args       args
		check      check
		wantAssets map[string]*unstructured.Unstructured
		wantErr    bool
	}{
		{
			name: "Parse",
			args: args{
				path:      "test",
				recursive: true,
				config: config{
					ManagedClusterName:          "mymanagedcluster",
					ManagedClusterNamespace:     "mymanagedclusterNS",
					BootstrapServiceAccountName: "mymanagedcluster",
				},
				templateProcessor: tp,
				values:            values,
			},
			check: check{
				kinds: kindsOrder,
			},
			wantErr: false,
		},
		{
			name: "Parse new order",
			args: args{
				path:      "test",
				recursive: true,
				config: config{
					ManagedClusterName:          "mymanagedcluster",
					ManagedClusterNamespace:     "mymanagedclusterNS",
					BootstrapServiceAccountName: "mymanagedcluster",
				},
				templateProcessor: tpNewOrder,
				values:            values,
			},
			check: check{
				kinds: kindsNewOrder,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAssets, err := tt.args.templateProcessor.TemplateResourcesInPathUnstructured(tt.args.path, nil, tt.args.recursive, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetsUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotAssets) != 3 {
				t.Errorf("The number of unstructured asset must be 3 got: %d", len(gotAssets))
				return
			}
			for i := range gotAssets {
				if gotAssets[i].GetKind() != tt.check.kinds[i] {
					t.Errorf("Sort is not correct wanted %s and got: %s", tt.check.kinds[i], gotAssets[i].GetKind())
				}
			}
		})
	}
}

func TestTemplateProcessor_TemplateAssetsInPathYaml(t *testing.T) {
	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	results := make([][]byte, 0)
	for _, y := range assets {
		results = append(results, []byte(y))
	}
	results[0] = append([]byte(copyright), results[0]...)
	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  TemplateProcessor
		args    args
		want    [][]byte
		wantErr bool
	}{
		{
			name:   "success",
			fields: *tp,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			want:    results,
			wantErr: false,
		},
		{
			name:   "failed missing values",
			fields: *tp,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    missingValues,
			},
			want:    results,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.TemplateResourcesInPathYaml(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("Applier.TemplateResourcesInPathYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if len(got) != len(tt.want) {
					t.Errorf("Applier.TemplateAssetsInPathYaml() returns %v yamls, want %v", len(got), len(tt.want))
				}
			}
		})
	}
}

func TestTemplateProcessor_Assets(t *testing.T) {
	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		t.Errorf("Unable to create templateProcessor %s", err.Error())
	}
	results := make([][]byte, 0)
	for _, y := range assets {
		results = append(results, []byte(y))
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
	}
	tests := []struct {
		name         string
		fields       TemplateProcessor
		args         args
		wantPayloads [][]byte
		wantErr      bool
	}{
		{
			name:   "success",
			fields: *tp,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
			},
			wantPayloads: results,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPayloads, err := tt.fields.Assets(tt.args.path, tt.args.excluded, tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("Applier.Assets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPayloads != nil {
				if len(gotPayloads) != len(tt.wantPayloads) {
					t.Errorf("Applier.TemplateAssetsInPathYaml() returns %v yamls, want %v", len(gotPayloads), len(tt.wantPayloads))
				}
			}
		})
	}
}

// func TestTemplateProcessor_TemplateBytesUnstructured(t *testing.T) {
// 	var values = struct {
// 		ManagedClusterName          string
// 		ManagedClusterNamespace     string
// 		BootstrapServiceAccountName string
// 	}{
// 		ManagedClusterName:          "mycluster",
// 		ManagedClusterNamespace:     "myclusterns",
// 		BootstrapServiceAccountName: "mysa",
// 	}

// 	type fields struct {
// 		reader  TemplateReader
// 		options *Options
// 	}
// 	type args struct {
// 		assets    []byte
// 		values    interface{}
// 		delimiter string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "succeed",
// 			fields: fields{
// 				reader:  NewTestReader(assets),
// 				options: &Options{},
// 			},
// 			args: args{
// 				assets:    assetsB,
// 				values:    values,
// 				delimiter: KubernetesYamlsDelimiter,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tp, _ := NewTemplateProcessor(
// 				tt.fields.reader,
// 				tt.fields.options)
// 			gotUs, err := tp.TemplateBytesUnstructured(tt.args.assets, tt.args.values, tt.args.delimiter)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TemplateProcessor.TemplateBytesUnstructured() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if len(gotUs) != 3 {
// 				t.Errorf("Got len %d, want 3", len(gotUs))
// 			}
// 		})
// 	}
// }

func TestTemplateProcessor_AssetNamesInPath(t *testing.T) {
	tpr := NewYamlStringReader(string(assetsB), KubernetesYamlsDelimiter)
	type fields struct {
		reader  TemplateReader
		options *Options
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "succeed",
			fields: fields{
				reader:  tpr,
				options: &Options{},
			},
			args: args{
				path:      ".",
				excluded:  nil,
				recursive: false,
			},
			want:    []string{"0", "1", "2"},
			wantErr: false,
		},
		{
			name: "succeed enpty",
			fields: fields{
				reader:  tpr,
				options: &Options{},
			},
			args: args{
				path:      "",
				excluded:  nil,
				recursive: true,
			},
			want:    []string{"0", "1", "2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp, _ := NewTemplateProcessor(
				tt.fields.reader,
				tt.fields.options)
			got, err := tp.AssetNamesInPath(tt.args.path, tt.args.excluded, tt.args.recursive)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateProcessor.AssetNamesInPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateProcessor.AssetNamesInPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_helpertpl(t *testing.T) {
	tpr := NewYamlFileReader("../../test/unit/resources/templates/withhelpers")
	tp, err := NewTemplateProcessor(tpr, &Options{})
	if err != nil {
		t.Error(err)
	}
	values := map[string]interface{}{
		"Values": map[string]string{
			"name": "TestTemplateProcessor_helpertpl",
		},
	}
	u, err := tp.TemplateResourcesInPathUnstructured(".", nil, false, values)
	if err != nil {
		t.Error(err)
	}
	if m, ok := u[0].Object["metadata"]; ok {
		metadata := m.(map[string]interface{})
		if metadata["name"].(string) != "Test" {
			t.Errorf("Expecting 'Test' got: %s", metadata["name"])
		}
	} else {
		t.Errorf("Malformed %v", u)
	}
}

func TestTemplateProcessor_BytesToUnstructured(t *testing.T) {
	var asset = []byte(``)
	tpr := NewYamlStringReader(string(asset), KubernetesYamlsDelimiter)
	type fields struct {
		reader  TemplateReader
		options *Options
	}
	type args struct {
		asset []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *unstructured.Unstructured
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				reader:  tpr,
				options: &Options{},
			},
			args: args{
				asset: []byte(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: myName
  namespace: myNamespace
---
`),
			},
			want: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ServiceAccount",
					"metadata": map[string]interface{}{
						"name":      "myName",
						"namespace": "myNamespace",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success empty",
			fields: fields{
				reader:  tpr,
				options: &Options{},
			},
			args: args{
				asset: []byte(`---
#apiVersion: v1
#kind: ServiceAccount
#metadata:
#  name: myName
#  namespace: myNamespace
---
`),
			},
			want:    &unstructured.Unstructured{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp, _ := NewTemplateProcessor(
				tt.fields.reader,
				tt.fields.options)
			got, err := tp.BytesToUnstructured(tt.args.asset)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateProcessor.BytesToUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateProcessor.BytesToUnstructured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateProcessor_BytesArrayToUnstructured(t *testing.T) {
	tpr := NewYamlStringReader(string(assetsB), KubernetesYamlsDelimiter)
	type fields struct {
		reader  TemplateReader
		options *Options
	}
	type args struct {
		assets [][]byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantUs  []*unstructured.Unstructured
		wantErr bool
	}{
		{
			name: "success, one is stripped because empty",
			fields: fields{
				reader:  tpr,
				options: &Options{},
			},
			args: args{
				assets: [][]byte{
					[]byte(`---
#apiVersion: v1
#kind: ServiceAccount
#metadata:
#  name: myName
#  namespace: myNamespace
---
`),
					[]byte(`---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: myName
  namespace: myNamespace
---
`)},
			},
			wantUs: []*unstructured.Unstructured{
				&unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ServiceAccount",
						"metadata": map[string]interface{}{
							"name":      "myName",
							"namespace": "myNamespace",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp, _ := NewTemplateProcessor(
				tt.fields.reader,
				tt.fields.options)
			gotUs, err := tp.BytesArrayToUnstructured(tt.args.assets)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateProcessor.BytesArrayToUnstructured() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUs, tt.wantUs) {
				t.Errorf("TemplateProcessor.BytesArrayToUnstructured() = %v, want %v", gotUs, tt.wantUs)
			}
		})
	}
}
