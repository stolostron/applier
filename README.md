[comment]: # ( Copyright Contributors to the Open Cluster Management project )

# IMPORTANT

This is the v2 version of the applier which took a totaly different approach to create and update resources on kubebernetes and it is not compatible with the v1 version. Now the applier relies on the openshift/libragy.go to apply the rendered files on kubeberentes.
You can fork the V1 latest version if you want to continue to improve it or switch to the V2 version.
# Applier

The applier applies templated resources on kubebernetes. It can be use as a CLI or as a package in your code allowing you to apply embeded templates to your clusters.
## Introduction to template

The template supports the [text/template](https://golang.org/pkg/text/template/) framework and so you can use statements defined in that framework.
As the [Mastermind/sprig](https://github.com/Masterminds/sprig) is also loaded, you can use any functions defined by that framework.
By enriching the [templatefunction.go](pkg/templateprocessor/templatefunction.go), you can also develop your own functions. Check for example the function `toYaml` in the [templatefunction.go](pkg/templateprocessor/templatefunction.go).
Available functions (available if `WithTemplateFuncMap(applier.FuncMap())` is called while building the applier):
- `toYaml` which marshal a Go object to yaml.
- `encodeBase64` which base64 encode a string, but `b64enc` from sprig can be used.
- `include` which include a template.

A Header file can be specified and if so will be included at the beginning of each template
## Template examples:

- `applier render --values examples/values.yaml --paths examples/simple`
- `pplier render --values examples/values.yaml --paths examples/header/templates --header examples/header/header.txt`

## command-line

A command-line is available to apply yaml files in a given directory. 
To generate it run either: 
- `make build`, the `applier` executable will be in the `$GOPATH/bin` directory
- `make install` to install from your local environment
- `make oc-plugin` to install as a `oc` plugin
- `make kubectl-plugin` to install as a `kubectl` plugin

To get the usage, run:
```
[oc|kubectl] applier -h 
```

For example you can run:

```bash
applier apply core-resources --path ./examples/simple --values ./examples/values.yaml
```
or
```bash
applier render --path ./examples/simple --values ./examples/values.yaml | kubectl apply -f - 
```

## Packages
### Methods

The package provides functions to apply or render resources. These functions must be called on a [Applier](hpkg/apply/apply.go#L133). An Applier can be build using the function NewApplierBuilder as follow:

```Go
applier := applierBuilder.
	WithClient(kubeClient, apiExtensionsClient, dynamicClient).
	Build()
```

There is other Withxxx functions you can call on the applierBuilder such as `WithTemplateFuncMap`, `WithOwner`, `WithCache`, `WithContext`.

Once you have the applier you can call one of the following method.
- [ApplyDirectly](pkg/apply/apply.go#L102) which takes kubernetes core resources from a reader such as namespace, secret... and apply them with the provided values.
- [ApplyCustomResources](pkg/apply/apply.go#L133) which takes custom resources from a reader and apply them with the provided values.
- [ApplyDeployments](pkg/apply/apply.go#L49) which teakes kubernetes Deployments from a reader and apply them with the provided values.
- [MustTemplateResources](pkg/apply/apply.go#L249) which takes resources from a reader and render it with the provided values.

### Readers

Two readers are available:
- `asset.NewDirectoriesReader()` which allows you to read files from given directories
- `asset.GetScenarioResourcesReader()` which allows you to read resources from your project. In order to use it you have to add such code in the directory you want to read:
```Go
// Copyright Contributors to the Open Cluster Management project
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

### Example:

Check the [command line apply code](pkg/cmd/apply/common/exec.go) to apply a list of files and this to just render [command line render code](/Users/dvernier/acm/applier/pkg/cmd/render/exec.go).