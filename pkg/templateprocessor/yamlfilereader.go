// Copyright Contributors to the Open Cluster Management project

package templateprocessor

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

//YamlFileReader defines a reader for yaml files
type YamlFileReader struct {
	path     string
	fileName string
}

var _ TemplateReader = &YamlFileReader{
	path:     "",
	fileName: "",
}

//Asset returns an asset
func (r *YamlFileReader) Asset(
	name string,
) ([]byte, error) {
	return ioutil.ReadFile(filepath.Clean(filepath.Join(r.path, name)))
}

//AssetNames returns the name of all assets
func (r *YamlFileReader) AssetNames() ([]string, error) {
	keys := make([]string, 0)
	var err error
	if r.fileName == "" {
		err = filepath.Walk(r.path, func(path string, info os.FileInfo, err error) error {
			if info != nil {
				if !info.IsDir() {
					newPath, err := filepath.Rel(r.path, path)
					if err != nil {
						return err
					}
					keys = append(keys, newPath)
				}
			}
			return nil
		})
	} else {
		helpersFile := filepath.Join(filepath.Base(r.path), "_helpers.tpl")
		if _, err := os.Stat(helpersFile); err == nil {
			keys = append(keys, "_helpers.tpl")
		}
		keys = append(keys, r.fileName)
	}
	return keys, err
}

//ToJSON converts to JSON
func (*YamlFileReader) ToJSON(
	b []byte,
) ([]byte, error) {
	b, err := yaml.YAMLToJSON(b)
	if err != nil {
		klog.Errorf("err:%s\nyaml:\n%s", err, string(b))
		return nil, err
	}
	return b, nil
}

//NewYamlFileReader constructs a new YamlFileReader
func NewYamlFileReader(
	path string,
) *YamlFileReader {
	reader := &YamlFileReader{
		path: path,
	}
	fi, err := os.Stat(path)
	if err != nil {
		klog.Fatal(err)
	}
	if !fi.Mode().IsDir() {
		reader.path = filepath.Dir(path)
		reader.fileName = filepath.Base(path)
	}
	return reader
}
