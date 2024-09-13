package utils

import (
    "context"
    "fmt"
    "log"
    "os"
	
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "task-management/models"
    "github.com/joho/godotenv"
)

var mongoClient *mongo.Client
var mongoCtx = context.Background()

func init() {
    // Load environment variables from the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Set MongoDB URI
    mongoURI := os.Getenv("MONGO_URI")

    // Initialize MongoDB client with the provided URI
    clientOptions := options.Client().ApplyURI(mongoURI)
    
    // Use '=' here because 'err' is already declared
    mongoClient, err = mongo.Connect(mongoCtx, clientOptions)
    if err != nil {
        fmt.Printf("Failed to connect to MongoDB: %s\n", err)
        return
    }

    // Check if MongoDB is reachable
    err = mongoClient.Ping(mongoCtx, nil)
    if err != nil {
        fmt.Printf("Failed to ping MongoDB: %s\n", err)
        return
    }

    fmt.Println("Successfully connected to MongoDB")
}

// Fetch all tasks from MongoDB
func GetAllTasksFromDB() ([]models.Task, error) {
    collection := mongoClient.Database("taskdb").Collection("tasks")
    
    var tasks []models.Task
    cursor, err := collection.Find(mongoCtx, bson.M{})
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve tasks from MongoDB: %w", err)
    }
    
    defer cursor.Close(mongoCtx)
    
    for cursor.Next(mongoCtx) {
        var task models.Task
        if err := cursor.Decode(&task); err != nil {
            return nil, fmt.Errorf("failed to decode task: %w", err)
        }
        tasks = append(tasks, task)
    }
    
    if err := cursor.Err(); err != nil {
        return nil, fmt.Errorf("cursor error: %w", err)
    }
    
    return tasks, nil
}

func CreateTaskInDB(task models.Task) error {
    collection := mongoClient.Database("taskdb").Collection("tasks")

    // Remove ID field if set
    task.ID = primitive.NilObjectID // Use nil ID to let MongoDB generate a new one

    _, err := collection.InsertOne(mongoCtx, task)
    if err != nil {
        return fmt.Errorf("failed to insert task into MongoDB: %w", err)
    }
    return nil
}


func UpdateTaskInDB(id string, task models.Task) error {
    collection := mongoClient.Database("taskdb").Collection("tasks")
    
    // Convert string ID to primitive.ObjectID
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid ID format: %w", err)
    }

    // Create filter to find the document with the specified ID
    filter := bson.M{"_id": objID}
    
    // Create update document
    update := bson.M{"$set": task}

    // Perform update operation
    result, err := collection.UpdateOne(mongoCtx, filter, update)
    if err != nil {
        return fmt.Errorf("failed to update task in MongoDB: %w", err)
    }

    // Check if the update was acknowledged
    if result.MatchedCount == 0 {
        return fmt.Errorf("no task found with ID: %s", id)
    }

    return nil
}


func DeleteTaskFromDB(id string) error {
    // Ensure the ID is not empty
    log.Printf("Received ID for deletion: %s", id)
    if id == "" {
        return fmt.Errorf("task ID cannot be empty")
    }

    // Define the MongoDB collection
    collection := mongoClient.Database("taskdb").Collection("tasks")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return fmt.Errorf("invalid ID format: %w", err)
    }
    // Create a filter to match the task by ID
    filter := bson.M{"_id": objID}
    
    // Attempt to delete the task
    result, err := collection.DeleteOne(mongoCtx, filter)
    if err != nil {
        // Log the error and return it
        log.Printf("Error deleting task from MongoDB: %v", err)
        return fmt.Errorf("failed to delete task from MongoDB: %w", err)
    }

    // Check if the task was actually deleted
    if result.DeletedCount == 0 {
        return fmt.Errorf("no task found with ID: %s", id)
    }

    return nil
}

func BulkCreateTasksInDB(tasks []models.Task) error {
    collection := mongoClient.Database("taskdb").Collection("tasks")
    var documents []interface{}

    for _, task := range tasks {
        task.ID = primitive.NilObjectID 

        documents = append(documents, task)
    }

    _, err := collection.InsertMany(mongoCtx, documents)
    if err != nil {
        return fmt.Errorf("failed to bulk insert tasks into MongoDB: %w", err)
    }
    return nil
}

