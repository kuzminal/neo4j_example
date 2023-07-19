package graph

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	people = []map[string]any{
		{"name": "Alice", "age": 42, "friends": []string{"Anna"}},
		{"name": "Bob", "age": 19, "friends": []string{"Peter"}},
		{"name": "Peter", "age": 50, "friends": []string{"Anna"}},
		{"name": "Anna", "age": 30},
	}
)

func CreateUsers(ctx context.Context, driver neo4j.DriverWithContext, dbName string) {
	for _, person := range people {
		_, err := neo4j.ExecuteQuery(ctx, driver,
			"MERGE (p:Person {name: $person.name, age: $person.age})",
			map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase(dbName))
		if err != nil {
			panic(err)
		}
	}
}

func CreateRelationships(ctx context.Context, driver neo4j.DriverWithContext, dbName string) {
	for _, person := range people {
		if person["friends"] != "" {
			_, err := neo4j.ExecuteQuery(ctx, driver, `
                MATCH (p:Person {name: $person.name})
                UNWIND $person.friends AS friend_name
                MATCH (friend:Person {name: friend_name})
                MERGE (p)-[:KNOWS]->(friend)
                `, map[string]any{
				"person": person,
			}, neo4j.EagerResultTransformer,
				neo4j.ExecuteQueryWithDatabase(dbName))
			if err != nil {
				panic(err)
			}
		}
	}
}
