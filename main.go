package main

import (
	"context"
	"fmt"
	"log"

    "github.com/Umar-Md/go-apis/controllers"  
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
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection established")

	userCollection = mongoClient.Database("userdb").Collection("users")

	userservice = services.NewUserService(userCollection, ctx)
	usercontroller = controllers.New(userservice)


	server = gin.Default()
}

func main() {

	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/v1")


	usercontroller.RegisterUserRoutes(basepath)


	log.Fatal(server.Run(":9090"))
}
