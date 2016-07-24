package subscription

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var db *sql.DB

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestSubscribe(t *testing.T) {
	if _, err := db.Exec("TRUNCATE TABLE subscribers CASCADE;"); err != nil {
		log.Fatal(err)
	}
	url := `?access_token=1ixLbltjWkzwqLMXT-8UF-UQeKRma0hOOWFA6o91oXw&subscriber_number=9171234567`
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	h := Handler(db)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	expected := "9171234567 is subscribed.\n"
	actual := w.Body.String()
	if w.Code != http.StatusOK || expected != actual {
		t.Errorf("Subscription fail.\nCode: %d\t Body: %s\n", w.Code, actual)
	}
}

func TestNotModified(t *testing.T) {
	if _, err := db.Exec("TRUNCATE TABLE subscribers CASCADE;"); err != nil {
		log.Fatal(err)
	}
	s := subscriber{
		AccesToken:       "1ixLbltjWkzwqLMXT-8UF-UQeKRma0hOOWFA6o91oXw",
		SubscriberNumber: "9171234567",
	}
	if err := subscribe(s, db); err != nil {
		log.Fatal(err)
	}
	url := "?access_token=" + s.AccesToken + "&subscriber_number=" + s.SubscriberNumber
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	h := Handler(db)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusNotModified {
		t.Errorf("Did not respond with HTTP 304.\nCode: %d\t Body: %s\n", w.Code, w.Body.String())
	}
}
