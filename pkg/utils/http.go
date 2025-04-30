package utils

import (
	"fmt"
)

func BuildFileURL(relativePath string) string {
	return fmt.Sprintf("/images/%s", relativePath)
}
