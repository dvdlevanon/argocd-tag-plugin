package imagetag

import (
	"argocd-tag-plugin/pkg/ecr"
	"fmt"
	"strings"
)

type ImageTag struct {
	imageUrl       string
	tagConstraints string
	tagSelector    string
}

func NewImageTag(imageUrl string, tagConstraints string, tagSelector string) ImageTag {
	return ImageTag{
		imageUrl:       imageUrl,
		tagConstraints: tagConstraints,
		tagSelector:    tagSelector,
	}
}

func (t ImageTag) Process() (string, error) {
	if ecr.IsEcrUrl(t.imageUrl) {
		tags, err := ecr.FindAvailableTags(t.imageUrl, t.tagConstraints)
		if err != nil {
			return "", err
		}
		return t.findMatch(tags)
	}

	return "", fmt.Errorf("unsupported container registry %s", t.imageUrl)
}

func (t ImageTag) findMatch(tags []string) (string, error) {
	for _, tag := range tags {
		if t.isMatch(tag, t.tagSelector) {
			return tag, nil
		}
	}

	return "", fmt.Errorf("no matched tag found for %s (tags: %v) (selector: %s)", t.imageUrl, tags, t.tagSelector)
}

func (t ImageTag) isMatch(tag string, selector string) bool {
	if strings.HasPrefix(selector, "prefix(") && strings.HasSuffix(selector, ")") {
		prefix := selector[len("prefix(") : len(selector)-len(")")]
		return strings.HasPrefix(tag, prefix)
	}

	return true
}
