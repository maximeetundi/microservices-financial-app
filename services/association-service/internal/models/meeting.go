package models

import "time"

// MeetingStatus defines the status of a meeting
type MeetingStatus string

const (
	MeetingStatusScheduled MeetingStatus = "scheduled"
	MeetingStatusCompleted MeetingStatus = "completed"
	MeetingStatusCancelled MeetingStatus = "cancelled"
)

// Meeting represents an association meeting
type Meeting struct {
	ID            string        `json:"id"`
	AssociationID string        `json:"association_id"`
	Title         string        `json:"title"`
	Date          time.Time     `json:"date"`
	Location      string        `json:"location"`
	Agenda        JSONB         `json:"agenda"`
	Minutes       string        `json:"minutes"`
	Attendance    JSONB         `json:"attendance"` // user_id => present/absent
	Status        MeetingStatus `json:"status"`
	CreatedBy     string        `json:"created_by"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// CreateMeetingRequest is the request to schedule a meeting
type CreateMeetingRequest struct {
	Title    string    `json:"title" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
	Location string    `json:"location"`
	Agenda   JSONB     `json:"agenda"`
}

// RecordAttendanceRequest is the request to record attendance
type RecordAttendanceRequest struct {
	Attendance map[string]bool `json:"attendance" binding:"required"` // user_id => present
}

// UpdateMinutesRequest is the request to update meeting minutes
type UpdateMinutesRequest struct {
	Minutes string `json:"minutes" binding:"required"`
}
