package main

import (
	"auth/api"
	"auth/database"
	"auth/domain/users"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	envVars := loadEnvVars()
	client := connectMongoDB(envVars)

	// Initialize the dependencies
	transactionRepository := database.NewTransactionRepository(client)
	userRepository := database.NewUserRepository(client)
	logRepository := database.NewLogRepository(client)

	userService := users.NewUserService(transactionRepository, logRepository, userRepository)

	appCtx := &api.AppContext{
		Port:        envVars.port,
		JwtSecret:   []byte(envVars.jwtSecret),
		UserService: userService,
	}

	api.Run(appCtx)
	fmt.Println("Hello, World!")
}

type envVars struct {
	dbUri        string
	dbReplicaSet string
	jwtSecret    string
	port         int
}

func loadEnvVars() *envVars {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	vars := os.Environ()
	envMap := make(map[string]string)

	for _, envVar := range vars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			envMap[key] = value
		}
	}

	lookup := func(key, defaultValue string) string {
		if value, ok := envMap[key]; ok {
			return value
		}
		return defaultValue
	}

	port, err := strconv.Atoi(lookup("PORT", "8080"))
	if err != nil {
		log.Fatalf("Error converting port to integer: %s", err)
	}

	return &envVars{
		dbUri:        lookup("DB_URI", "mongodb://localhost:27017/"),
		dbReplicaSet: lookup("DB_REPLICA_SET", "myReplicaSet"),
		jwtSecret:    lookup("JWT_SECRET", ""),
		port:         port,
	}

}

func connectMongoDB(envVars *envVars) *mongo.Client {
	clientOptions := options.Client().ApplyURI(envVars.dbUri)
	clientOptions.SetDirect(true)
	clientOptions.SetReplicaSet(envVars.dbReplicaSet)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	return client
}
