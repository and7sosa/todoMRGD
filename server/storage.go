package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	collection *mongo.Collection
}

func ConnectDB(name string, uri string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	sOptions := options.ServerAPI(options.ServerAPIVersion1)
	cOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(sOptions)
	client, err := mongo.Connect(ctx, cOptions)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(name)
	cmd := bson.D{{Key: "ping", Value: 1}}
	err = db.RunCommand(ctx, cmd).Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database successfully pinged. MongoDB connected!")
	return db
}

func NewStorage(name string, db *mongo.Database) *Storage {
	collection := db.Collection(name)
	return &Storage{
		collection: collection,
	}
}

func (s *Storage) CreateTodo(todo *Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := s.collection.InsertOne(ctx, todo)
	if err != nil {
		return err
	}
	id := res.InsertedID.(primitive.ObjectID)
	todo.ID = id
	return nil
}

func (s *Storage) GetAllTodos(todos *[]Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	filter := bson.D{{}}
	defer cancel()
	cursor, err := s.collection.Find(ctx, filter)

	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var todo Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return err
		}
		*todos = append(*todos, todo)
	}
	return nil
}

func (s *Storage) GetTodoById(idHex string, todo *Todo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	res := s.collection.FindOne(ctx, filter)
	err = res.Decode(&todo)
	if err != nil {
		return err
	}
	return nil
}
func (s *Storage) UpdateTodo(idHex string, todo *Todo) error {
	// var oldTodo Todo
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	// todo.ID = id
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "task", Value: todo.Task},
		{Key: "completed", Value: todo.Completed},
	}}}
	res := s.collection.FindOneAndUpdate(ctx, filter, update)
	// err = res.Decode(&oldTodo)
	err = res.Decode(&todo)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteTodo(idHex string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: id}}
	res, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("error deleting todo item. Item with id %v not found", idHex)
	}
	return nil
}
