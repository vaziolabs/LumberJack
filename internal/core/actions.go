package core

import (
	"fmt"
	"time"
)

// GenerateID generates a unique ID for a user
func GenerateID() string {
	return fmt.Sprintf("user-%d", time.Now().UnixNano())
}

// StartEvent starts a new event or schedules it for the future
func (n *Node) StartEvent(eventID string, plannedStart, plannedEnd *time.Time, metadata map[string]interface{}) error {
	if n.Type != LeafNode {
		return fmt.Errorf("cannot add event to non-leaf node")
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	event := Event{
		Metadata: metadata,
		Status:   EventPending,
	}

	// Handle category if provided in metadata
	if category, ok := metadata["category"].(string); ok {
		event.Category = category
	}

	// Handle frequency if provided in metadata
	if frequency, ok := metadata["frequency"].(string); ok {
		event.Frequency = frequency
	}

	// Handle custom pattern if provided in metadata
	if pattern, ok := metadata["custom_pattern"].(string); ok {
		event.Pattern = pattern
	}

	if plannedStart == nil || time.Now().After(*plannedStart) {
		now := time.Now()
		event.StartTime = &now
		event.Status = EventOngoing
	}

	n.Events[eventID] = event
	return nil
}

// EndEvent marks an event as finished
func (n *Node) EndEvent(eventID string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	event, exists := n.Events[eventID]
	if !exists {
		return fmt.Errorf("event not found: %s", eventID)
	}

	if event.StartTime == nil {
		return fmt.Errorf("cannot end event that hasn't started")
	}

	now := time.Now()
	event.EndTime = &now
	event.Status = EventFinished
	n.Events[eventID] = event
	return nil
}

// AppendToEvent adds a new entry to an ongoing event
func (n *Node) AppendToEvent(eventID string, content interface{}, metadata map[string]interface{}, userID string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	event, exists := n.Events[eventID]
	if !exists {
		return fmt.Errorf("event not found: %s", eventID)
	}

	if event.EndTime != nil {
		return fmt.Errorf("cannot append to finished event")
	}

	if event.StartTime == nil {
		return fmt.Errorf("cannot append to event that hasn't started")
	}

	entry := Entry{
		Timestamp: time.Now(),
		Content:   content,
		Metadata:  metadata,
		UserID:    userID,
	}

	event.Entries = append(event.Entries, entry)
	n.Events[eventID] = event
	return nil
}

// PlanEvent plans a future event
func (n *Node) PlanEvent(eventID string, plannedStart, plannedEnd *time.Time, metadata map[string]interface{}) error {
	if n.Type != LeafNode {
		return fmt.Errorf("cannot plan event for non-leaf node")
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	event := Event{
		Metadata: metadata,
		Status:   EventPending,
	}

	event.StartTime = plannedStart
	event.EndTime = plannedEnd

	n.PlannedEvents[eventID] = event
	return nil
}

// CheckPermission checks if a user has permission to perform an action on the node
func (n *Node) CheckPermission(userID string, permission Permission) bool {
	for _, user := range n.Users {
		if user.ID == userID {
			for _, perm := range user.Permissions {

				if perm == permission {
					return true
				}
			}
		}
	}
	return false
}

// AssignUser assigns a user to the node with permission checking
func (n *Node) AssignUser(user User, permission Permission) error {
	// Check if the user already exists in the node
	for i, existingUser := range n.Users {
		if existingUser.ID == user.ID {
			for _, perm := range existingUser.Permissions {
				if perm == permission {
					return fmt.Errorf("user %s already has permission %d", user.Username, permission)
				}
			}
			// Add the new permission to the existing user's permissions
			n.Users[i].Permissions = append(n.Users[i].Permissions, permission)
			return nil
		}
	}

	// If the user does not exist, add them to the node with the specified permission
	user.Permissions = []Permission{permission}
	n.Users = append(n.Users, user)
	return nil
}

func (n *Node) StartTimeTracking(userID string) *Entry {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	entry := Entry{
		Timestamp: time.Now(),
		UserID:    userID,
		Content:   "start_time_entry",
	}

	n.Entries = append(n.Entries, entry)
	return &entry
}

// StopTimeTracking stops tracking time for the node
func (n *Node) StopTimeTracking(userID string) *Entry {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	entry := Entry{
		Timestamp: time.Now(),
		UserID:    userID,
		Content:   "end_time_entry",
	}

	n.Entries = append(n.Entries, entry)

	return &entry
}
