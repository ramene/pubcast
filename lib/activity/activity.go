package activity

import (
	"fmt"
	"time"
	"net/url"
	"strings"

	"github.com/ramene/pubcast/config"
	"github.com/ramene/pubcast/data"
	// "github.com/ramene/pubcast/data/models"
	"github.com/spf13/viper"
)

// RWType determines if the type of the object is read only or read-write
type RWType bool

// EventsTable holds the internals of the table, i.e,
// the manager of this instance's database pool (Roach).
// Here you could also add things like a `logger` with
// some predefined fields (for structured logging with
// context).
type EventsTable struct {
	pool *data.ActivityPool
}

type EventsTableConfig struct {
	ActivityPool *data.ActivityPool
}

// EventRow represents in a `struct` the information we
// can get from the table (some fields are insertable but
// not all - ID and CreatedAt are generated when we `insert`,
// thus, these can only be retrieved).
type EventRow struct {
	Id        int64
	Type      string
	CreatedAt time.Time
}

const (
	// Read indicates the object is only being read.
	Read RWType = false
	// ReadWrite indicates the object is being mutated as well.
	ReadWrite = true
)

// Owns determines wether the app owns an IRI, or Internationalized Resource ID
func Owns(id *url.URL) bool {
	actual := id.Host
	expected := viper.GetString(config.ServerHostname) + ":" + viper.GetString(config.ServerPort)

	return actual == expected
}

// Has determines if the app has ActivityStream data at the IRI (Internationalized Resource ID)
// We expect IRIs to have a path like `/activity/<object>/<value>/<inbox|outbox>`
// func Has(id *url.URL) (bool, error) {
// 	if !matchesURLSpec(id) {
// 		return false, nil
// 	}

// 	group, err := models.GetGroup(data.NewDB(), getSlug(id))
// 	if err != nil {
// 		return false, err
// 	}
// 	if group == nil {
// 		return false, nil
// 	}

// 	return true, nil
// }

func (table *EventsTable) GetEventsByType(eventType string) (rows []EventRow, err error) {
	if eventType == "" {
		err = errors.Errorf("Can't get event rows with empty type")
		return
	}

	return
}

// Deletes empty strings from an array of strings
// ["", "dogs", "oh"] => ["dogs", "oh"]
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func matchesURLSpec(id *url.URL) bool {
	fragments := deleteEmpty(strings.Split(id.Path, "/"))

	// Fragments need four pieces "activity", "<object>", "<value>", "inbox or outbox"
	if len(fragments) != 4 {
		fmt.Println("length", len(fragments))
		return false
	}

	// We only accept paths starting with /activity
	if fragments[0] != "activity" {
		return false
	}

	// The ActivityStreams object we're trying to reference
	object := fragments[1]
	if !(object == "organization" || object == "group") {
		return false
	}

	// The forth piece must be "inbox" or "outbox"
	if !(fragments[3] == "inbox" || fragments[3] == "outbox") {
		return false
	}

	return true
}

func getSlug(id *url.URL) string {
	fragments := deleteEmpty(strings.Split(id.Path, "/"))
	if len(fragments) != 4 {
		return "" // Likely bad.
	}

	return fragments[2]
}
