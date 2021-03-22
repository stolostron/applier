// Copyright Contributors to the Open Cluster Management project

package templateprocessor

import (
	"fmt"

	"github.com/ghodss/yaml"
)

//YamlFileReader defines a reader for map of string
type MapReader struct {
	assets map[string]string
}

var _ TemplateReader = &MapReader{assets: map[string]string{}}

//Asset returns an asset
func (r *MapReader) Asset(name string) ([]byte, error) {
	if s, ok := r.assets[name]; ok {
		return []byte(s), nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

//AssetNames returns the name of all assets
func (r *MapReader) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	for k := range r.assets {
		keys = append(keys, k)
	}
	return keys, nil
}

//ToJSON converts to JSON
func (r *MapReader) ToJSON(b []byte) ([]byte, error) {
	return yaml.YAMLToJSON(b)
}

//NewTestReader constructs a new YamlFileReader
func NewTestReader(assets map[string]string) *MapReader {
	return &MapReader{assets}
}
