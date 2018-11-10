package data

import (
	"database/sql"
    "fmt"

    "github.com/DATA-DOG/go-txdb"
    "github.com/davecgh/go-spew/spew"
    "github.com/pkg/errors"

    // Using the blank identifier in order to solely
    // provide the side-effects of the package.
    // Eseentially the side effect is calling the `init()`
    // method of `lib/pq`:
    //  func init () {  sql.Register("postgres", &Driver{} }
    // which you can see at `github.com/lib/pq/conn.go`
    _ "github.com/lib/pq"
)

// ActivityPool holds the connection pool to the database
type ActivityPool struct {
	// Db holds a sql.DB pointer that represents a pool of zero or more
	// underlying connections - safe for concurrent use by multiple
	// goroutines -, with freeing/creation of new connections all managed
	// by `sql/database` package.
	Db  *sql.DB
	cfg ActivityPub
}

type ActivityPub struct {
	host            string
	port            int
	username        string
	password        string
}

func GetPool() (s *ActivityPool, err error) {
	if s.Db == nil {
		return
	}

	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)",
			spew.Sdump(s.Db))
	}

	return
}

// NewDB opens a standard DB
func NewDB() (*sql.DB, error) {
	const (
		host   = "postgres"
		port   = 5432
		user   = "postgres"
		password = "postgres"
		dbname = "pubcast"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return sql.Open("postgres", psqlInfo)
}

// SetupTestDB is used to setup a transactional database.
// Use it inside of an `init` function in a test file.
func SetupTestDB() {
    const (
        host   = "postgres"
        port   = 5432
        user   = "postgres"
        password = "postgres"
        dbname = "pubcast_test"
    )

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

        // we register an sql driver named "txdb"
        txdb.Register("txdb", "postgres", psqlInfo)
}

func NewTestDB() (pool ActivityPool, err error) {

	// The first argument corresponds to the driver name that the driver
	// (in this case, `lib/pq`) used to register itself in `database/sql`.
	// The next argument specifies the parameters to be used in the connection.
	// Details about this string can be seen at https://godoc.org/github.com/lib/pq
	
	// return sql.Open("txdb", "twelve")
	
	db, err := sql.Open("txdb", "twelve")
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)",
			spew.Sdump(err))
		return
	}

	// Ping verifies if the connection to the database is alive or if a
	// new connection can be made.
	if err = db.Ping(); err != nil {
			err = errors.Wrapf(err,
				"Couldn't ping postgre database (%s)",
				spew.Sdump(err))
			return
	}

	pool.Db = db
	return
}

// ConnectToTestDB creates a new test db pool and sets it to our data pool ActivityPool
// Call this if you're using ActivityPool somewhere inside a function and want your test
// to use our test db.
func ConnectToTestDB() {

    db, err := NewTestDB()
    if err != nil {
		err = errors.Wrapf(err,
			"Couldn't open connection to postgre database (%s)",
			spew.Sdump(db))
	}

	fmt.Println("Connection successful.")
}

func (r *ActivityPool) Close() (err error) {
	if r.Db == nil {
		return
	}

	if err = r.Db.Close(); err != nil {
		err = errors.Wrapf(err,
			"Errored closing database connection",
			spew.Sdump(r.cfg))
	}

	return
}