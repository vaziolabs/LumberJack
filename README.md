# LumberJack API Documentation
![lumberjack500](https://github.com/user-attachments/assets/1f1c2f9e-a550-47eb-8f18-a11abbc1cf58)
## Overview
The LumberJack API provides a hierarchical event tracking system where nodes can have multiple parents and events can be tracked across different organizational paths.

## Installation & Usage

### Using as a Package
```bash
go get github.com/vaziolabs/lumberjack
```

### Building from Source
```bash
git clone https://github.com/vaziolabs/lumberjack.git
cd lumberjack
go build
```

### Getting Started
Starting LumberJack is as simple as running the following command:
```bash
./lumberjack
```

To create a new server configuration:
```bash
./lumberjack create
```

To start the server:
```bash
./lumberjack start

# Or with dashboard
./lumberjack start -d
```

To list current configuration:
```bash
./lumberjack list
```

To delete configuration:
```bash
./lumberjack delete
```

## TODOS:
 - [ ] Improved Testing
    - [ ] Fix Testing Logging and Scoping to create Run directives
    - [ ] Remove redundant tests
    - [ ] Finish incomplete tests
 - [ ] Refactor CLI types to work with the Core for directory structure
    - [ ] Ensure logger is logging to the correct file
    - [ ] Ensure the dat file is created and updated in the correct directory
 - [ ] Add Core Logger
    - [ ] Integrate debug logging into log file
 - [ ] Improve CLI
   - [X] Allow Multiple Databases and configs
   - [ ] Add proper linux directory structure
   - [ ] Test all Help commands
   - [ ] Ensure delete commands require admin confirmation
   - [ ] Update `list` command to use the proper ID
   - [ ] Fix `delete` command to not delete all databases
   - [ ] Test ALL commands
   - [ ] have CLI daemonize API and Dashboard
   - [X] Remove Admin from Config
   - [ ] Add Windows Support
 - [ ] Test Dashboard Data Display and Interaction
    - [ ] Add Dashboard Login
    - [ ] Add API Event Logging
 - [ ] Create Typescript module for direct API integration
 - [ ] Add TLS
 - [ ] Improve Session Based Authentication
    - [X] Add JWT
    - [ ] Allow for Certificate Authentication

## Core Concepts

### Nodes
- **Branch Node**: Can contain other nodes
- **Leaf Node**: End points for tracking events
- Each node can have multiple parents, enabling flexible organizational structures

### Events
An Event represents a tracked activity with start/end times and associated entries.

#### Attributes:
- `StartTime`: When the event begins
- `EndTime`: When the event concludes
- `Entries`: List of timestamped records
- `Metadata`: Custom event data
- `Status`: pending/ongoing/finished

### Entries
Timestamped records within an event.

#### Attributes:
- `Timestamp`: Creation time
- `Content`: Entry data
- `Metadata`: Additional entry info
- `UserID`: Creator identifier

## API Endpoints

### Node Management

#### Create Node
```bash
curl -X POST http://localhost:8080/create_node \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "path": "work/projects",
    "name": "project-alpha",
    "type": "leaf"
  }'
```

#### Assign User
```bash
curl -X POST http://localhost:8080/assign_user \
  -H "Content-Type: application/json" \
  -H "X-User-ID: admin" \
  -d '{
    "path": "work/projects/project-alpha",
    "user": {
      "id": "user123",
      "username": "john_doe",
      "permissions": ["read", "write"]
    },
    "permission": "write"
  }'
```

### Event Management

#### Start Event
```bash
curl -X POST http://localhost:8080/start_event \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1",
    "metadata": {
      "type": "sprint",
      "duration": "2 weeks",
      "team": "alpha"
    }
  }'
```

#### Plan Event
```bash
curl -X POST http://localhost:8080/plan_event/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-2",
    "start_time": "2024-01-15T09:00:00Z",
    "end_time": "2024-01-29T17:00:00Z",
    "metadata": {
      "type": "sprint",
      "frequency": "bi-weekly",
      "custom_pattern": "0900",
      "category": "work::projects::sprints"
    }
  }'
```

#### Append to Event
```bash
curl -X POST http://localhost:8080/append_event/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1",
    "content": "Completed user authentication feature",
    "metadata": {
      "type": "milestone",
      "story_points": 5
    }
  }'
```

#### End Event
```bash
curl -X POST http://localhost:8080/end_event \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

### Time Tracking

#### Start Time Tracking
```bash
curl -X POST http://localhost:8080/start_time_tracking/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

#### Stop Time Tracking
```bash
curl -X POST http://localhost:8080/stop_time_tracking/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

#### Get Time Tracking Summary
```bash
curl -X GET http://localhost:8080/get_time_tracking/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

### Queries

#### Get Event Entries
```bash
curl -X POST http://localhost:8080/get_event_entries \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

#### Get Event Summary
```bash
curl -X GET http://localhost:8080/get_event_summary/ \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
    "path": "work/projects/project-alpha",
    "event_id": "sprint-1"
  }'
```

## Response Formats

### Event Summary Response
```json
{
  "event_id": "sprint-1",
  "status": "finished",
  "start_time": "2024-01-01T09:00:00Z",
  "end_time": "2024-01-14T17:00:00Z",
  "entries_count": 15,
  "metadata": {
    "type": "sprint",
    "team": "alpha"
  }
}
```

### Time Tracking Summary Response
```json
[
  {
    "start_time": "2024-01-04T09:00:00Z",
    "end_time": "2024-01-04T17:00:00Z",
    "duration": 28800000000000
  }
]
```

## Error Handling
All endpoints return standard HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

Error responses include a message:
```json
{
  "error": "Invalid event ID format"
}
```
