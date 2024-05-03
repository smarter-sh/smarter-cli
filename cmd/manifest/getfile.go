/*
Copyright © 2024 Lawrence McDaniel <lawrence@querium.com>
*/
package manifest

import (
	"fmt"
	"os"
)

func getYamlFileContents(kind string) (string, error) {
	filePath := fmt.Sprintf("./data/manifests/%s.yaml", kind)
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}
