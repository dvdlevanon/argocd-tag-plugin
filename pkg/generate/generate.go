package generate

import (
	"argocd-tag-plugin/pkg/imagetag"
	"argocd-tag-plugin/pkg/manifests"
	"reflect"
	"regexp"
)

var extractRegex = regexp.MustCompile(`^(.*):<image-tag-plugin:(.*?)#(.*?)>.*$`)
var replaceRegex = regexp.MustCompile(`<image-tag-plugin:.*?>`)

func ProcessManifests(manifests manifests.Manifests) error {
	for _, manifest := range manifests {
		if err := process(&manifest.Object); err != nil {
			return err
		}
	}

	return nil
}

func process(node *map[string]interface{}) error {
	obj := *node
	for key, value := range obj {
		valueType := reflect.ValueOf(value).Kind()
		if valueType == reflect.Map {
			inner, ok := value.(map[string]interface{})
			if !ok {
				continue
			}
			if err := process(&inner); err != nil {
				return err
			}
		} else if valueType == reflect.Slice {
			for idx, elm := range value.([]interface{}) {
				switch elm := elm.(type) {
				case map[string]interface{}:
					{
						if err := process(&elm); err != nil {
							return err
						}
					}
				case string:
					{
						replacement, err := processString(key, elm)
						if err != nil {
							return err
						}
						value.([]interface{})[idx] = replacement
					}
				}
			}
		} else if valueType == reflect.String {
			replacement, err := processString(key, value.(string))
			if err != nil {
				return err
			}
			obj[key] = replacement
		}
	}

	return nil
}

func processString(key string, value string) (string, error) {
	matches := extractRegex.FindStringSubmatch(value)
	if len(matches) < 4 {
		return value, nil
	}

	dynamicTag := imagetag.NewImageTag(matches[1], matches[2], matches[3])
	tag, err := dynamicTag.Process()
	if err != nil {
		return value, err
	}

	return replaceRegex.ReplaceAllString(value, tag), nil
}
