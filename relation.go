package main

type cardinality int

const (
	ONE cardinality = iota
	MANY
)

type RelationshipCrud struct {
	firstTableName         string
	firstTableCardinality  cardinality
	secondTableName        string
	secondTableCardinality cardinality
	firstTableMethodName   string
	secondTableMethodName  string
	relationshipType       string
	pivotTableName         string
}

func (rc *RelationshipCrud) Construct(firstTableName, secondTableName string, c1, c2 cardinality) {
	rc.firstTableName = firstTableName
	rc.secondTableName = secondTableName
	rc.firstTableCardinality = c1
	rc.secondTableCardinality = c2

	if c1 == ONE && c2 == ONE {
		rc.relationshipType = "one-to-one"
	} else if c1 == MANY && c2 == ONE {
		rc.firstTableName, rc.secondTableName = rc.secondTableName, rc.firstTableName
		c1, c2 = c2, c1
		rc.relationshipType = "one-to-many"
	} else if c1 == ONE && c2 == MANY {
		rc.relationshipType = "one-to-many"
	} else if c1 == MANY && c2 == MANY {
		rc.relationshipType = "many-to-many"
	}

	rc.firstTableMethodName = getSingularCamelCasedName(firstTableName)
	rc.secondTableMethodName = getSingularCamelCasedName(secondTableName)
}

func (rc *RelationshipCrud) SetPivotTableName(pivotTableName string) {
	rc.pivotTableName = pivotTableName
}

//func (rc *RelationshipCrud) prepareQueries() ([]string, []string) {
//	var relationQueries []string
//	var reverseQueries []string
//
//	if rc.relationshipType == "one-to-one" {
//		idField := getSingular(rc.firstTableName) + "_id"
//		query := fmt.Sprintf("-- name: Get%sFor%s :one\n"+
//			"SELECT * FROM %s WHERE %s = $1;",
//			rc.secondTableMethodName, rc.firstTableMethodName, rc.secondTableName, idField)
//		reverseQueries = append(reverseQueries, query)
//	} else if rc.relationshipType == "one-to-many" {
//		idField := getSingular(rc.firstTableName) + "_id"
//		pluralMethodName := getPlural(rc.secondTableMethodName)
//		query := fmt.Sprintf("-- name: GetAll%sFor%s :many\n"+
//			"SELECT * FROM %s WHERE %s = $1;",
//			pluralMethodName, rc.firstTableMethodName, rc.secondTableName, idField)
//		reverseQueries = append(reverseQueries, query)
//		query = fmt.Sprintf("-- name: Get%sOf%s :many\n"+
//			"SELECT * FROM %s WHERE id = $1;",
//			rc.firstTableMethodName, rc.secondTableMethodName, rc.firstTableName)
//		relationQueries = append(relationQueries, query)
//	} else if rc.relationshipType == "many-to-many" {
//		//check for pivot table name
//		//get all from the first one
//		//get all from the second one
//	}
//}
