package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"neo4j/internal/model"
	"neo4j/internal/mutation"
	"neo4j/internal/query"
)

func main() {
	ctx := context.Background()
	//dbName := "users"
	dbName := "neo4j"
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	// создание первоначального графа пользлователей
	//graph.CreateUsers(ctx, driver, dbName)

	// создание связей между пользаками
	//graph.CreateRelationships(ctx, driver, dbName)

	//поиск по графу
	userName := "Peter"
	age := 20
	friends := query.FindFriendsUnderAge(ctx, driver, userName, age)
	fmt.Printf("\n%s's friends uder %d age:\n", userName, age)
	for _, friend := range friends {
		fmt.Printf("Name: %s, age: %d\n", friend.Name, friend.Age)
	}
	// друзья друзей
	ffriends := query.FindFriendsOfFriends(ctx, driver, userName, dbName)
	fmt.Printf("\n%s's friends of friends:\n", userName)
	for _, friend := range ffriends {
		fmt.Printf("Name: %s, age: %d\n", friend.Name, friend.Age)
	}

	//создать пользователя
	createdUser, _ := mutation.CreateUser(model.Person{Name: "Alex", Age: 37}, ctx, driver, dbUser)

	// добавить друга пользователю
	mutation.AddFriends(model.Person{Name: "Peter"}, []model.Person{createdUser}, ctx, driver, dbName)
	mutation.AddFriends(model.Person{Name: "Alex"}, []model.Person{{Name: "Alice"}}, ctx, driver, dbName)

	//удалить связь
	mutation.RemoveRelationWithUser(model.Person{Name: "Alex"}, model.Person{Name: "Alice"}, ctx, driver, dbName)
}
