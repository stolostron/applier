// Copyright Contributors to the Open Cluster Management project

package applier

import (
	"context"
	"reflect"
	"testing"

	"github.com/open-cluster-management/applier/pkg/templateprocessor"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestApplierClient_CreateOrUpdateInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewTestReader(assets)

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            values.BootstrapServiceAccountName,
			Namespace:       values.ManagedClusterNamespace,
			ResourceVersion: "0",
		},
	}
	saSecrets := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            values.BootstrapServiceAccountName,
			Namespace:       values.ManagedClusterNamespace,
			ResourceVersion: "0",
		},
		Secrets: []corev1.ObjectReference{
			{Name: "objectname"},
		},
	}
	clientUpdate := fake.NewFakeClient(sa)

	aUpdate, err := NewApplier(reader, nil, clientUpdate, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	// clientUpdateNoMerger := fake.NewFakeClient(sa)

	// aUpdateNoMerger, err := NewApplier(reader, nil, clientUpdateNoMerger, nil, nil, nil, nil)
	// if err != nil {
	// 	t.Errorf("Unable to create applier %s", err.Error())
	// }

	clientUpdateMerged := fake.NewFakeClient(saSecrets)

	aUpdateMerged, err := NewApplier(reader, nil, clientUpdateMerged, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
		{
			name:   "success update",
			fields: *aUpdate,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
		// {
		// 	name:   "success update no merger",
		// 	fields: *aUpdateNoMerger,
		// 	args: args{
		// 		path:      "test",
		// 		excluded:  nil,
		// 		recursive: false,
		// 		values:    values,
		// 	},
		// 	wantErr: true,
		// },
		{
			name:   "success update merged",
			fields: *aUpdateMerged,
			args: args{
				path: "test",
				excluded: []string{
					"test/clusterrolebinding",
					"test/clusterrole",
				},
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.CreateOrUpdateInPath(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplierClient.CreateOrUpdateInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				sa := &corev1.ServiceAccount{}
				err := client.Get(context.TODO(), types.NamespacedName{
					Name:      values.BootstrapServiceAccountName,
					Namespace: values.ManagedClusterNamespace,
				}, sa)
				if err != nil {
					t.Error(err)
				}
				r := &rbacv1.ClusterRole{}
				err = client.Get(context.TODO(), types.NamespacedName{
					Name: values.ManagedClusterName,
				}, r)
				if err != nil {
					t.Error(err)
				}
				rb := &rbacv1.ClusterRoleBinding{}
				err = client.Get(context.TODO(), types.NamespacedName{
					Name: values.ManagedClusterName,
				}, rb)
				if err != nil {
					t.Error(err)
				}
				if rb.RoleRef.Name != values.ManagedClusterName {
					t.Errorf("Expecting %s got %s", values.ManagedClusterName, rb.RoleRef.Name)
				}
				switch tt.name {
				case "success update":
					if len(sa.Secrets) == 0 {
						t.Error("Not merged as no secrets found")
					}
				case "success update merged":
					if sa.Secrets[0].Name != "mysecret" {
						t.Errorf("Not merged secrets=%#v", sa.Secrets[0])
					}
				}
			}
		})
	}
}

func TestNewApplier(t *testing.T) {
	client := fake.NewFakeClient([]runtime.Object{}...)
	owner := &corev1.Secret{}
	scheme := &runtime.Scheme{}
	type args struct {
		reader                   templateprocessor.TemplateReader
		templateProcessorOptions *templateprocessor.Options
		client                   crclient.Client
		owner                    metav1.Object
		scheme                   *runtime.Scheme
	}
	tests := []struct {
		name    string
		args    args
		want    *Applier
		wantErr bool
	}{
		{
			name: "Succeed",
			args: args{
				reader:                   templateprocessor.NewTestReader(assets),
				templateProcessorOptions: nil,
				client:                   client,
				owner:                    owner,
				scheme:                   scheme,
			},
			want: &Applier{
				templateProcessor: &templateprocessor.TemplateProcessor{},
				client:            client,
				owner:             owner,
				scheme:            scheme,
			},
			wantErr: false,
		},
		{
			name: "Failed no templateProcessor",
			args: args{
				reader:                   nil,
				templateProcessorOptions: nil,
				client:                   client,
				owner:                    owner,
				scheme:                   scheme,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failed no client",
			args: args{
				reader:                   nil,
				templateProcessorOptions: nil,
				client:                   nil,
				owner:                    owner,
				scheme:                   scheme,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewApplier(tt.args.reader, tt.args.templateProcessorOptions, tt.args.client, tt.args.owner, tt.args.scheme, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewApplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.templateProcessor, tt.want.templateProcessor) &&
					!reflect.DeepEqual(got.client, tt.want.client) &&
					!reflect.DeepEqual(got.owner, tt.want.owner) &&
					!reflect.DeepEqual(got.scheme, tt.want.scheme) {
					t.Errorf("NewApplier() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestApplier_setControllerReference(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "set-controller-reference",
			Namespace: values.ManagedClusterNamespace,
		},
	}

	type fields struct {
		templateProcessor *templateprocessor.TemplateProcessor
		client            crclient.Client
		owner             metav1.Object
		scheme            *runtime.Scheme
	}
	type args struct {
		u *unstructured.Unstructured
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: nil,
				owner:  sa,
				scheme: testscheme,
			},
			args: args{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": corev1.SchemeGroupVersion.String(),
						"kind":       "ServiceAccount",
						"metadata": map[string]interface{}{
							"name":      "set-controller-reference",
							"namespace": "myclusterns",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "success no owner",
			fields: fields{
				client: nil,
				owner:  nil,
				scheme: testscheme,
			},
			args: args{
				u: &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": corev1.SchemeGroupVersion.String(),
						"kind":       "ServiceAccount",
						"metadata": map[string]interface{}{
							"name":      "set-controller-reference",
							"namespace": "myclusterns",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
			}
			err := a.setControllerReference(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("Applier.setControllerReference() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				t.Log(tt.args.u.GetOwnerReferences())
				switch tt.name {
				case "success":
					if len(tt.args.u.GetOwnerReferences()) == 0 {
						t.Error("No ownerReference set")
					}
				case "success no owner":
					{
						if len(tt.args.u.GetOwnerReferences()) != 0 {
							t.Error("ownerReference found")
						}
					}
				}
			}
		})
	}
}

func TestApplier_CreateInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewTestReader(assets)

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      values.BootstrapServiceAccountName,
			Namespace: values.ManagedClusterNamespace,
		},
	}
	clientUpdate := fake.NewFakeClient(sa)

	aUpdate, err := NewApplier(reader, nil, clientUpdate, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
		{
			name:   "fail update",
			fields: *aUpdate,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
				applierOptions:    tt.fields.applierOptions,
			}
			if err := a.CreateInPath(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Applier.CreateInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplier_UpdateInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewTestReader(assets)

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	saSecrets := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            values.BootstrapServiceAccountName,
			Namespace:       values.ManagedClusterNamespace,
			ResourceVersion: "0",
		},
		Secrets: []corev1.ObjectReference{
			{Name: "objectname"},
		},
	}

	clientUpdateMerged := fake.NewFakeClient(saSecrets)

	aUpdateMerged, err := NewApplier(reader, nil, clientUpdateMerged, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}
	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "fail",
			fields: *a,
			args: args{
				path:      "test",
				excluded:  nil,
				recursive: false,
				values:    values,
			},
			wantErr: true,
		},
		{
			name:   "success update merged",
			fields: *aUpdateMerged,
			args: args{
				path: "test",
				excluded: []string{
					"test/clusterrolebinding",
					"test/clusterrole",
				},
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
				applierOptions:    tt.fields.applierOptions,
			}
			if err := a.UpdateInPath(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Applier.UpdateInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplier_DeleteInPath(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewTestReader(assets)

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      values.BootstrapServiceAccountName,
			Namespace: values.ManagedClusterNamespace,
		},
	}

	client := fake.NewFakeClient(sa)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	type args struct {
		path      string
		excluded  []string
		recursive bool
		values    interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success delete",
			fields: *a,
			args: args{
				path: "test",
				excluded: []string{
					"test/clusterrolebinding",
					"test/clusterrole",
				},
				recursive: false,
				values:    values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
				applierOptions:    tt.fields.applierOptions,
			}
			if err := a.DeleteInPath(tt.args.path, tt.args.excluded, tt.args.recursive, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Applier.DeleteInPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			saa := &corev1.ServiceAccount{}
			err := client.Get(context.TODO(),
				types.NamespacedName{
					Name:      sa.GetName(),
					Namespace: sa.GetNamespace()},
				saa)
			if err != nil && !errors.IsNotFound(err) {
				t.Error(err)
			}
		})
	}
}

func TestApplier_CreateOrUpdateResources(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRole{})
	testscheme.AddKnownTypes(rbacv1.SchemeGroupVersion, &rbacv1.ClusterRoleBinding{})
	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewYamlStringReader(assetsYaml, templateprocessor.KubernetesYamlsDelimiter)

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            values.BootstrapServiceAccountName,
			Namespace:       values.ManagedClusterNamespace,
			ResourceVersion: "0",
		},
	}
	clientUpdate := fake.NewFakeClient(sa)

	aUpdate, err := NewApplier(reader, nil, clientUpdate, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	type args struct {
		assets []string
		values interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				assets: []string{"0", "1", "2"},
				values: values,
			},
			wantErr: false,
		},
		{
			name:   "success update",
			fields: *aUpdate,
			args: args{
				assets: []string{"0", "1", "2"},
				values: values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
				applierOptions:    tt.fields.applierOptions,
			}
			if err := a.CreateOrUpdateResources(tt.args.assets, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Applier.CreateOrUpdateResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplier_CreateOrUpdateResource(t *testing.T) {
	testscheme := scheme.Scheme

	testscheme.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.ServiceAccount{})

	reader := templateprocessor.NewYamlStringReader(assetYaml, templateprocessor.KubernetesYamlsDelimiter)

	client := fake.NewFakeClient([]runtime.Object{}...)

	a, err := NewApplier(reader, nil, client, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	sa := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            values.BootstrapServiceAccountName,
			Namespace:       values.ManagedClusterNamespace,
			ResourceVersion: "0",
		},
	}
	clientUpdate := fake.NewFakeClient(sa)

	aUpdate, err := NewApplier(reader, nil, clientUpdate, nil, nil, nil)
	if err != nil {
		t.Errorf("Unable to create applier %s", err.Error())
	}

	type args struct {
		asset  string
		values interface{}
	}
	tests := []struct {
		name    string
		fields  Applier
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: *a,
			args: args{
				asset:  "0",
				values: values,
			},
			wantErr: false,
		},
		{
			name:   "success update",
			fields: *aUpdate,
			args: args{
				asset:  "0",
				values: values,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Applier{
				templateProcessor: tt.fields.templateProcessor,
				client:            tt.fields.client,
				owner:             tt.fields.owner,
				scheme:            tt.fields.scheme,
				applierOptions:    tt.fields.applierOptions,
			}
			if err := a.CreateOrUpdateResource(tt.args.asset, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("Applier.CreateOrUpdateResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
