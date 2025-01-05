package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"forestree"
)

var (
	testStateFile string // Global state file for all tests
)

func TestMain(m *testing.M) {
	// Setup test environment
	tmpDir, err := os.MkdirTemp("", "forestree-test-*")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	testStateFile = filepath.Join(tmpDir, "test_state.dat")

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll(tmpDir)
	os.Exit(code)
}

func setupTestForest(t *testing.T) *App {
	app := NewApp()
	logger := newTestLogger(t)
	logger.info("Setting up test forest")
	root := app.Forest

	// Create basic structure with admin user
	root.ID = "root"
	root.Name = "root"
	root.Type = forestree.BranchNode
	root.Children = make(map[string]*forestree.Node)

	adminUser := forestree.User{
		ID:          "admin",
		Username:    "admin",
		Permissions: []forestree.Permission{forestree.AdminPermission},
	}
	root.Users = []forestree.User{adminUser}

	// Create test node with admin permissions
	testNode := forestree.NewNode(forestree.LeafNode, "test-node")
	testNode.ID = "test-node"
	testNode.Users = []forestree.User{adminUser}

	// Add child and set up parent reference
	root.AddChild(testNode)

	// Save initial state
	if err := app.WriteChangesToFile(app.Forest, testStateFile); err != nil {
		logger.failure("Failed to write initial state: %v", err)
	} else {
		logger.success("Saved initial state to file")
	}

	logger.info("Test forest setup complete")
	return app
}

func TestForestOperations(t *testing.T) {
	app := setupTestForest(t)
	logger := newTestLogger(t)
	logger.enter("ForestOperations")
	defer logger.exit("ForestOperations")

	rootNode := app.Forest

	logger.enter("Node Creation")
	childNode1 := forestree.NewNode(forestree.BranchNode, "child1")
	childNode1.ID = "child1"
	childNode2 := forestree.NewNode(forestree.LeafNode, "child2")
	childNode2.ID = "child2"
	childNode3 := forestree.NewNode(forestree.LeafNode, "child3")
	childNode3.ID = "child3"

	rootNode.Children[childNode1.ID] = childNode1
	rootNode.Children[childNode2.ID] = childNode2
	rootNode.Children[childNode3.ID] = childNode3

	if err := app.WriteChangesToFile(rootNode, testStateFile); err != nil {
		logger.failure("Failed to write state to file: %v", err)
	} else {
		logger.success("State saved to file")
	}
	logger.exit("Node Creation")

	logger.info("Creating Nodes")

	logger.info("Setting Children Nodes")
	rootNode.Children["child1"] = childNode1
	rootNode.Children["child2"] = childNode2
	rootNode.Children["child3"] = childNode3

	// Verify the node structure
	if len(rootNode.Children) != 4 {
		logger.failure("Expected 4 children, got %d", len(rootNode.Children))
		t.Errorf("Expected 4 children, got %d", len(rootNode.Children))
	} else {
		logger.success("Found 4 children")
		// Log the children
		for _, child := range rootNode.Children {
			logger.info("Child: %s", child.Name)
		}
	}

	logger.info("Adding Users")
	// Test user assignment and permission checking
	user1 := forestree.User{ID: "user1", Permissions: []forestree.Permission{forestree.ReadPermission, forestree.WritePermission}}
	user2 := forestree.User{ID: "user2", Permissions: []forestree.Permission{forestree.ReadPermission}}
	user3 := forestree.User{ID: "user3", Permissions: []forestree.Permission{forestree.AdminPermission}}

	if err := childNode2.AssignUser(user1, forestree.ReadPermission); err != nil {
		logger.failure("Failed to assign user1 to childNode2: %v", err)
		t.Error(err)
	} else {
		logger.success("User1 assigned to childNode2")
	}
	if err := childNode2.AssignUser(user2, forestree.ReadPermission); err != nil {
		logger.failure("Failed to assign user2 to childNode2: %v", err)
		t.Error(err)
	} else {
		logger.success("User2 assigned to childNode2")
	}
	if err := childNode3.AssignUser(user3, forestree.AdminPermission); err != nil {
		logger.failure("Failed to assign user3 to childNode3: %v", err)
		t.Error(err)
	} else {
		logger.success("User3 assigned to childNode3")
	}
	logger.exit("User Management")

	// Check permissions
	if !childNode2.CheckPermission("user1", forestree.ReadPermission) {
		logger.failure("user1 should have read permission on childNode2")
		t.Errorf("user1 should have read permission on childNode2")
	} else if !childNode3.CheckPermission("user3", forestree.AdminPermission) {
		logger.failure("user3 should have admin permission on childNode3")
		t.Errorf("user3 should have admin permission on childNode3")
	} else if childNode2.CheckPermission("user2", forestree.WritePermission) {
		logger.failure("user2 should not have write permission on childNode2")
		t.Errorf("user2 should not have write permission on childNode2")
	} else {
		logger.success("Permission checks passed successfully")
	}

	logger.enter("Event Management")
	logger.info("Starting Event 1")
	// Test event management
	now := time.Now()
	event1ID := "event1"
	if err := childNode2.StartEvent(event1ID, &now, nil, map[string]interface{}{"title": "Test Event 1"}); err != nil {
		logger.failure("Failed to start event: %v", err)
		t.Error(err)
	} else {
		logger.success("Event started successfully")
	}

	logger.info("Appending Entry to Event 1")
	if err := childNode2.AppendToEvent(event1ID, "Event entry 1", map[string]interface{}{"note": "First entry"}, "user1"); err != nil {
		logger.failure("Failed to append to event: %v", err)
		t.Errorf("Failed to append to event: %v", err)
	} else {
		logger.success("Appended entry to event")
	}

	logger.info("Ending Event 1")
	if err := childNode2.EndEvent(event1ID); err != nil {
		logger.failure("Failed to end event: %v", err)
		t.Errorf("Failed to end event: %v", err)
	} else {
		logger.success("Event ended successfully")
	}

	// Verify event storage and persistence
	originalEvent, exists := childNode2.Events[event1ID]
	if !exists {
		logger.failure("Event %s not found in storage", event1ID)
		t.Errorf("Event %s not found in storage", event1ID)
	} else {
		logger.success("Event found in storage")
	}

	if err := app.WriteChangesToFile(rootNode, testStateFile); err != nil {
		logger.failure("Failed to write state to file: %v", err)
	} else {
		logger.success("State saved to file")
	}

	// Get the stored event after reload
	storedEvent, exists := childNode2.Events[event1ID]
	if !exists {
		logger.failure("Event %s not found in storage after reload", event1ID)
		t.Errorf("Event %s not found in storage after reload", event1ID)
	} else {
		logger.success("Event found in storage after reload")
	}

	// Verify event properties against the original stored version
	if !storedEvent.StartTime.Equal(*originalEvent.StartTime) {
		logger.failure("Start time mismatch after storage: expected %v, got %v",
			originalEvent.StartTime, storedEvent.StartTime)
		t.Errorf("Start time mismatch after storage: expected %v, got %v",
			originalEvent.StartTime, storedEvent.StartTime)
	} else {
		logger.success("Event start time is correct")
	}

	logger.info("Getting Event Summary")
	// Get the event summary
	eventSummary, err := childNode2.GetEventSummary(event1ID)
	if err != nil {
		logger.failure("Failed to get event summary: %v", err)
		t.Errorf("Failed to get event summary: %v", err)
	} else if eventSummary.Status != forestree.EventFinished {
		logger.failure("Expected event status to be 'finished', got '%s'", eventSummary.Status)
		t.Errorf("Expected event status to be 'finished', got '%s'", eventSummary.Status)
	} else if eventSummary.EntriesCount != 1 {
		logger.failure("Expected 1 entry in the event, got %d", eventSummary.EntriesCount)
		t.Errorf("Expected 1 entry in the event, got %d", eventSummary.EntriesCount)
	} else {
		logger.success("Event summary is correct")
	}

	// Test planned events
	event2ID := "event2"
	plannedStart := now.Add(time.Hour)
	plannedEnd := now.Add(2 * time.Hour)
	if err := childNode2.PlanEvent(event2ID, &plannedStart, &plannedEnd, map[string]interface{}{"title": "Test Event 2"}); err != nil {
		logger.failure("Failed to plan event: %v", err)
		t.Errorf("Failed to plan event: %v", err)
	} else {
		logger.success("Planned event successfully")
	}

	// Compare planned vs actual events
	match, err := childNode2.CompareEvents(event2ID, event1ID)
	if match {
		logger.failure("Expected events to not match, but they did")
		t.Errorf("Expected events to not match, but they did")
	} else {
		logger.success("Found expected differences: %s", err)
	}

	// Test time tracking
	start_entry := childNode2.StartTimeTracking("user1")
	if start_entry == nil {
		logger.failure("Failed to start time tracking")
		t.Error("Failed to start time tracking")
	} else {
		logger.success("Time tracking started")
	}

	end_entry := childNode2.StopTimeTracking("user1")
	if end_entry == nil {
		logger.failure("Failed to stop time tracking")
		t.Error("Failed to stop time tracking")
	} else {
		logger.success("Time tracking stopped")
	}

	summary := childNode2.GetTimeTrackingSummary("user1")
	if len(summary) != 1 {
		logger.failure("Expected 1 time tracking entry, got %d: %s", len(summary), summary)
		t.Errorf("Expected 1 time tracking entry, got %d: %s", len(summary), summary)
	} else {
		logger.success("Time tracking summary is correct")
	}
	logger.exit("Time Tracking")

	logger.enter("Persistence Verification")
	if err := app.WriteChangesToFile(rootNode, testStateFile); err != nil {
		logger.failure("Failed to save final state: %v", err)
		t.Error(err)
	} else {
		logger.success("State saved successfully")
	}

	// Create new app and load state
	newApp := NewApp()
	var loadedForest forestree.Node
	if err := newApp.ReadChangesFromFile(testStateFile, &loadedForest); err != nil {
		t.Fatalf("Failed to load state from file: %v", err)
	}
	newApp.forest = &loadedForest

	// Verify the loaded forest
	if _, err := newApp.forest.GetNode(childNode2.ID); err != nil {
		t.Errorf("Failed to find childNode2 in the loaded forest: %v", err)
	} else {
		logger.success("Found childNode2 in the loaded forest")
	}
	logger.exit("Persistence Verification")
}

func TestJsonSerialization(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("JsonSerialization")
	defer logger.exit("JsonSerialization")

	logger.enter("Node Setup")
	node := forestree.NewNode(forestree.BranchNode, "test-node")
	node.Users = []forestree.User{
		{ID: "user1", Username: "User 1", Email: "user1@example.com", Permissions: []forestree.Permission{forestree.ReadPermission, forestree.WritePermission}},
		{ID: "user2", Username: "User 2", Email: "user2@example.com", Permissions: []forestree.Permission{forestree.ReadPermission}},
	}

	// Add some events
	now := time.Now()
	one_hour := now.Add(time.Hour)
	node.Events = map[string]forestree.Event{
		"event1": {
			StartTime: &now,
			EndTime:   &one_hour,
			Entries: []forestree.Entry{
				{Content: "Entry 1", Timestamp: now, UserID: "user1"},
				{Content: "Entry 2", Timestamp: now.Add(10 * time.Minute), UserID: "user2"},
			},
			Metadata: map[string]interface{}{"title": "Test Event 1"},
			Status:   forestree.EventFinished,
		},
	}
	logger.exit("Node Setup")

	logger.enter("Serialization")
	jsonData, err := json.Marshal(node)
	if err != nil {
		logger.failure("Failed to marshal node: %v", err)
		t.Error(err)
	} else {
		logger.success("Node marshalled successfully")
	}

	// Deserialize the node from JSON
	var loadedNode forestree.Node
	if err := json.Unmarshal(jsonData, &loadedNode); err != nil {
		logger.failure("Failed to unmarshal node from JSON: %v", err)
		t.Error(err)
	}

	// Verify the deserialized node
	if loadedNode.Name != "test-node" {
		logger.failure("Expected node name to be 'test-node', got '%s'", loadedNode.Name)
		t.Errorf("Expected node name to be 'test-node', got '%s'", loadedNode.Name)
	} else {
		logger.success("Node name is correct")
	}
	if len(loadedNode.Users) != 2 {
		logger.failure("Expected 2 users, got %d", len(loadedNode.Users))
		t.Errorf("Expected 2 users, got %d", len(loadedNode.Users))
	} else {
		logger.success("Found 2 users")
	}
	if len(loadedNode.Events) != 1 {
		logger.failure("Expected 1 event, got %d", len(loadedNode.Events))
		t.Errorf("Expected 1 event, got %d", len(loadedNode.Events))
	} else {
		logger.success("Found 1 event")
	}
	event, ok := loadedNode.Events["event1"]
	if !ok {
		logger.failure("Expected event with ID 'event1' to be present")
		t.Errorf("Expected event with ID 'event1' to be present")
	} else {
		logger.success("Event found")
	}
	if event.Entries[0].Content != "Entry 1" {
		logger.failure("Expected first event entry to have content 'Entry 1', got '%s'", event.Entries[0].Content)
		t.Errorf("Expected first event entry to have content 'Entry 1', got '%s'", event.Entries[0].Content)
	} else {
		logger.success("Event entry content is correct")
	}
	logger.exit("Serialization")
}

func TestHandleAssignUser(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleAssignUser")
	defer logger.exit("HandleAssignUser")

	app := setupTestForest(t)

	body := map[string]interface{}{
		"path":        "test-node",
		"assignee_id": "test_user",
		"permission":  forestree.WritePermission,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/assign_user", bytes.NewBuffer(bodyBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handleAssignUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v\nBody: %v",
			status, http.StatusOK, rr.Body.String())
		t.Errorf("Handler returned wrong status code: got %v want %v\nBody: %v",
			status, http.StatusOK, rr.Body.String())
	} else {
		logger.success("User assigned successfully")
	}

	// Save state
	if err := app.WriteChangesToFile(app.Forest, testStateFile); err != nil {
		logger.failure("Failed to save state: %v", err)
		t.Fatalf("Failed to save state: %v", err)
	} else {
		logger.success("State saved successfully")
	}
}

func TestHandleGetTimeTracking(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleGetTimeTracking")
	defer logger.exit("HandleGetTimeTracking")

	app := setupTestForest(t)

	// Start time tracking
	startBody := map[string]interface{}{
		"path": "test-node",
	}
	startBytes, _ := json.Marshal(startBody)
	startReq := httptest.NewRequest("POST", "/start_time_tracking", bytes.NewBuffer(startBytes))
	startReq.Header.Set("X-User-ID", "admin")
	startReq.Header.Set("Content-Type", "application/json")
	app.handleStartTimeTracking(httptest.NewRecorder(), startReq)

	// Stop time tracking
	stopReq := httptest.NewRequest("POST", "/stop_time_tracking", bytes.NewBuffer(startBytes))
	stopReq.Header.Set("X-User-ID", "admin")
	stopReq.Header.Set("Content-Type", "application/json")
	app.handleStopTimeTracking(httptest.NewRecorder(), stopReq)

	// Get summary
	req := httptest.NewRequest("POST", "/get_time_tracking", bytes.NewBuffer(startBytes))
	req.Header.Set("X-User-ID", "admin")
	rr := httptest.NewRecorder()
	app.handleGetTimeTracking(rr, req)

	var summary []map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&summary)

	if len(summary) == 0 {
		logger.failure("Expected non-empty summary")
		t.Error("Expected non-empty summary")
	} else {
		logger.success("Time tracking summary is correct")
		logger.info("Summary: %v", summary)
	}
}

func TestHandleGetEventEntries(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleGetEventEntries")
	defer logger.exit("HandleGetEventEntries")

	app := setupTestForest(t)

	// First create an event
	startBody := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
		"metadata": map[string]interface{}{
			"test": "data",
		},
	}
	startBytes, _ := json.Marshal(startBody)
	startReq := httptest.NewRequest("POST", "/start_event", bytes.NewBuffer(startBytes))
	startReq.Header.Set("X-User-ID", "admin")
	startReq.Header.Set("Content-Type", "application/json")

	startRR := httptest.NewRecorder()
	app.handleStartEvent(startRR, startReq)

	if startRR.Code != http.StatusOK {
		logger.failure("Failed to start event: %v", startRR.Body.String())
		t.Fatalf("Failed to start event: %v", startRR.Body.String())
	}

	// Now get the entries
	body := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/get_event_entries", bytes.NewBuffer(bodyBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.handleGetEventEntries(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var entries []interface{}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		logger.failure("Failed to decode response: %v", err)
		t.Errorf("Failed to decode response: %v", err)
	} else {
		logger.success("Successfully retrieved event entries")
		logger.info("Entries: %v", entries)
	}
}

func TestHandleEndEvent(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleEndEvent")
	defer logger.exit("HandleEndEvent")

	app := setupTestForest(t)

	// Start event
	startBody := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
		"metadata": map[string]interface{}{"test": "data"},
	}
	startBytes, _ := json.Marshal(startBody)
	startReq := httptest.NewRequest("POST", "/start_event", bytes.NewBuffer(startBytes))
	startReq.Header.Set("X-User-ID", "admin")
	startReq.Header.Set("Content-Type", "application/json")
	app.handleStartEvent(httptest.NewRecorder(), startReq)

	// End event
	endBody := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
	}
	endBytes, _ := json.Marshal(endBody)
	endReq := httptest.NewRequest("POST", "/end_event", bytes.NewBuffer(endBytes))
	endReq.Header.Set("X-User-ID", "admin")
	endReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.handleEndEvent(rr, endReq)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("Event ended successfully")
		logger.info("Response: %s", rr.Body.String())
	}
}

func TestHandleAppendToEvent(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleAppendToEvent")
	defer logger.exit("HandleAppendToEvent")

	app := setupTestForest(t)

	// First start an event
	startBody := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
		"metadata": map[string]interface{}{
			"test": "data",
		},
	}
	startBytes, _ := json.Marshal(startBody)
	startReq := httptest.NewRequest("POST", "/start_event", bytes.NewBuffer(startBytes))
	startReq.Header.Set("X-User-ID", "admin")
	startReq.Header.Set("Content-Type", "application/json")

	startRR := httptest.NewRecorder()
	app.handleStartEvent(startRR, startReq)

	if startRR.Code != http.StatusOK {
		logger.failure("Failed to start event: %v", startRR.Body.String())
		t.Fatalf("Failed to start event: %v", startRR.Body.String())
	}

	// Now append to the event
	body := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
		"content":  "test content",
		"metadata": map[string]interface{}{
			"test": "data",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/append_event", bytes.NewBuffer(bodyBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.handleAppendToEvent(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v, body: %v",
			status, http.StatusOK, rr.Body.String())
		t.Errorf("Handler returned wrong status code: got %v want %v, body: %v",
			status, http.StatusOK, rr.Body.String())
	} else {
		logger.success("Successfully appended to event")
	}

	// Verify the entry was added
	getBody := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
	}
	getBytes, _ := json.Marshal(getBody)
	getReq := httptest.NewRequest("POST", "/get_event_entries", bytes.NewBuffer(getBytes))
	getReq.Header.Set("X-User-ID", "admin")
	getReq.Header.Set("Content-Type", "application/json")

	getRR := httptest.NewRecorder()
	app.handleGetEventEntries(getRR, getReq)

	var entries []interface{}
	if err := json.NewDecoder(getRR.Body).Decode(&entries); err != nil {
		logger.failure("Failed to decode entries: %v", err)
		t.Errorf("Failed to decode entries: %v", err)
	}

	if len(entries) != 1 {
		logger.failure("Expected 1 entry, got %d", len(entries))
		t.Errorf("Expected 1 entry, got %d", len(entries))
	} else {
		logger.success("Verified correct number of entries")
		logger.info("Entries: %v", entries)
	}
}

func TestHandleStartEvent(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleStartEvent")
	defer logger.exit("HandleStartEvent")

	app := setupTestForest(t)

	body := map[string]interface{}{
		"path":     "test-node",
		"event_id": "test-event",
		"metadata": map[string]interface{}{
			"test": "data",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/start_event", bytes.NewBuffer(bodyBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	app.handleStartEvent(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("Event started successfully")
		logger.info("Response: %s", rr.Body.String())
	}
}

func TestHandleAppendToEventWithCategories(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandleAppendToEventWithCategories")
	defer logger.exit("HandleAppendToEventWithCategories")

	app := setupTestForest(t)
	nodePath := "health/exercise/workout"
	logger.info("Looking for node at path: %s", nodePath)

	// Create branch structure: health::exercise::workout
	healthNode := forestree.NewNode(forestree.BranchNode, "health")
	healthNode.ID = "health"
	exerciseNode := forestree.NewNode(forestree.BranchNode, "exercise")
	exerciseNode.ID = "exercise"
	workoutNode := forestree.NewNode(forestree.LeafNode, "workout")
	workoutNode.ID = "workout"

	// Add admin user to all nodes
	adminUser := forestree.User{
		ID:          "admin",
		Username:    "admin",
		Permissions: []forestree.Permission{forestree.AdminPermission},
	}
	healthNode.Users = []forestree.User{adminUser}
	exerciseNode.Users = []forestree.User{adminUser}
	workoutNode.Users = []forestree.User{adminUser}

	// Set up node hierarchy
	app.Forest.AddChild(healthNode)
	healthNode.AddChild(exerciseNode)
	exerciseNode.AddChild(workoutNode)

	nodePath := "health/exercise/workout"
	logger.info("Testing with node path: %s", nodePath)

	// Start event
	startBody := map[string]interface{}{
		"path":     nodePath,
		"event_id": "morning_run",
		"metadata": map[string]interface{}{
			"type":           "event",
			"frequency":      "weekly",
			"custom_pattern": "MWF@0700",
			"category":       "health::exercise::workout",
		},
	}
	startBytes, _ := json.Marshal(startBody)
	startReq := httptest.NewRequest("POST", "/start_event", bytes.NewBuffer(startBytes))
	startReq.Header.Set("X-User-ID", "admin")
	startReq.Header.Set("Content-Type", "application/json")
	app.handleStartEvent(httptest.NewRecorder(), startReq)

	// Append to event
	appendBody := map[string]interface{}{
		"path":     nodePath,
		"event_id": "morning_run",
		"content":  "Completed morning run",
		"metadata": map[string]interface{}{
			"type":     "event",
			"status":   "completed",
			"category": "health::exercise::workout",
		},
	}
	appendBytes, _ := json.Marshal(appendBody)
	req := httptest.NewRequest("POST", "/append_event", bytes.NewBuffer(appendBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	app.handleAppendToEvent(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("Event appended successfully")
		logger.info("Response: %s", rr.Body.String())
	}
}

func TestHandlePlanEvent(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("HandlePlanEvent")
	defer logger.exit("HandlePlanEvent")

	app := setupTestForest(t)

	logger.info("Setting up test environment")

	// Create branch structure: study::languages::spanish
	studyNode := forestree.NewNode(forestree.BranchNode, "study")
	languagesNode := forestree.NewNode(forestree.BranchNode, "languages")
	spanishNode := forestree.NewNode(forestree.LeafNode, "spanish")

	app.Forest.Children["study"] = studyNode
	studyNode.Children["languages"] = languagesNode
	languagesNode.Children["spanish"] = spanishNode

	// Plan a study routine
	start := time.Now().Add(time.Hour)
	end := time.Now().Add(2 * time.Hour)
	body := map[string]interface{}{
		"path":       "study/languages/spanish",
		"event_id":   "vocabulary_practice",
		"start_time": start.Format(time.RFC3339),
		"end_time":   end.Format(time.RFC3339),
		"metadata": map[string]interface{}{
			"frequency":      "daily",
			"custom_pattern": "2000",
			"category":       "study::languages::spanish",
			"message":        "Plan vocabulary practice",
			"status":         "pending",
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/plan_event", bytes.NewBuffer(bodyBytes))
	req.Header.Set("X-User-ID", "admin")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.handlePlanEvent)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		logger.failure("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("Event planned successfully")
		logger.info("Response: %s", rr.Body.String())
	}
}

func TestMultipleParentNodes(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("MultipleParentNodes")
	defer logger.exit("MultipleParentNodes")

	app := setupTestForest(t)
	root := app.Forest

	// Create branch structure for multiple paths to savings
	lifeNode := forestree.NewNode(forestree.BranchNode, "life")
	workNode := forestree.NewNode(forestree.BranchNode, "work")
	paydayNode := forestree.NewNode(forestree.BranchNode, "payday")
	financeNode := forestree.NewNode(forestree.BranchNode, "finance")
	wealthNode := forestree.NewNode(forestree.BranchNode, "wealth")
	houseNode := forestree.NewNode(forestree.BranchNode, "house")
	fundNode := forestree.NewNode(forestree.BranchNode, "fund")
	savingsNode := forestree.NewNode(forestree.LeafNode, "savings")

	// Setup path: life::work::payday::savings
	root.Children["life"] = lifeNode
	lifeNode.Children["work"] = workNode
	workNode.Children["payday"] = paydayNode
	paydayNode.Children["savings"] = savingsNode

	// Setup path: life::finance::savings
	lifeNode.Children["finance"] = financeNode
	financeNode.Children["savings"] = savingsNode

	// Setup path: life::wealth::house::fund::savings
	lifeNode.Children["wealth"] = wealthNode
	wealthNode.Children["house"] = houseNode
	houseNode.Children["fund"] = fundNode
	fundNode.Children["savings"] = savingsNode

	// Verify the node structure
	if len(root.Children["life"].Children) != 3 {
		logger.failure("Expected life node to have 3 children, got %d", len(root.Children["life"].Children))
		t.Errorf("Expected life node to have 3 children, got %d", len(root.Children["life"].Children))
	} else {
		logger.success("Life node has expected number of children")
	}

	// Verify savings node is accessible from all paths
	paths := [][]string{
		{"life", "work", "payday", "savings"},
		{"life", "finance", "savings"},
		{"life", "wealth", "house", "fund", "savings"},
	}

	logger.enter("Path Verification")
	for _, path := range paths {
		currentNode := root
		for _, nodeName := range path {
			if next, exists := currentNode.Children[nodeName]; exists {
				currentNode = next
				logger.success("Path %v is valid at node %s", path, nodeName)
			} else {
				logger.failure("Path %v is broken at node %s", path, nodeName)
				t.Errorf("Path %v is broken at node %s", path, nodeName)
			}
		}
	}
	logger.exit("Path Verification")

	logger.enter("Event Propagation")
	// Test event propagation through all parents
	savingsNode.StartEvent("deposit", nil, nil, map[string]interface{}{"amount": 1000})

	// Verify event is accessible from all paths
	for _, path := range paths {
		currentNode := root
		for _, nodeName := range path {
			currentNode = currentNode.Children[nodeName]
		}
		if _, err := currentNode.GetEventSummary("deposit"); err != nil {
			logger.failure("Event not accessible through path %v: %v", path, err)
			t.Errorf("Event not accessible through path %v: %v", path, err)
		} else {
			logger.success("Event accessible through path %v", path)
		}
	}
	logger.exit("Event Propagation")
}

func TestUserCreationAndAuthentication(t *testing.T) {
	logger := newTestLogger(t)
	logger.enter("UserCreationAndAuthentication")
	defer logger.exit("UserCreationAndAuthentication")

	app := setupTestForest(t)

	logger.enter("User Creation")
	// Test user creation
	createBody := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "securepass123",
	}
	createBytes, _ := json.Marshal(createBody)

	createReq := httptest.NewRequest("POST", "/create_user", bytes.NewBuffer(createBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createRR := httptest.NewRecorder()
	app.handleCreateUser(createRR, createReq)

	if status := createRR.Code; status != http.StatusOK {
		logger.failure("User creation failed with status code: got %v want %v", status, http.StatusOK)
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("User created successfully")
	}
	logger.exit("User Creation")

	logger.enter("Login Tests")
	// Test successful login
	loginBody := map[string]interface{}{
		"username": "testuser",
		"password": "securepass123",
	}
	loginBytes, _ := json.Marshal(loginBody)

	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBytes))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRR := httptest.NewRecorder()
	app.handleLogin(loginRR, loginReq)

	if status := loginRR.Code; status != http.StatusOK {
		logger.failure("Login failed with status code: got %v want %v", status, http.StatusOK)
		t.Errorf("login handler returned wrong status code: got %v want %v", status, http.StatusOK)
	} else {
		logger.success("Login successful")
	}

	var loginResponse map[string]string
	if err := json.NewDecoder(loginRR.Body).Decode(&loginResponse); err != nil {
		logger.failure("Failed to decode login response: %v", err)
		t.Errorf("Failed to decode login response: %v", err)
	} else {
		logger.success("Login response decoded successfully")
	}

	if _, exists := loginResponse["token"]; !exists {
		logger.failure("Login response missing token")
		t.Error("Login response missing token")
	} else {
		logger.success("JWT token received")
	}

	logger.enter("Failed Login Test")
	// Test failed login with wrong password
	wrongLoginBody := map[string]interface{}{
		"username": "testuser",
		"password": "wrongpass",
	}
	wrongLoginBytes, _ := json.Marshal(wrongLoginBody)

	wrongLoginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(wrongLoginBytes))
	wrongLoginReq.Header.Set("Content-Type", "application/json")
	wrongLoginRR := httptest.NewRecorder()
	app.handleLogin(wrongLoginRR, wrongLoginReq)

	if status := wrongLoginRR.Code; status != http.StatusUnauthorized {
		logger.failure("Wrong password login returned unexpected status: got %v want %v",
			status, http.StatusUnauthorized)
		t.Errorf("login handler with wrong password returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	} else {
		logger.success("Wrong password correctly rejected")
	}
	logger.exit("Failed Login Test")
	logger.exit("Login Tests")
}
