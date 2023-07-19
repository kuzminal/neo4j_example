package query

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"neo4j/internal/model"
	"neo4j/internal/util"
)

func FindFriendsUnderAge(ctx context.Context, driver neo4j.DriverWithContext, user string, age int) (friends []model.Person) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `
        MATCH (p:Person {name: $name})-[:KNOWS]-(friend:Person)
        WHERE friend.age < $age
        RETURN friend
        `, map[string]any{
		"name": user,
		"age":  age,
	}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))
	if err != nil {
		panic(err)
	}

	for _, record := range result.Records {
		person, _ := record.Get("friend")
		node := person.(dbtype.Node)
		p, _ := util.ConvertToPerson(node)
		friends = append(friends, p)
	}
	return friends
}

func FindFriendsOfFriends(ctx context.Context, driver neo4j.DriverWithContext, user string, dbName string) (friends []model.Person) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `
       MATCH (p:Person {name: "Peter"})-[:KNOWS]-(friend:Person)-[:KNOWS]-(ffriend:Person) 
       RETURN ffriend
        `, map[string]any{
		"name": user,
	}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbName))
	if err != nil {
		panic(err)
	}

	for _, record := range result.Records {
		person, _ := record.Get("ffriend")
		node := person.(dbtype.Node)
		p, _ := util.ConvertToPerson(node)
		friends = append(friends, p)
	}
	return friends
}
