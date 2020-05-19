package exports

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/innosat-mats/rac-extract-payload/internal/common"
)

// GetJSONFilename replaces extension in template name with .json
func GetJSONFilename(templateName string) string {
	ext := filepath.Ext(templateName)
	return fmt.Sprintf(
		"%v.json",
		templateName[0:len(templateName)-len(ext)],
	)
}

// WriteJSON into target
func WriteJSON(target io.Writer, pkg *common.DataRecord, jsonFileName string) {
	err := json.NewEncoder(target).Encode(pkg)
	if err != nil {
		log.Printf("failed to encode json into %s", jsonFileName)
	}
}
