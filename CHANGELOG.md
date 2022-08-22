[comment]: # ( Copyright Red Hat )
# Release Content
- Add WithRestConfig(), it can be used instead of WithClient()
- Add WithOwner() at the applier level.
- Add WithContext at the applier level.
- Add WithKindOrder() at the applier level.
- Keep the original order when running with `--sort-on-kind=false`
- Now AssetNames with a nil prefix will skip the prefix test
## Additions

## Breaking changes

## Bug fixes

- Resources was not rendered when starting with "---"

## Internal changes
