package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Task struct {
    ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Title       string             `json:"title" bson:"title"`
    Description string             `json:"description" bson:"description"`
    Priority    string             `json:"priority" bson:"priority"`
    Status      string             `json:"status" bson:"status"`
    Deadline    time.Time          `json:"deadline" bson:"deadline"`
}
