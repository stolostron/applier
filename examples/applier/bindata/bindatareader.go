// Copyright Contributors to the Open Cluster Management project

package bindata

import (
	"github.com/ghodss/yaml"
)

type Bindata struct{}

func (*Bindata) Asset(name string) ([]byte, error) {
	return Asset(name)
}

func (*Bindata) AssetNames() ([]string, error) {
	return AssetNames(), nil
}

func (*Bindata) ToJSON(b []byte) ([]byte, error) {
	return yaml.YAMLToJSON(b)
}

func NewBindataReader() *Bindata {
	return &Bindata{}
}
