package main

import (
	"encoding/json"
	"fmt"
	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v2"
	_ "gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

type MultiVariableString []string

func (m *MultiVariableString) Set(str string, opt getopt.Option) error {
	*m = append(*m, str)
	_ = opt
	return nil
}

func (m *MultiVariableString) String() string {
	return strings.Join(m.Array(), ", ")
}

func (m *MultiVariableString) Array() []string {
	return *m
}

func (m *MultiVariableString) ParseMultipleOptions() []string {
	return m.Array()
}

type TwoVariableString []string

func (t *TwoVariableString) Set(str string, opt getopt.Option) error {
	if len(*t) < 2 {
		*t = append(*t, str)
	}
	_ = opt
	return nil
}

func (t *TwoVariableString) String() string {
	return strings.Join(t.Array(), ", ")
}

func (t *TwoVariableString) Array() []string {
	return *t
}

func (t *TwoVariableString) ParseMultipleOptions() []string {
	return t.Array()
}

type sqlcAware interface {
	getQueriesPath() string
	getOutputPath() string
}

type SqlcYamlConf struct {
	SQL []struct {
		Queries string `yaml:"queries"`
		Gen     struct {
			Go struct {
				Out string `yaml:"out"`
			} `yaml:"go"`
		} `yaml:"gen"`
	} `yaml:"sql"`
}

func (syc SqlcYamlConf) getQueriesPath() string {
	return syc.SQL[0].Queries
}

func (syc SqlcYamlConf) getOutputPath() string {
	return syc.SQL[0].Gen.Go.Out
}

type SqlcJsonConf struct {
	SQL []struct {
		Queries string `yaml:"queries"`
		Gen     struct {
			Go struct {
				Out string `yaml:"out"`
			} `yaml:"go"`
		} `yaml:"gen"`
	} `json:"sql"`
}

func (sjc SqlcJsonConf) getQueriesPath() string {
	return sjc.SQL[0].Queries
}

func (sjc SqlcJsonConf) getOutputPath() string {
	return sjc.SQL[0].Gen.Go.Out
}

func main() {
	var tableNames MultiVariableString
	var relation = false
	var oneTableName TwoVariableString
	var manyTableNames TwoVariableString
	var pivotTableName string
	var help = false
	getopt.FlagLong(&tableNames, "table-name", 't', "Use this option to provide a "+
		"single table name. This option can be used multiple times.")
	getopt.FlagLong(&help, "help", 'h', "Shows the help menu.")
	getopt.FlagLong(&relation, "relation", 'r', "Activates relationships mode. Relationship mode "+
		"will ignore all the --table-name (-t) options")
	getopt.FlagLong(&oneTableName, "one", 'o', "This option allows you to specify the name of "+
		"the \"one\" side table in a one-to-many or one-to-one relationship. Only works with relationships mode enabled.")
	getopt.FlagLong(&manyTableNames, "many", 'm', "This option allows you to specify the name of "+
		"the \"many\" side table in a one-to-many or many-to-many relationship. If you are dealing with a "+
		"many-to-many relationship you need to use this option exactly twice. Only works with relationships mode enabled.")
	getopt.FlagLong(&pivotTableName, "pivot-name", 'p', "Use this option to indicate the pivot table"+
		" name in your many-to-many relationship. Only works with relationships mode enabled.")

	getopt.Parse()

	if help {
		getopt.PrintUsage(os.Stdout)
		return
	}

	sqlcConf := readSqlcConfFile()
	modelPath := strings.TrimSuffix(sqlcConf.getOutputPath(), "/")
	modelPath = fmt.Sprintf("%s/models.go", modelPath)

	if !relation {
		tableNames.ParseMultipleOptions()

		for _, tableName := range tableNames {
			fmt.Printf("Table Name: %s\n", getSingularCamelCasedName(tableName))
			bc := BasicCrud{}

			bc.Construct(tableName, modelPath)
			bc.CreateQueriesFile(sqlcConf.getQueriesPath())
		}
	} else {
		manyTableNames.ParseMultipleOptions()
	}

}

func readSqlcConfFile() sqlcAware {
	fileName := ""
	fileType := ""
	if _, err := os.Stat("sqlc.yaml"); err == nil {
		fileType = "yaml"
		fileName = "sqlc.yaml"
	} else if _, err := os.Stat("sqlc.json"); err == nil {
		fileType = "json"
		fileName = "sqlc.json"
	}

	if fileName == "" {
		log.Panic("sqlc configuration file does not exists")
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Panicf("encountered error while reading %s. error: %v", fileName, err)
	}
	var conf sqlcAware
	if fileType == "yaml" {
		yamlConf := SqlcYamlConf{}
		err = yaml.Unmarshal(file, &yamlConf)
		if err != nil {
			log.Panicf("encountered error while parsing %s. error: %v", fileName, err)
		}
		conf = yamlConf
	} else if fileType == "json" {
		jsonConf := SqlcJsonConf{}
		err = json.Unmarshal(file, &jsonConf)
		if err != nil {
			log.Panicf("encountered error while parsing %s. error: %v", fileName, err)
		}
		conf = jsonConf
	}
	return conf
}
