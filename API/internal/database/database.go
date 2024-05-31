package database


//using mongo

import (
	"context"
	
	"log"
	"time"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)





func OpenConnection() (*mongo.Client, error) {
	 
    //coonect
client, err := mongo.NewClient(options.Client().ApplyURI( "mongodb://127.0.0.1:27018/"))
if err != nil {
	log.Fatal(err)
}
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel() // Call the cancel function to avoid a context leak
err = client.Connect(ctx)
if err != nil {
	log.Fatal(err)
}
    
    log.Println("Connected to MongoDB!")
	return client, nil
}


func Close(client *mongo.Client) {
    err := client.Disconnect(context.Background())
    if err != nil {
        log.Println("Error closing connection to MongoDB:", err)
        return
    }
    log.Println("Connection to MongoDB closed.")
}

func SaveCollection(client *mongo.Client, collectionName string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	
	collection := client.Database("blockchain").Collection(collectionName)
	
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}


