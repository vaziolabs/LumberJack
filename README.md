# LumberJack API Documentation

## Overview
The LumberJack API provides a hierarchical event tracking system where nodes can have multiple parents and events can be tracked across different organizational paths.

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
curl -X POST http://localhost:8080/plan_event/work/projects/project-alpha \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
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
curl -X POST http://localhost:8080/append_event/work/projects/project-alpha \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user123" \
  -d '{
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
curl -X POST http://localhost:8080/start_time_tracking/work/projects/project-alpha \
  -H "X-User-ID: user123"
```

#### Stop Time Tracking
```bash
curl -X POST http://localhost:8080/stop_time_tracking/work/projects/project-alpha \
  -H "X-User-ID: user123"
```

#### Get Time Tracking Summary
```bash
curl -X GET http://localhost:8080/get_time_tracking/work/projects/project-alpha \
  -H "X-User-ID: user123"
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
curl -X GET "http://localhost:8080/get_event_summary/work/projects/project-alpha?event_id=sprint-1" \
  -H "X-User-ID: user123"
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
