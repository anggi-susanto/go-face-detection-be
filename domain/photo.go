package domain

import "time"

type Photo struct {
	ID            int64     `json:"id" bson:"_id, omitempty"`
	FilePath      string    `json:"photo_url" bson:"photo_url"`
	TimeStamp     time.Time `json:"timestamp" bson:"timestamp"`
	Status        string    `json:"status" bson:"status"`
	FacesDetected int       `json:"faces_detected" bson:"faces_detected"`
}
