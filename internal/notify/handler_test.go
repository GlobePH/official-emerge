package notify

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jeepers-creepers/emerge/internal/sms"

	"github.com/jackc/pgx"
)

var handler http.Handler

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("$DATABASE_URL must be set")
	}

	cfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		panic(err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:   cfg,
		AfterConnect: sms.PrepareStatements,
	})
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if _, err = pool.Exec("TRUNCATE table subscribers CASCADE;"); err != nil {
		panic(err)
	}

	if _, err := pool.Exec("INSERT INTO subscribers (subscriber_number, access_token) VALUES ('9171234567', '9171234567');"); err != nil {
		panic(err)
	}

	handler = New(pool)

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
		panic(err)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fail()
	}
}
