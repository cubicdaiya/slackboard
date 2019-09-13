package slackboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// Log is io.Writer
type Log struct {
	bytes.Buffer
}

func TestNotifyDirectlyHandler_PostMessage(t *testing.T) {
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
			if auth := r.Header.Get("Authorization"); auth != "Bearer xxxx-testtoken" {
				t.Errorf("authorization header: got %v want %v", auth, "Bearer xxxx-testtoken")
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

			fmt.Fprint(w, `{"ok":true}`)
		}))
		defer ts.Close()
		SlackPostMessageAPIURL = ts.URL
		ConfSlackboard.Core.SlackToken = "xxxx-testtoken"

		// setup a test client
		req, err := http.NewRequest(
			http.MethodPost,
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

func TestNotifyDirectlyHandler_IncomingWebhook(t *testing.T) {
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
		ConfSlackboard.Core.SlackToken = ""

		// setup a test client
		req, err := http.NewRequest(
			http.MethodPost,
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

func TestNotifyDirectlyHandlerSlackServerError_PostMessage(t *testing.T) {
	// setup a mock slack server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"ok":false, "error":"error-message"}`)
	}))
	defer ts.Close()
	SlackPostMessageAPIURL = ts.URL
	ConfSlackboard.Core.SlackToken = "xxxx-testtoken"

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
		http.MethodPost,
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

func TestNotifyDirectlyHandlerSlackServerError_IncomingWebhook(t *testing.T) {
	// setup a mock slack server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Server down", http.StatusInternalServerError)
	}))
	defer ts.Close()
	ConfSlackboard.Core.SlackURL = ts.URL
	ConfSlackboard.Core.SlackToken = ""

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
		http.MethodPost,
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

	var testData = []struct {
		in  map[string]interface{}
		out map[string]interface{}
	}{
		{
			// expect to disable qps
			map[string]interface{}{
				"qps":      0,
				"sync":     "true",
				"parallel": 2,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "",
			},
		},
		{
			// expect to reject a request
			map[string]interface{}{
				"qps":      1,
				"sync":     "true",
				"parallel": 2,
			},
			map[string]interface{}{
				"code":  http.StatusTooManyRequests,
				"body":  `{"message":"failed to post message to slack"}`,
				"error": "[warning] Reject a sync message due to ratelimit:{\"channel\":\"random\",\"username\":\"slackboard\",\"icon_emoji\":\":clipboard:\",\"text\":\"notification text\",\"parse\":\"full\",\"attachments\":null}\n[error] failed to post message to slack\n",
			},
		},
		{
			// expect to disable qps even core.qps > 0
			map[string]interface{}{
				"qps":      1,
				"sync":     "false",
				"parallel": 2,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "",
			},
		},
		{
			// expect to reject a request
			map[string]interface{}{
				"qps":                1,
				"sync":               "false",
				"max_delay_duration": 0,
				"parallel":           2,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "[warning] Reject an async message due to ratelimit:{\"channel\":\"random\",\"username\":\"slackboard\",\"icon_emoji\":\":clipboard:\",\"text\":\"notification text\",\"parse\":\"full\",\"attachments\":null}\n[error] failed to post message to slack:QPS ratelimit error\n",
			},
		},
		{
			// expect to ensure posting a message
			// because it is acceptable to have 1 in the queue
			map[string]interface{}{
				"qps":                1,
				"sync":               "false",
				"max_delay_duration": 1,
				"parallel":           2,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "",
			},
		},
		{
			// expect to reject 3 requests
			// because it is not acceptable to have more than 1 in the queue
			map[string]interface{}{
				"qps":                1,
				"sync":               "false",
				"max_delay_duration": 1,
				"parallel":           5,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": strings.Repeat("[warning] Reject an async message due to ratelimit:{\"channel\":\"random\",\"username\":\"slackboard\",\"icon_emoji\":\":clipboard:\",\"text\":\"notification text\",\"parse\":\"full\",\"attachments\":null}\n[error] failed to post message to slack:QPS ratelimit error\n", 3),
			},
		},
		{
			// expect to accept 2 requests at the same time
			map[string]interface{}{
				"qps":      2,
				"sync":     "true",
				"parallel": 2,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "",
			},
		},
		{
			// expect to reject 1 request
			// because it is not acceptable to have more than 2 in the queue
			map[string]interface{}{
				"qps":                2,
				"sync":               "false",
				"max_delay_duration": 1,
				"parallel":           5,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": "[warning] Reject an async message due to ratelimit:{\"channel\":\"random\",\"username\":\"slackboard\",\"icon_emoji\":\":clipboard:\",\"text\":\"notification text\",\"parse\":\"full\",\"attachments\":null}\n[error] failed to post message to slack:QPS ratelimit error\n",
			},
		},
		{
			// expect to reject 40 requests
			// because it is not acceptable to have more than 30 in the queue
			map[string]interface{}{
				"qps":                10,
				"sync":               "false",
				"max_delay_duration": 3,
				"parallel":           80,
			},
			map[string]interface{}{
				"code":  http.StatusOK,
				"body":  `{"message":"ok"}`,
				"error": strings.Repeat("[warning] Reject an async message due to ratelimit:{\"channel\":\"random\",\"username\":\"slackboard\",\"icon_emoji\":\":clipboard:\",\"text\":\"notification text\",\"parse\":\"full\",\"attachments\":null}\n[error] failed to post message to slack:QPS ratelimit error\n", 40),
			},
		},
	}

	for _, tt := range testData {
		// setup a qpsend
		ConfSlackboard.Core.QPS = tt.in["qps"].(int)
		if maxWait := tt.in["max_delay_duration"]; maxWait != nil {
			ConfSlackboard.Core.MaxDelayDuration = maxWait.(int)
		} else {
			ConfSlackboard.Core.MaxDelayDuration = -1
		}
		QPSEnd = NewQPSPerSlackEndpoint(ConfSlackboard)

		// setup a logger
		buf := &Log{}
		LogError.Out = buf

		// setup a test client
		ch := make(chan *httptest.ResponseRecorder)
		qcount := tt.in["parallel"].(int)
		for i := 0; i < qcount; i++ {
			go func() {
				req, err := http.NewRequest(
					http.MethodPost,
					"/notify-directly",
					bytes.NewBuffer([]byte(createInJSONStr(tt.in["sync"].(string)))),
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
			if 1 <= i {
				// always ok
				if status := rr.Code; status != http.StatusOK {
					t.Errorf("status code: got %v want %v: in %v", status, http.StatusOK, tt.in)
				}

				expected := `{"message":"ok"}`
				jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expected))
				if err != nil {
					t.Fatal(err)
				}
				if !jsonIsEqual {
					t.Errorf("unexpected body: got %v want %v: in %v", rr.Body.String(), expected, tt.in)
				}
			} else {
				// depending on a qps setting
				expectedCode := tt.out["code"].(int)
				if status := rr.Code; status != expectedCode {
					t.Errorf("status code: got %v want %v: in %v", status, expectedCode, tt.in)
				}

				expectedBody := tt.out["body"].(string)
				jsonIsEqual, err := jsonBytesEqual(rr.Body.Bytes(), []byte(expectedBody))
				if err != nil {
					t.Fatal(err)
				}
				if !jsonIsEqual {
					t.Errorf("unexpected body: got %v want %v: in %v", rr.Body.String(), expectedBody, tt.in)
				}
			}
		}

		if tt.in["sync"].(string) == "false" {
			waitSleep := 1
			if maxWait := tt.in["max_delay_duration"]; maxWait != nil {
				waitSleep += maxWait.(int)
				waitSleep *= 2 // it depends on heuristic
			}
			time.Sleep(time.Duration(waitSleep) * time.Second)
		}
		expectedErrorLog := tt.out["error"].(string)
		if errorLog := buf.String(); !multiLineStringEqual(errorLog, expectedErrorLog) {
			t.Errorf("errorLog: got %v want %v: in %v", errorLog, expectedErrorLog, tt.in)
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

func multiLineStringEqual(a, b string) bool {
	m := strings.Split(a, "\n")
	m2 := strings.Split(b, "\n")
	sort.Strings(m)
	sort.Strings(m2)
	return reflect.DeepEqual(m, m2)
}

func createInJSONStr(sync string) string {
	return `{
        "payload": {
            "channel": "random",
            "username": "slackboard",
            "icon_emoji": ":clipboard:",
            "text": "notification text",
            "parse": "full"
        },
        "sync": ` + sync +
		`}`
}
