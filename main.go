package main

import (
	"context"
	"fmt"
	"log"

    "github.com/Umar-Md/go-apis/controllers"  // Correctly import controllers package
	"github.com/Umar-Md/go-apis/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontroller controllers.UserController
	ctx            context.Context
	userCollection *mongo.Collection
	mongoClient    *mongo.Client
	err            error
)

func init() {
	// Initialize context
	ctx = context.TODO()

	// MongoDB connection setup
	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection established")

	// Initialize MongoDB collection
	userCollection = mongoClient.Database("userdb").Collection("users")

	// Initialize the services and controllers
	userservice = services.NewUserService(userCollection, ctx)
	usercontroller = controllers.New(userservice)

	// Initialize Gin server
	server = gin.Default()
}

func main() {
	// Defer MongoDB disconnection
	defer mongoClient.Disconnect(ctx)

	// Define the base path for API versioning
	basepath := server.Group("/v1")

	// Register routes with the controller
	usercontroller.RegisterUserRoutes(basepath)

	// Start the server
	log.Fatal(server.Run(":9090"))
}
