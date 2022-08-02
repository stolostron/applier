// Copyright Contributors to the Open Cluster Management project
package asset

type ScenarioReader interface {
	//Retrieve an asset from the data source
	Asset(templatePath string) ([]byte, error)
	//List all available assets in the data source
	AssetNames(excluded []string) ([]string, error)
	ExtractAssets(prefix, dir string, excluded []string) error
	ToJSON(b []byte) ([]byte, error)
}
