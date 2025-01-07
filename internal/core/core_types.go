package core

import (
	"sync"
	"time"
)

// Permission represents a permission level for users
type Permission int

// LeafType represents the type of a leaf in the tree-forest
type LeafType int

// EventStatus represents the current status of a timed event
type EventStatus string

// Admin represents an admin user
type Admin struct {
	User
	Permissions []Permission `json:"permissions"`
}

// User represents a user in the system
type User struct {
	ID           string       `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Organization string       `json:"organization"`
	Phone        string       `json:"phone"`
	Permissions  []Permission `json:"permissions"`
}

// Event represents an event with start/end times and entries
type Event struct {
	StartTime *time.Time             `json:"start_time,omitempty"`
	EndTime   *time.Time             `json:"end_time,omitempty"`
	Entries   []Entry                `json:"entries"`
	Metadata  map[string]interface{} `json:"metadata"`
	Status    EventStatus            `json:"status"`
	Category  string                 `json:"category,omitempty"`
	Frequency string                 `json:"frequency,omitempty"`
	Pattern   string                 `json:"pattern,omitempty"`
}

// EventSummary provides a summary of the event's timing and status
type EventSummary struct {
	Status         EventStatus `json:"status"`
	Duration       *string     `json:"duration,omitempty"`
	RemainingTime  *string     `json:"remaining_time,omitempty"`
	EntriesCount   int         `json:"entries_count"`
	LastUpdateTime *time.Time  `json:"last_update_time,omitempty"`
}

// NodeType represents the type of a node in the tree-forest
type NodeType int

// Entry represents an entry in the node
type Entry struct {
	Content   interface{}            `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	UserID    string                 `json:"user_id"`
	Timestamp time.Time              `json:"timestamp"`
}

// Node represents a node in the tree-forest
type Node struct {
	ID            string            `json:"id"`
	Type          NodeType          `json:"type"`
	Name          string            `json:"name"`
	Parents       map[string]string `json:"parents"`
	Children      map[string]*Node  `json:"children"`
	Events        map[string]Event  `json:"events"`
	PlannedEvents map[string]Event  `json:"planned_events"`
	Users         []User            `json:"users"`
	Entries       []Entry           `json:"entries"`
	mutex         sync.RWMutex      `json:"-"`
}
