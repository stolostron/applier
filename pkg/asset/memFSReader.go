// Copyright Red Hat

package asset

type MemFS struct {
	data map[string][]byte
}

var _ ScenarioReader = &ScenarioResourcesReader{
	files: nil,
}

func NewMemFSReader() *MemFS {
	return &MemFS{
		data: make(map[string][]byte, 0),
	}
}

func (r *MemFS) AddAssetsFromScenarioReader(reader ScenarioReader) error {
	assets, err := reader.AssetNames(nil, nil)
	if err != nil {
		return err
	}
	for _, asset := range assets {
		b, err := reader.Asset(asset)
		if err != nil {
			return err
		}
		r.AddAsset(asset, b)
	}
	return nil
}

func (r *MemFS) AddAsset(fileName string, data []byte) {
	r.data[fileName] = data
}

func (r *MemFS) Asset(name string) ([]byte, error) {
	return r.data[name], nil
}

func (r *MemFS) AssetNames(prefixes, excluded []string) ([]string, error) {
	assetNames := make([]string, 0)

	for f := range r.data {
		if !isExcluded(f, prefixes, excluded) {
			assetNames = append(assetNames, f)
		}
	}
	return assetNames, nil
}
