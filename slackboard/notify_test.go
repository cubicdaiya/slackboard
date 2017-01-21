package slackboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNotifyDirectlyHandler(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{
			`{
                "payload": {
                    "channel": "random",
                    "username": "slackboard",
                    "icon_emoji": ":clipboard:",
                    "text": "notification text",
                    "parse": "full"
                },
                "sync": true
            }`,
			`{
                "channel":"random",
                "username":"slackboard",
                "icon_emoji":":clipboard:",
                "text":"notification text",
                "parse":"full",
                "attachments":null
            }`,
		},
		{
			`{
                "payload": {
                    "channel": "general",
                    "username": "bot",
                    "icon_emoji": ":information_desk_person:",
                    "text": "notification general text",
                    "parse": "full"
                },
                "sync": true
            }`,
			`{
                "channel":"general",
                "username":"bot",
                "icon_emoji":":information_desk_person:",
                "text":"notification general text",
                "parse":"full",
                "attachments":null
            }`,
		},
	}

	for _, tt := range testData {
		// setup a mock slack server
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ctype := r.Header.Get("Content-Type"); ctype != "application/json" {
				t.Errorf("content type header: got %v want %v", ctype, "application/json")
			}

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}

			jsonIsEqual, err := jsonBytesEqual(reqBody, []byte(tt.out))
			if err != nil {
				t.Fatal(err)
			}
			if !jsonIsEqual {
				t.Errorf("unexpected body: got %v want %v", string(reqBody), tt.out)
			}

			fmt.Fprint(w, "ok")
		}))
		defer ts.Close()
		ConfSlackboard.Core.SlackURL = ts.URL

		// setup a test client
		req, err := http.NewRequest(
			"POST",
			"/notify-directly",
			bytes.NewBuffer([]byte(tt.in)),
		)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(NotifyDirectlyHandler)
		handler.ServeHTTP(rr, req)

		// check a response
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code: got %v want %v", status, http.StatusOK)
		}

		expected := `{"message":"ok"}`
		jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expected))
		if err != nil {
			t.Fatal(err)
		}
		if !jsonIsEqual {
			t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	}
}

func TestNotifyDirectlyHandlerSlackServerError(t *testing.T) {
	// setup a mock slack server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server down", http.StatusInternalServerError)
	}))
	defer ts.Close()
	ConfSlackboard.Core.SlackURL = ts.URL

	// setup a test client
	inJSONStr := `{
        "payload": {
            "channel": "random",
            "username": "slackboard",
            "icon_emoji": ":clipboard:",
            "text": "notification text",
            "parse": "full"
        },
        "sync": true
    }`
	req, err := http.NewRequest(
		"POST",
		"/notify-directly",
		bytes.NewBuffer([]byte(inJSONStr)),
	)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotifyDirectlyHandler)
	handler.ServeHTTP(rr, req)

	// check a response
	if status := rr.Code; status != http.StatusBadGateway {
		t.Errorf("status code: got %v want %v", status, http.StatusBadGateway)
	}

	expected := `{"message":"failed to post message to slack"}`
	jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expected))
	if err != nil {
		t.Fatal(err)
	}
	if !jsonIsEqual {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestNotifyDirectlyHandlerQPS(t *testing.T) {
	// setup a mock slack server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	}))
	defer ts.Close()
	ConfSlackboard.Core.SlackURL = ts.URL

	inJSONStr := `{
        "payload": {
            "channel": "random",
            "username": "slackboard",
            "icon_emoji": ":clipboard:",
            "text": "notification text",
            "parse": "full"
        },
        "sync": true
    }`

	var testData = []struct {
		in  int
		out map[string]interface{}
	}{
		{
			0,
			map[string]interface{}{
				"code": http.StatusOK,
				"body": `{"message":"ok"}`,
			},
		},
		{
			1,
			map[string]interface{}{
				"code": http.StatusTooManyRequests,
				"body": `{"message":"failed to post message to slack"}`,
			},
		},
		{
			2,
			map[string]interface{}{
				"code": http.StatusOK,
				"body": `{"message":"ok"}`,
			},
		},
	}

	for _, tt := range testData {
		// setup a qpsend
		ConfSlackboard.Core.QPS = tt.in
		QPSEnd = NewQPSPerSlackEndpoint(ConfSlackboard)

		// setup a test client
		ch := make(chan *httptest.ResponseRecorder)
		qcount := 2
		for i := 0; i < qcount; i++ {
			go func() {
				req, err := http.NewRequest(
					"POST",
					"/notify-directly",
					bytes.NewBuffer([]byte(inJSONStr)),
				)
				if err != nil {
					t.Fatal(err)
				}

				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(NotifyDirectlyHandler)
				handler.ServeHTTP(rr, req)
				ch <- rr
			}()
		}

		// check a response
		for i := 0; i < qcount; i++ {
			rr := <-ch
			if i == 1 {
				// always ok
				if status := rr.Code; status != http.StatusOK {
					t.Errorf("status code: got %v want %v", status, http.StatusOK)
				}

				expected := `{"message":"ok"}`
				jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expected))
				if err != nil {
					t.Fatal(err)
				}
				if !jsonIsEqual {
					t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
				}
			} else {
				// depending on a qps setting
				expectedCode := tt.out["code"].(int)
				if status := rr.Code; status != expectedCode {
					t.Errorf("status code: got %v want %v", status, expectedCode)
				}

				expectedBody := tt.out["body"].(string)
				jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expectedBody))
				if err != nil {
					t.Fatal(err)
				}
				if !jsonIsEqual {
					t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expectedBody)
				}
			}
		}
	}
}

func jsonBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
