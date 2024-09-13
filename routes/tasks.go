package routes

import (
    "os"
    "encoding/json"
    "net/http"
    "task-management/models"
    "task-management/utils"
    "github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

func SetupRouter() http.Handler {
    router := mux.NewRouter()

    // Define your routes
    router.HandleFunc("/tasks", GetTasks).Methods("GET")
    router.HandleFunc("/tasks", CreateTask).Methods("POST")
    router.HandleFunc("/tasks/{id:[a-zA-Z0-9]+}", UpdateTask).Methods("PUT")
    router.HandleFunc("/tasks/{id:[a-zA-Z0-9]+}", DeleteTask).Methods("DELETE")
    router.HandleFunc("/tasks/bulk", BulkCreateTasks).Methods("POST")

    // Get the frontend URL from the environment variable
    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:3000" // Fallback in case env var isn't set
    }

    // Configure CORS settings
    corsOptions := handlers.AllowedOrigins([]string{frontendURL})
    corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
    corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

    // Apply CORS middleware
    return handlers.CORS(corsOptions, corsHeaders, corsMethods)(router)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := utils.GetTasksFromCache()
    if err != nil {
        utils.RespondWithError(w, "Unable to fetch tasks: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if tasks == nil {
        // Cache miss, retrieve from database
        tasks, err = utils.GetAllTasksFromDB()
        if err != nil {
            utils.RespondWithError(w, "Unable to fetch tasks from database: "+err.Error(), http.StatusInternalServerError)
            return
        }
        // Update the cache
        utils.SetTasksInCache(tasks)
    }
    // utils.SetTasksInCache(tasks)
    utils.RespondWithSuccess(w, "Tasks fetched successfully", http.StatusOK, tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
    var task models.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        utils.RespondWithError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    err := utils.CreateTaskInDB(task)
    if err != nil {
        utils.RespondWithError(w, "Unable to create task: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Invalidate cache and update it with the latest data
    tasks, err := utils.GetAllTasksFromDB()
    if err != nil {
        utils.RespondWithError(w, "Failed to fetch tasks after creation", http.StatusInternalServerError)
        return
    }
    utils.SetTasksInCache(tasks)

    utils.RespondWithSuccess(w, "Task created successfully", http.StatusCreated, nil)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var task models.Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        utils.RespondWithError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    err := utils.UpdateTaskInDB(id, task)
    if err != nil {
        utils.RespondWithError(w, "Unable to update task: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Invalidate cache and update it with the latest data
    tasks, err := utils.GetAllTasksFromDB()
    if err != nil {
        utils.RespondWithError(w, "Failed to fetch tasks after update", http.StatusInternalServerError)
        return
    }
    utils.SetTasksInCache(tasks)

    utils.RespondWithSuccess(w, "Task updated successfully", http.StatusOK, tasks)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    err := utils.DeleteTaskFromDB(id)
    if err != nil {
        utils.RespondWithError(w, "Unable to delete task: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Invalidate cache and update it with the latest data
    tasks, err := utils.GetAllTasksFromDB()
    if err != nil {
        utils.RespondWithError(w, "Failed to fetch tasks after deletion", http.StatusInternalServerError)
        return
    }
    utils.SetTasksInCache(tasks)

    utils.RespondWithSuccess(w, "Task deleted successfully", http.StatusOK, nil)
}

func BulkCreateTasks(w http.ResponseWriter, r *http.Request) {
    var tasks []models.Task
    if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
        utils.RespondWithError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
        return
    }

    err := utils.BulkCreateTasksInDB(tasks)
    
    if err != nil {
        utils.RespondWithError(w, "Unable to create tasks: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Invalidate cache and update it with the latest data
    allTasks, err := utils.GetAllTasksFromDB()
    if err != nil {
        utils.RespondWithError(w, "Failed to fetch tasks after bulk creation", http.StatusInternalServerError)
        return
    }
    utils.SetTasksInCache(allTasks)

    utils.RespondWithSuccess(w, "Tasks created successfully", http.StatusCreated, nil)
}
