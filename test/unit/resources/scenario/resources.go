// Copyright Red Hat
package scenario

import (
	"embed"

	"github.com/stolostron/applier/pkg/asset"
)

//go:embed musttemplateasset ownerref multicontent render/results
var files embed.FS

func GetScenarioResourcesReader() *asset.ScenarioResourcesReader {
	return asset.NewScenarioResourcesReader(&files)
}
