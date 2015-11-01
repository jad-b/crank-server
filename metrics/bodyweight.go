package metrics

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jad-b/torque"
	"github.com/jad-b/torque/users"
	"github.com/jmoiron/sqlx"
)

const (
	// BodyweightSQL is the SQL required to create the Bodyweight table.
	bodyweightTable = `
  user_id integer NOT NULL,
  timestamp timestamp with time zone NOT NULL,
  weight numeric(5,2) NOT NULL CHECK (weight < 1000),
  comment text,
  UNIQUE(user_id, timestamp)
`
)

var (
	bodyweightTableName = fmt.Sprintf("%s.bodyweight", Schema)
)

// CreateTableBodyweight creates the Bodyweight table
func CreateTableBodyweight(db *sqlx.DB) error {
	return torque.CreateTable(
		db,
		Schema,
		bodyweightTableName,
		bodyweightTable,
		true)
}

// Bodyweight is a timestamped bodyweight record, with optional comment.
type Bodyweight struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Weight    float64   `json:"weight"`
	Comment   string    `json:"comment"`
}

/*
	DBResource
*/

// Create inserts a new bodyweight entry into the DB.
func (bw *Bodyweight) Create(db *sqlx.DB) error {
	// Need to ensure we only record second-precision
	bw.Timestamp = bw.Timestamp.Truncate(time.Second)
	_, err := db.NamedExec(fmt.Sprintf(`
	INSERT INTO %s.%s (
		user_id,
		timestamp,
		weight,
		comment
	) VALUES (
		:user_id,
		:timestamp,
		:weight,
		:comment
	)`, Schema, bodyweightTableName), bw)
	return err
}

// Retrieve does a lookup for the corresponding bodyweight record by timestamp.
func (bw *Bodyweight) Retrieve(db *sqlx.DB) error {
	return db.Get(
		bw,
		fmt.Sprintf(`
		SELECT
			user_id,
			timestamp,
			weight,
			comment
		FROM %s.%s
		WHERE user_id=$1 AND timestamp=$2`,
			Schema, bodyweightTableName),
		bw.UserID, bw.Timestamp)
}

// Update modifies the matching row in the DB by timestamp.
func (bw *Bodyweight) Update(db *sqlx.DB) error {
	_, err := db.NamedExec(
		fmt.Sprintf(`
			UPDATE %s.%s
			SET
				user_id=:user_id,
				weight=:weight,
				comment=:comment
			WHERE timestamp=:timestamp`,
			Schema, bodyweightTableName),
		bw)
	return err
}

// Delete removes the row from the DB
func (bw *Bodyweight) Delete(db *sqlx.DB) error {
	stmt := fmt.Sprintf(`
			DELETE FROM %s.%s
			WHERE timestamp=:timestamp`,
		Schema, bodyweightTableName)
	res, err := db.NamedExec(stmt, bw)
	// DELETE had no affect
	if i, _ := res.RowsAffected(); i == 0 {
		return errors.New("Resource does not exist")
	}
	return err
}

/*
	RESTfulHandler
*/

// HandlePost creates a new bodyweight record.
func (bw Bodyweight) HandlePost(w http.ResponseWriter, req *http.Request) {
	log.Print("Request: Create Bodyweight")
	err := torque.ReadJSONRequest(req, bw)
	if err != nil {
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	if err = bw.Create(torque.DB); err != nil {
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	log.Printf("Created %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandleGet returns the related bodyweight record
// Lookup performed by timestamp and user id
func (bw Bodyweight) HandleGet(w http.ResponseWriter, req *http.Request) {
	torque.LogRequest(req)
	var err error
	// Get timestamp from query params
	bw.Timestamp, err = torque.GetTimestampQuery(req)
	if err != nil {
		log.Print(err)
		http.Error(w, "Missing timestamp in query parameters", http.StatusBadRequest)
		return
	}
	// Get user ID from query params OR auth token
	var uID int
	uIDparam := req.URL.Query().Get("user_id")
	if uIDparam == "" { // Not found; default to auth token lookup
		// Auth header guaranteed by this point
		authHeader := req.Header.Get(torque.HeaderAuthorization)
		authToken := users.ParseAuthToken(authHeader)
		uID, err = users.SwapTokenForID(torque.DB, authToken)
		if err != nil {
			http.Error(w, "No user found for that auth token", http.StatusBadRequest)
			return
		}
	} else { // Found, but a string
		uID, err = strconv.Atoi(uIDparam)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s is an invalid user ID", uID), http.StatusBadRequest)
			return
		}
	}
	bw.UserID = uID
	// DB retrieval
	log.Printf("Retrieving %+v", bw)
	if err := (&bw).Retrieve(torque.DB); err != nil {
		log.Print(err)
		torque.BadRequest(w, req, "No record found")
		return
	}
	log.Printf("Retrieved %+v", bw)
	torque.WriteOkayJSON(w, bw)
}

// HandlePut updates a Bodyweight resource.
func (bw Bodyweight) HandlePut(w http.ResponseWriter, req *http.Request) {
	// Parse body of PUT request into a Bodyweight struct
	err := torque.ReadJSONRequest(req, &bw)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	// Update in DB
	if err = (&bw).Update(torque.DB); err != nil {
		log.Print(err)
		http.Error(w, "Failed to write record to database", http.StatusInternalServerError)
		return
	}
	torque.WriteOkayJSON(w, bw)
}

// HandleDelete removes the bodyweight record from the database.
func (bw Bodyweight) HandleDelete(w http.ResponseWriter, req *http.Request) {
	// Parse DELETE body into Bodyweight struct
	err := torque.ReadJSONRequest(req, &bw)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to parse JSON from request", http.StatusBadRequest)
		return
	}
	// Delete from DB
	if err = (&bw).Delete(torque.DB); err != nil {
		log.Print(err)
		http.NotFound(w, req)
		return
	}
	log.Printf("Deleted bodyweight @ %s", bw.Timestamp)
	torque.WriteJSON(w, http.StatusNoContent, nil)
}

/*
	RESTfulResource
*/

// GetResourceName returns the name the resource wishes to be refered to by in
// the URL
func (bw Bodyweight) GetResourceName() string {
	return torque.SlashJoin(Category, "bodyweight/")
}

// RegisterURL sets up the handler for the Bodyweight reosurce on the server.
func (bw *Bodyweight) RegisterURL() error { return nil }

/*
	HTTP{Poster, Getter, Updater, Deleter}
*/

// Get prepares a URL for  GET'ing a bodyweight record from the server
// This includes:
// - Timestamp truncation to seconds
// - Setting timestamp query field
func (bw *Bodyweight) Get(earl url.URL) url.URL {
	q := earl.Query()
	q.Set("timestamp", torque.Stamp(bw.Timestamp.Truncate(time.Second)))
	q.Set("user_id", strconv.Itoa(bw.UserID))
	earl.RawQuery = q.Encode()
	return earl
}
