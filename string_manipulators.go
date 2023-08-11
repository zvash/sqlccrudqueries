package main

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"log"
	"os"
	"regexp"
	"strings"
)

func getSingularCamelCasedName(tableName string) string {
	parts := strings.Split(tableName, "_")

	pluralizeClient := pluralize.NewClient()

	parts[len(parts)-1] = pluralizeClient.Singular(parts[len(parts)-1])
	tableName = strings.Join(parts, "_")
	return strcase.ToCamel(tableName)
}

func parseModel(modelName, modelFilePath string) []string {
	var res []string
	buf, err := os.ReadFile(modelFilePath)
	if err != nil {
		log.Panicf("Cannot open %s file. It's needed for creating the queries!", modelFilePath)
	}
	models := string(buf)
	pattern := regexp.MustCompile(fmt.Sprintf(`(?s)type %s struct \{\n(.*?)\n\}`, modelName))
	matched := pattern.FindStringSubmatch(models)
	if matched != nil && len(matched) > 1 {
		attributesPattern := regexp.MustCompile(`(?s)\t([A-Z][\w]*)`)
		for _, field := range attributesPattern.FindAllString(matched[1], -1) {
			res = append(res, strcase.ToSnake(strings.TrimSpace(field)))
		}
	}
	return res
}

func getPlural(word string) string {
	pluralizeClient := pluralize.NewClient()
	return pluralizeClient.Plural(word)
}

func getSingular(word string) string {
	pluralizeClient := pluralize.NewClient()
	return pluralizeClient.Singular(word)
}
