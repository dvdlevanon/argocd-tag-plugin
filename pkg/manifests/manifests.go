package manifests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8yaml "k8s.io/apimachinery/pkg/util/yaml"
)

type Manifests []unstructured.Unstructured

const StdIn = "-"

func WriteManifests(writer io.Writer, manifests Manifests) error {
	for _, manifest := range manifests {
		if err := writeManifest(writer, manifest.Object); err != nil {
			return err
		}
	}

	return nil
}

func writeManifest(writer io.Writer, manifest map[string]interface{}) error {
	res, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("ToYAML: could not export into YAML: %s", err)
	}

	output := string(res)
	fmt.Fprintf(writer, "%s---\n", output)
	return nil
}

func ReadManifests(path string, reader io.Reader) (Manifests, error) {
	var manifests Manifests
	var err error

	if path == StdIn {
		manifests, err = readManifestData(reader)
		if err != nil {
			return manifests, err
		}
	} else {
		files, err := listFiles(path)
		if len(files) < 1 {
			return manifests, fmt.Errorf("no YAML or JSON files were found in %s", path)
		}
		if err != nil {
			return manifests, err
		}

		var errs []error
		manifests, errs = readFilesAsManifests(files)
		if len(errs) != 0 {
			errMessages := make([]string, len(errs))
			for idx, err := range errs {
				errMessages[idx] = err.Error()
			}
			return manifests, fmt.Errorf("could not read YAML/JSON files:\n%s", strings.Join(errMessages, "\n"))
		}
	}

	return manifests, nil
}

func listFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}

func readFilesAsManifests(paths []string) (result []unstructured.Unstructured, errs []error) {
	for _, path := range paths {
		rawdata, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		manifest, err := readManifestData(bytes.NewReader(rawdata))
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		result = append(result, manifest...)
	}

	return result, errs
}

func readManifestData(yamlData io.Reader) ([]unstructured.Unstructured, error) {
	decoder := k8yaml.NewYAMLOrJSONDecoder(yamlData, 1)

	var manifests []unstructured.Unstructured
	for {
		nxtManifest := unstructured.Unstructured{}
		err := decoder.Decode(&nxtManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(nxtManifest.Object) > 0 {
			manifests = append(manifests, nxtManifest)
		}
	}

	return manifests, nil
}
