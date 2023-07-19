package mutation

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"neo4j/internal/model"
	"neo4j/internal/util"
)

func AddFriends(user model.Person, friends []model.Person, ctx context.Context, driver neo4j.DriverWithContext, dbName string) error {
	for _, friend := range friends {
		_, err := neo4j.ExecuteQuery(ctx, driver, `
                MATCH (p:Person {name: $person})
                MATCH (friend:Person {name: $friend})
                MERGE (p)-[:KNOWS]->(friend)
                `, map[string]any{
			"person": user.Name,
			"friend": friend.Name,
		}, neo4j.EagerResultTransformer,
			neo4j.ExecuteQueryWithDatabase(dbName))
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateUser(user model.Person, ctx context.Context, driver neo4j.DriverWithContext, dbName string) (person model.Person, err error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		"MERGE (user:Person {name: $personName, age: $personAge}) RETURN user",
		map[string]any{
			"personName": user.Name,
			"personAge":  user.Age,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbName))
	if err != nil {
		panic(err)
	}
	if result.Records[0] == nil {
		return person, fmt.Errorf("could not create user %s", user.Name)
	}
	record := result.Records[0]
	p, _ := record.Get("user")
	node := p.(dbtype.Node)
	person, _ = util.ConvertToPerson(node)

	return person, nil
}

func RemoveRelationWithUser(user model.Person, exFriend model.Person, ctx context.Context, driver neo4j.DriverWithContext, dbName string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		"MATCH (u:Person {name: $userName})-[r:KNOWS]->(e:Person {name: $exFriend}) DELETE r",
		map[string]any{
			"userName": user.Name,
			"exFriend": exFriend.Name,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbName))
	if err != nil {
		return err
	}
	return nil
}
