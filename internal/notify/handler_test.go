package notify

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

const test_json = `{
	"inboundSMSMessageList": {
		"inboundSMSMessage": [{
			"dateTime":"Sat Jul 23 2016 14:06:48 GMT+0000 (UTC)",
			"destinationAddress":"tel:29290586859",
			"messageId":"579379f8db2c71010040c546",
			"message":"Testing receive",
			"resourceURL":null,
			"senderAddress":"tel:+639171234567"
		}],
		"numberOfMessagesInThisBatch": 1,
		"resourceURL": null,
		"totalNumberOfPendingMessages": 0
	}
}`

func TestNotify(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(test_json))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec("TRUNCATE TABLE subscribers CASCADE;"); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO subscribers (subscriber_number, access_token) VALUES ('9171234567', '9171234567');"); err != nil {
		log.Fatal(err)
	}
	h := Handler(db)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Notify fail.\nCode: %d\tBody: %s\n", w.Code, w.Body.String())
	}
}
