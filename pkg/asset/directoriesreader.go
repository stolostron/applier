// Copyright Red Hat

package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//YamlFileReader defines a reader for yaml files
type YamlFileReader struct {
	header string
	paths  []string
	files  []string
}

var _ ScenarioReader = &YamlFileReader{
	header: "",
	paths:  []string{},
}

//NewDirectoriesReader constructs a new YamlFileReader
func NewDirectoriesReader(
	header string,
	paths []string,
) *YamlFileReader {
	reader := &YamlFileReader{
		header: header,
		paths:  paths,
	}
	return reader
}

//Asset returns an asset
func (r *YamlFileReader) Asset(
	name string,
) ([]byte, error) {
	files, err := r.AssetNames([]string{name}, nil)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 || len(files) > 1 {
		return nil, fmt.Errorf("file %s is not part of the assets", name)
	}
	return ioutil.ReadFile(filepath.Clean(name))
}

//AssetNames returns the name of all assets
func (r *YamlFileReader) AssetNames(prefixes, excluded []string) ([]string, error) {
	resultFiles := make([]string, 0)
	if len(r.header) != 0 {
		resultFiles = append(resultFiles, r.header)
	}
	visit := func(path string, fileInfo os.FileInfo, err error) error {
		if fileInfo == nil {
			return fmt.Errorf("paths %s doesn't exist", path)
		}
		if fileInfo.IsDir() {
			return nil
		}
		if isExcluded(path, prefixes, excluded) {
			return nil
		}
		resultFiles = append(resultFiles, path)
		return nil
	}

	for _, p := range r.paths {
		if err := filepath.Walk(p, visit); err != nil {
			return resultFiles, err
		}
	}
	return resultFiles, nil
}
