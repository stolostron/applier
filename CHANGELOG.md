[comment]: # ( Copyright Red Hat )
# Release Content

## Additions
- Add `--sort-on-kind`option which sort the files based on their kind. For example a namespace resource will be created before a serviceaccount resource. 
- Add `--output-dir` option in the `applier render` allowing you to get the renderered files in a directory instead of a single output stream.
- Set `--sort-on-kind` default true
- By default, `ApplyDirectly` will sort the files. To Not sort it use `WithKindSort(apply.NoCreateUpdateKindSort)`
- The `--values` content can be passed using a pipe (ie: `cat values.yaml | applier render --path <file>` )
- Allowing resources files having multiple resource definition each being separated by `---`
- Refactor the reader interface, Only 2 methods must be implemented `Asset` and `AssetNames`
- Add `Apply` methods which runs `ApplyDirectly`, `ApplyCustomResources` and `ApplyDeployments`.
- Add command `applier apply ...` which call the `Apply` method.
- The `files` parameter in `Apply`, `ApplyDirectly`, `ApplyDeployments` and `ApplyCustomResources` can be an array of directory or file rather then only files.
- Add `--header` in the command-lines.
- Add `--excluded` to exclude files included by the `--path`

## Breaking changes

## Bug fixes

## Internal changes
