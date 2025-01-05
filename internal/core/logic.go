package core

import (
	"fmt"
	"reflect"

	"golang.org/x/crypto/bcrypt"
)

// CompareEvents compares the planned event to the actual event and reports differences
func (n *Node) CompareEvents(plannedEventID, actualEventID string) (bool, error) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	plannedEvent, plannedExists := n.PlannedEvents[plannedEventID]
	actualEvent, actualExists := n.Events[actualEventID]

	if !plannedExists || !actualExists {
		return false, fmt.Errorf("one or both events not found: plannedEventID=%s, actualEventID=%s", plannedEventID, actualEventID)
	}

	// Initialize a slice to hold differences
	var differences []string

	// Compare StartTime
	if (plannedEvent.StartTime == nil && actualEvent.StartTime != nil) || (plannedEvent.StartTime != nil && actualEvent.StartTime == nil) {
		differences = append(differences, "StartTime differs")
	} else if plannedEvent.StartTime != nil && actualEvent.StartTime != nil && !plannedEvent.StartTime.Equal(*actualEvent.StartTime) {
		differences = append(differences, fmt.Sprintf("StartTime differs: planned=%v, actual=%v", *plannedEvent.StartTime, *actualEvent.StartTime))
	}

	// Compare EndTime
	if (plannedEvent.EndTime == nil && actualEvent.EndTime != nil) || (plannedEvent.EndTime != nil && actualEvent.EndTime == nil) {
		differences = append(differences, "EndTime differs")
	} else if plannedEvent.EndTime != nil && actualEvent.EndTime != nil && !plannedEvent.EndTime.Equal(*actualEvent.EndTime) {
		differences = append(differences, fmt.Sprintf("EndTime differs: planned=%v, actual=%v", *plannedEvent.EndTime, *actualEvent.EndTime))
	}

	// Compare Status
	if plannedEvent.Status != actualEvent.Status {
		differences = append(differences, fmt.Sprintf("Status differs: planned=%s, actual=%s", plannedEvent.Status, actualEvent.Status))
	}

	// Compare Metadata
	if !reflect.DeepEqual(plannedEvent.Metadata, actualEvent.Metadata) {
		differences = append(differences, "Metadata differs")
	}

	// Compare Entries
	if len(plannedEvent.Entries) != len(actualEvent.Entries) {
		differences = append(differences, fmt.Sprintf("Entries count differs: planned=%d, actual=%d", len(plannedEvent.Entries), len(actualEvent.Entries)))
	} else {
		for i := range plannedEvent.Entries {
			if !reflect.DeepEqual(plannedEvent.Entries[i], actualEvent.Entries[i]) {
				differences = append(differences, fmt.Sprintf("Entry %d differs", i))
			}
		}
	}

	// If there are differences, return them
	if len(differences) > 0 {
		return false, fmt.Errorf("differences found: %v", differences)
	}

	return true, nil
}

// SetPassword sets the password for the user
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword verifies the password for the user
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
