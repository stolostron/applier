[comment]: # ( Copyright Red Hat )

# IMPORTANT

This is the next generation of the applier which took a totaly different approach to create and update resources on kubebernetes and it is not compatible with the v1.0.1 version. Now the applier relies on the openshift/libragy.go to apply the rendered files on kubeberentes.
You can fork the V1.0.1 latest version if you want to continue to improve it or switch to the V1.1.0 version.
# Applier

The applier applies templated resources on kubebernetes. It can be use as a CLI or as a Go package in your code allowing you to apply embeded templates to your clusters.
## Introduction to template

The template supports the [text/template](https://golang.org/pkg/text/template/) framework and so you can use statements defined in that framework.
As the [Mastermind/sprig](https://github.com/Masterminds/sprig) is also loaded, you can use any functions defined by that framework.
By enriching the [templatefunction.go](pkg/templateprocessor/templatefunction.go), you can also develop your own functions. Check for example the function `toYaml` in the [templatefunction.go](pkg/templateprocessor/templatefunction.go).
Available functions (available if `WithTemplateFuncMap(applier.FuncMap())` is called while building the applier):
- `toYaml` which marshal a Go object to yaml.
- `encodeBase64` which base64 encode a string, but `b64enc` from sprig can be used.
- `include` which include a template.

A Header file can be specified containing go/text `block` or `define` and if so will be included at the beginning of each template. This allows you to add extra business logic in your templates.
## Template examples:

- `applier render --values examples/values.yaml --paths examples/simple`
- `applier render --values examples/values.yaml --paths examples/header/templates --header examples/header/header.txt`


## Packages
### Methods

The package provides methods to apply or render resources. These functions must be called on a [Applier](pkg/apply/apply.go#L133). An Applier can be build using the function NewApplierBuilder as follow:

```Go
applier := applierBuilder.
	WithClient(kubeClient, apiExtensionsClient, dynamicClient).
	Build()
```

There is other WithXxxx functions you can call on the applierBuilder or on the applier itself such as `WithTemplateFuncMap`, `WithOwner`, `WithCache`, `WithContext`, `WithKindOrder`...

Once you have the applier you can call one of the following method.
- [Apply](pkg/apply/apply.go) which will call `ApplyDirectly`, `ApplyCustomResources` or `ApplyDeployments` depending on the kind of resources.
- [ApplyDirectly](pkg/apply/apply.go) which takes kubernetes core resources from a reader such as namespace, secret... and apply them with the provided values.
- [ApplyCustomResources](pkg/apply/apply.go) which takes custom resources from a reader and apply them with the provided values.
- [ApplyDeployments](pkg/apply/apply.go) which teakes kubernetes Deployments from a reader and apply them with the provided values.
- [MustTemplateResources](pkg/apply/apply.go) which takes resources from a reader and render it with the provided values.

### Readers

Three readers are available:
- `asset.NewMemFSReader()` which allows you to read files from given directories
- `asset.NewDirectoriesReader()` which allows you to read files from given directories
- `asset.GetScenarioResourcesReader()` which allows you to read resources from your project. In order to use it you have to add such code in the directory you want to read:
```Go
// Copyright Red Hat
package scenario

import (
	"embed"

	"github.com/stolostron/applier/pkg/asset"
)

//go:embed musttemplateasset ownerref
var files embed.FS

func GetScenarioResourcesReader() *asset.ScenarioResourcesReader {
	return asset.NewScenarioResourcesReader(&files)
}
```
and then call the GetScenarioResourcesReader() to get the reader.

### Examples:

Check the [command line apply code](pkg/cmd/apply/common/exec.go) to apply a list of files and this to just render [command line render code](/Users/dvernier/acm/applier/pkg/cmd/render/exec.go).

## command-line

A command-line is available to apply or render yaml files of a given directory. 
To generate the command line you can clone this project and then run either: 
- `make install` to install from your local environment
- `make oc-plugin` to install as a `oc` plugin
- `make kubectl-plugin` to install as a `kubectl` plugin

or you can run

```
kubectl krew install applier
```
To install the krew plugin follow: [krew quickstart](https://krew.sigs.k8s.io/docs/user-guide/quickstart/)

To get the usage, run:
```
[oc|kubectl] applier -h 
```

## apply command 

The apply command can be use as is or with one of these 3 subcommands `core-reources`, `custom-resources` or `deployments`. Using it directly as `applier apply [options]` allows you to have a mix of core, custom and deployment resources in the `--path` option. The applier will sort the resource depending on their kind before applying them.

By default, the option `--sort-on-kind` is set to true and so the files will be sorted based on the kind. For example, namespace will be placed before serviceaccount.

For example you can run:

```bash
applier apply core-resources --path ./examples/simple --values ./examples/values.yaml
```
or
```
cat ./examples/values.yaml | applier apply core-resources --path ./examples/simple
```

The generated yaml file can be shown with option `--output-file`.
Dry-run can be enabled with the option `--dry-run`.
The combination of `--dry-run` and `--output-file /dev/stdout` (as the bellow `render` command) with ` | kubectl apply -f  -` allows to apply apply any kind of resources and not only `core`, `custom` and `deployments` as the resources template in that case will be only rendered.

For example:
```
applier apply core-resources --path ./examples/simple --values ./examples/values.yaml --dry-run --output-file /dev/stdout
```
will not apply the resources and will display the following:

```
# Copyright Red Hat

apiVersion: v1
kind: Namespace
metadata:
  name: "my-ns"
---
# Copyright Red Hat

apiVersion: v1
kind: ServiceAccount
metadata:
  name: "my-sa"
  namespace: "my-ns"
secrets:
- name: mysecret
---
```
## render command

The `render` command is similar than using the `apply` command with the options `--dry-run` and `--output-file /dev/stdout`

```
applier render --path ./examples/simple --values ./examples/values.yaml
```
or
```
cat ./example/values.yaml | apply render --path ./examples/simple
```

The result will be:

```
# Copyright Red Hat

apiVersion: v1
kind: Namespace
metadata:
  name: "my-ns"
---
# Copyright Red Hat

apiVersion: v1
kind: ServiceAccount
metadata:
  name: "my-sa"
  namespace: "my-ns"
secrets:
- name: mysecret
---
```

The `render` subcommand can be use in conjunction with `| kubectl apply -f -` to apply the generated yaml file.



