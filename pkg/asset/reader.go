// Copyright Red Hat
package asset

type ScenarioReader interface {
	//Retrieve an asset from the data source
	Asset(templatePath string) ([]byte, error)
	//List all available assets in the data source
	AssetNames(excluded []string) ([]string, error)
}
