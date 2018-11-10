package models

import (
	"database/sql"
	"time"

	slugify "github.com/gosimple/slug"
	txdb "github.com/ramene/pubcast/data"
)

// EventsTable holds the internals of the table, i.e,
// the manager of this instance's database pool (Roach).
// Here you could also add things like a `logger` with
// some predefined fields (for structured logging with
// context).
type EventsTable struct {
	txdb *txdb.dbPool
}


// Group is a collection of Organizations
// Refers to the https://www.w3.org/TR/activitystreams-vocabulary/#dfn-group
// Also refers to the Groups table in the database
type Group struct {
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetGroup returns a single Group object or nil
func GetGroup(db *sql.DB, slug string) (*Group, error) {
	row := db.QueryRow(`
		select slug, name, note, created_at, updated_at
		from groups where slug = $1
	`, slug)

	var group Group
	err := row.Scan(&group.Slug,
		&group.Name, &group.Note, &group.CreatedAt, &group.UpdatedAt)

	// This is not an error from the user's perspective
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// PutGroup creates a group with this name and note
func PutGroup(db *sql.DB, name string, note string) (string, error) {
	slug := slugify.MakeLang(name, "en")

	query := `
		INSERT INTO groups (slug, name, note)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(query, slug, name, note)
	return slug, err
}

// Organization is someone who owns a episodes of podcasts
// Refers to the https://www.w3.org/TR/activitystreams-vocabulary/#dfn-organization
// Also refers to the Organizations table in the database
type Organization struct {
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetOrganization gets an organization at any slug
func GetOrganization(db *sql.DB, slug string) (*Organization, error) {
	row := db.QueryRow(`
		select slug, name, note, created_at, updated_at
		from organizations where slug = $1
	`, slug)

	var org Organization
	err := row.Scan(&org.Slug,
		&org.Name, &org.Note, &org.CreatedAt, &org.UpdatedAt)

	// This is not an error from the user's perspective
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &org, nil
}

// PutOrganization creates a group with this name and note
func PutOrganization(db *sql.DB, name string, note string) (string, error) {
	slug := slugify.MakeLang(name, "en")

	query := `
		INSERT INTO organizations (slug, name, note)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(query, slug, name, note)
	return slug, err
}

// Podcast is a something with an audio link, a name, and a note
// Refers to the Podcasts table in the database
type Podcast struct {
	Slug         string    `json:"slug"`
	Name         string    `json:"name"`
	Note         string    `json:"note"`
	ThumbnailURL string    `json:"thumbnail_url"`
	AudioURL     string    `json:"audio_url"`
	MediaType    string    `json:"media_type"`
	PostedAt     time.Time `json:"posted_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
