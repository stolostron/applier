// Copyright Red Hat
package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/klog/v2"
)

func ToJSON(b []byte) ([]byte, error) {
	b, err := yaml.YAMLToJSON(b)
	if err != nil {
		klog.Errorf("err:%s\nyaml:\n%s", err, string(b))
		return nil, err
	}
	return b, nil
}

func ExtractAssets(r ScenarioReader, prefix, dir string, excluded []string) error {
	assetNames, err := r.AssetNames([]string{prefix}, excluded)
	if err != nil {
		return err
	}
	for _, assetName := range assetNames {
		relPath, err := filepath.Rel(prefix, assetName)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, relPath)

		if relPath == "." {
			path = filepath.Join(dir, filepath.Base(assetName))
		}
		err = os.MkdirAll(filepath.Dir(path), os.FileMode(0700))
		if err != nil {
			return err
		}
		data, err := r.Asset(assetName)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, os.FileMode(0600))
		if err != nil {
			return err
		}
	}
	return nil
}

func isExcluded(f string, files, excluded []string) bool {
	isExcluded := false
	for _, e := range excluded {
		if f == e {
			isExcluded = true
		}
	}
	if isExcluded {
		return true
	}
	isExcluded = true
	for _, d := range files {
		if strings.HasPrefix(f, d) {
			isExcluded = false
		}
	}
	return isExcluded
}
