package utils

import (
    "context"
    "encoding/json"
    "log"
    "github.com/go-redis/redis/v8"
    "task-management/models"
    "github.com/joho/godotenv"
    "time"
    "fmt"
    "os"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
    // Load environment variables from the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    // Directly set Redis address
    redisAddr := os.Getenv("REDIS_ADDR")

    // Initialize Redis client with the provided address
    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })

    // Optionally check if Redis is reachable (useful for debugging)
    _, err = redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
}

func GetTasksFromCache() ([]models.Task, error) {
    // Try to get the tasks from Redis cache
    val, err := redisClient.Get(ctx, "tasks").Result()
    if err == redis.Nil {
        log.Println("Cache miss: no tasks found in cache.") // Log cache miss for visibility
        // Cache miss, return nil with no error
        return nil, nil
    } else if err != nil {
        log.Printf("Error fetching tasks from cache: %v", err) // Log Redis error for debugging
        // Return Redis error
        return nil, fmt.Errorf("failed to fetch tasks from cache: %w", err)
    }

    // Unmarshal the cached JSON data into tasks
    var tasks []models.Task
    err = json.Unmarshal([]byte(val), &tasks)
    if err != nil {
        log.Printf("Error unmarshaling tasks from cache: %v", err) // Log unmarshaling error
        // Return unmarshaling error
        return nil, fmt.Errorf("failed to unmarshal tasks from cache: %w", err)
    }

    log.Println("Tasks successfully fetched from cache.") // Log success
    return tasks, nil
}

func SetTasksInCache(tasks []models.Task) error {
    data, err := json.Marshal(tasks)
    if err != nil {
        log.Printf("Error marshaling tasks to JSON: %v", err) 
        return fmt.Errorf("failed to marshal tasks for cache: %w", err) 
    }

    err = redisClient.Set(ctx, "tasks", data, time.Hour).Err()
    if err != nil {
        log.Printf("Error setting tasks in cache: %v", err)
        return fmt.Errorf("failed to set tasks in cache: %w", err) 
    }

    log.Println("Tasks successfully cached.") 
    return nil
}

func ClearCache() {
    redisClient.Del(ctx, "tasks")
}
