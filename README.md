# How to run
```shell script
go run github.com/smahjoub/events-api
```
or 
```shell script
go build -o ./bin/main
```

## How to test
```shell script
go test -v server.go main.go handlers_test.go  -covermode=count  -coverprofile=./bin/coverage.out
```

## Initial Structure
```
    .
    ├── bin
    │   ├── coverage.out
    │   └── main
    ├── errors
    │   └── errors.go
    ├── handlers
    │   └── handlers.go
    ├── objects
    │   ├── event.go
    │   └── requests.go
    ├── store
    │   ├── postgres.go
    │   └── store.go
    ├── .gitignore
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── handlers_test.go
    ├── LICENSE
    ├── main.go
    ├── README.md
    └── server.go
```

### Rest api
**Object: Event**
```go
package objects

import (
	"time"
)

// EventStatus defines the status of the event
type EventStatus string

const (
	// Some default event status
	Original    EventStatus = "original"
	Cancelled   EventStatus = "cancelled"
	Rescheduled EventStatus = "rescheduled"
)

type TimeSlot struct {
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

// Event object for the API
type Event struct {
	// Identifier
	ID string `gorm:"primary_key" json:"id,omitempty"`

	// General details
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty"`
	Address     string `json:"address,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`

	// Event slot duration
	Slot *TimeSlot `gorm:"embedded" json:"slot,omitempty"`

	// Change status
	Status EventStatus `json:"status,omitempty"`

	// Meta information
	CreatedOn     time.Time `json:"created_on,omitempty"`
	UpdatedOn     time.Time `json:"updated_on,omitempty"`
	CancelledOn   time.Time `json:"cancelled_on,omitempty"`
	RescheduledOn time.Time `json:"rescheduled_on,omitempty"`
}
```

#### Endpoints

**Get event**
```http request
GET http://localhost:8080/api/v1/event?id=20200829011748
Accept: application/json
###
```

**Create an event**
```http request
POST http://localhost:8080/api/v1/event
Content-Type: application/json

{
    "name": "Inaugration",
    "description": "Inaugration of Yes Bank",
    "slot": {
        "start_time": "2020-12-11T09:00:00+05:30",
        "end_time": "2020-12-11T15:00:00+05:30"
    },
    "website": "https://yesbank.com",
    "address": "Yes City"
}
###
```

**List at max 42 events after the event: 20200828011748**
```http request
GET http://localhost:8080/api/v1/events?limit=42&after=20200828011748
Accept: application/json
###
```

**Update event's general details**
```http request
PUT http://localhost:8080/api/v1/event/details
Content-Type: application/json

{
    "id": "20200829011748",
    "name": "Inaugration of Yes Bank",
    "description": "Inaugration of New Yes Bank",
    "website": "https://yesbank.com",
    "address": "Yes City",
    "phone_number": "+2348932"
}
###
```

**Reschedule the event**
```http request
PATCH http://localhost:8080/api/v1/reschedule/cancel?id=20200829011748
Content-Type: application/json

{
    "id": "20200829011748",
    "new_slot": {
        "start_time": "2020-12-12T09:00:00+05:30",
        "end_time": "2020-12-12T15:00:00+05:30"
    }
}
###
```

**Cancel the event**
```http request
PATCH http://localhost:8080/api/v1/event/cancel?id=20200829011748
Content-Type: application/json

###
```

**Delete the event**
```http request
DELETE http://localhost:8080/api/v1/event?id=20200829011748
Content-Type: application/json

###
```

