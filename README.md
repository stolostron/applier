[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# Applier

The applier applies/creates/deletes templated resources on kubebernetes. A reader needs to be implemented to read the templated resources. Readers such as to read a directory structure or a string are provided and reader example for 'bindata' or 'go:embed' are provided.

## Introduction to template

The file [templateprocessor](pkg/templateprocessor) contains an number of methods allowing you to render template yaml files. 
The template supports the [text/template](https://golang.org/pkg/text/template/) framework and so you can use statements defined in that framework.
As the [Mastermind/sprig](https://github.com/Masterminds/sprig) is also loaded, you can use any functions defined by that framework.
By enriching the [templatefunction.go](pkg/templateprocessor/templatefunction.go), you can also develop your own functions. Check for example the function `toYaml` in the [templatefunction.go](pkg/templateprocessor/templatefunction.go).
Available functions:
- `toYaml` which marshal a Go object to yaml.
- `encodeBase64` which base64 encode a string, but `b64enc` from sprig can be used.
- `include` which include a template.
A `_helpers.tpl` file can also be added to define your own functions.
The resources are read by an Go object satisfying the [TemplateReader](pkg/templateprocessor/templateProcessor.go) reader.  
The reader is embedded in a applier.TemplateProcessor object
The resources are sorted in order to be applied in a kubernetes environment using a applier.Client

## Template examples:

- [example without _helpers.tpl](examples/applier/resources/yamlfilereader)
- [example with _helpers.tpl](test/functional/resources/sample)

## command-line

A command-line is available to apply yaml files in a given directory. 
To generate it run either: 
- `make build`, the `applier` executable will be in the `bin` directory
- `make install` to install from your local environment
- `make oc-plugin` to install as a `oc` plugin
- `make kubectl-plugin` to install as a `kubectl` plugin

To get the usage, run:
```
[oc|kubectl] applier -h 
```
The CLI accept values from pipe. These values are appened to the provided values.yaml. As the piped values are added at the end of the provided values.yaml, the piped values could override the values provided in values.yaml.

For example:
`echo "att1: val1" | applier -d <mydir>`

## Go Package
### Implementing a reader

A reader allows the applier to read the templates. Two readers are provided in the package:

- A Files reader: The files reader reads manage a single file template or all templates in a given directory. The applier has the capability to walk recursively in the directory if required. A new instance of the reader can be created with [NewYamlFileReader](pkg/templateprocessor/yamlfilereader.go)
- A String reader: The string reader reads templates from a string. Each template are separated by a delimiter. A new instance of the reader can be created with [NewYamlStringReader](pkg/templateprocessor/yamlstringreader.go)

A reader will read assets from a data source. You can find [testreader.go](pkg/templateprocessor/testreader.go) an example of a reader which reads the data from memory.

A `bindata` implementation can be found at [bindata](examples/applier/bindata/bindatareader.go)
A `embed` implementatation can be found at [go-embed](https://github.com/stolostron/cm-cli/blob/main/pkg/resources/resources.go)

### Methods

- In [templateprocessor](pkg/templateprocessor) there are methods which process the yaml templates, return them as a list of yamls or list of `unstructured.Unstructured`.
- There are also methods that sort these processed yaml templates depending of their `kind`. The order is defined in `kindOrder` variable which can be override.
Methods such as `CreateOrUpdateInPath` or `DeleteInPath` which `creates/update` or `delete` all resources definedd in a specific path. Other methods are available in the file [applier.go](pkg/applier/applier.go)

#### Example 1: Generate a templated yaml

```
	values := struct {
		ManagedClusterName      string
		ManagedClusterNamespace string
	}{
		ManagedClusterName:      saNsN.Name,
		ManagedClusterNamespace: saNsN.Namespace,
	}
	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		return nil, err
	}
	result, err := tp.TemplateAsset("hub/managedcluster/manifests/managedcluster-service-account.yaml", values)
	if err != nil {
		return nil, err
	}
```
The result contains a `[]byte` representing the templated yaml with the provided config.

#### Example 2: Generate a list of templated yaml

```
	values := struct {
		KlusterletNamespace   string
		BootstrapSecretName   string
		BootstrapSecretToken  string
		BootstrapSecretCaCert string
		ImagePullSecretName   string
		ImagePullSecretData   string
		ImagePullSecretType   corev1.SecretType
	}{
		KlusterletNamespace:   klusterletNamespace,
		BootstrapSecretName:   managedCluster.Name,
		BootstrapSecretToken:  base64.StdEncoding.EncodeToString(bootStrapSecret.Data["token"]),
		BootstrapSecretCaCert: base64.StdEncoding.EncodeToString(bootStrapSecret.Data["ca.crt"]),
		ImagePullSecretName:   imagePullSecret.Name,
		ImagePullSecretData:   base64.StdEncoding.EncodeToString(imagePullSecret.Data[".dockerconfigjson"]),
		ImagePullSecretType:   imagePullSecret.Type,
	}

	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		return nil, err
	}

	results, err := tp.TemplateResources([]string{
		"klusterlet/namespace.yaml",
		"klusterlet/image_pull_secret.yaml",
		"klusterlet/bootstrap_secret.yaml",
		"klusterlet/cluster_role.yaml",
		"klusterlet/cluster_role_binding.yaml",
		"klusterlet/service_account.yaml",
		"klusterlet/operator.yaml",
	}, values)

```
results contains a non-sorted `[][]bytes` each element is the templated yaml files using the provided values.

#### Example 3: Retreive a list of yaml files

```
	tp, err := NewTemplateProcessor(NewTestReader(assets), nil, nil)
	if err != nil {
		return nil, nil, err
	}
	crds, err = tp.Assets("klusterlet/crds", nil, true)
	if err != nil {
		return nil, nil, err
	}
```
The crds contains a `[][]byte` (non-sorted) of all yaml files found in `klusterlet/crds` directory and sub-directory using the provided config.

#### Example 4: Generate a sorted list of yaml files based using all templates in a given directory

```
	values := struct {
		KlusterletNamespace   string
		BootstrapSecretName   string
		BootstrapSecretToken  string
		BootstrapSecretCaCert string
		ImagePullSecretName   string
		ImagePullSecretData   string
		ImagePullSecretType   corev1.SecretType
	}{
		KlusterletNamespace:   klusterletNamespace,
		BootstrapSecretName:   managedCluster.Name,
		BootstrapSecretToken:  base64.StdEncoding.EncodeToString(bootStrapSecret.Data["token"]),
		BootstrapSecretCaCert: base64.StdEncoding.EncodeToString(bootStrapSecret.Data["ca.crt"]),
		ImagePullSecretName:   imagePullSecret.Name,
		ImagePullSecretData:   base64.StdEncoding.EncodeToString(imagePullSecret.Data[".dockerconfigjson"]),
		ImagePullSecretType:   imagePullSecret.Type,
	}

	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		return nil, nil, err
	}

	resutls, err := tp.TemplateResourcesInPathYaml(
		"klusterlet", nil, false, values)
	if err != nil {
		return nil, nil, err
	}
```
The results contains a `[][]byte`. The yaml files are sorted based on the Kind, Namespace and Name of the resource. All yaml files come from the `resources/klusterlet` (non-recursive) using the provided values.

#### Example 5: Create or update all resources defined in a directory

```
var merger bindata.Merger = func(current,
	new *unstructured.Unstructured,
) (
	future *unstructured.Unstructured,
	update bool,
) {
	if spec, ok := want.Object["spec"]; ok && 
	!reflect.DeepEqual(spec, current.Object["spec"]) {
		update = true
		current.Object["spec"] = spec
	}
	if rules, ok := want.Object["rules"]; ok && 
	!reflect.DeepEqual(rules, current.Object["rules"]) {
		update = true
		current.Object["rules"] = rules
	}
	if roleRef, ok := want.Object["roleRef"]; ok && 
	!reflect.DeepEqual(roleRef, current.Object["roleRef"]) {
		update = true
		current.Object["roleRef"] = roleRef
	}
	if subjects, ok := want.Object["subjects"]; ok && 
	!reflect.DeepEqual(subjects, current.Object["subjects"]) {
		update = true
		current.Object["subjects"] = subjects
	}
	return current, update
}
...
	values := struct {
		ManagedClusterName          string
		ManagedClusterNamespace     string
		BootstrapServiceAccountName string
	}{
		ManagedClusterName:          instance.Name,
		ManagedClusterNamespace:     instance.Name,
		BootstrapServiceAccountName: instance.Name + bootstrapServiceAccountNamePostfix,
	}

	tp, err := NewTemplateProcessor(NewTestReader(assets), nil)
	if err != nil {
		return nil, nil, err
	}

	a, err := applier.NewApplier(templateprocessor.NewTestReader(assets), nil, r.client, instance, r.scheme, merger)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = a.CreateOrUpdateInPath(
		"hub/managedcluster/manifests",
		[]string{"hub/managedcluster/manifests/managedcluster-service-account.yaml"},
		false,
		values,
	)

	if err != nil {
		return reconcile.Result{}, err
	}
```

This will create or update all resources located in the `hub/managedcluster/manifests` directory (non-recursive) except `hub/managedcluster/manifests/managedcluster-service-account.yaml`. The resources are sorted based on their Kind, Namespace and Name. A Merger function is passed as parameter to define if the update must occur or not and how to merge the current resource with the new resource.
