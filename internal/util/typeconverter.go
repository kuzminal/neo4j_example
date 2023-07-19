package util

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"neo4j/internal/model"
)

func ConvertToPerson(node neo4j.Node) (model.Person, error) {
	person := model.Person{}
	if node.Props == nil {
		return person, fmt.Errorf("could not convert record to Person")
	}
	person.Name = node.Props["name"].(string)
	if node.Props["age"] != nil {
		person.Age = node.Props["age"].(int64)
	}
	person.Id = node.GetElementId()
	return person, nil
}
