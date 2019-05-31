package gobblredis

import (
	"reflect"
	"testing"

	"github.com/calebhiebert/gobbl/session"
	"github.com/go-redis/redis"
)

func TestGet(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	_, err := sessionStore.Get("dummy-id")
	if err != sess.ErrSessionNonexistant {
		t.Error("Get is not returning an error for empty sessions")
	}
}

func TestCreate(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sessionStore.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	session, err := sessionStore.Get("test-id")
	if err != nil {
		t.Error("Received error on session retrieval")
	}

	s := session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestUpdate(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	testData := map[string]interface{}{
		"to_be_deleted":     "pickles",
		"to_be_overwritten": "not pickles",
	}

	err := sessionStore.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	updatedTestData := map[string]interface{}{
		"to_be_overwritten": "definitely pickles",
	}

	err = sessionStore.Update("test-id", &updatedTestData)
	if err != nil {
		t.Error("Received error on session updating")
	}

	session, err := sessionStore.Get("test-id")
	if err != nil {
		t.Error("Received error on session retrieval")
	}

	s := session

	if s["to_be_overwritten"] != "definitely pickles" {
		t.Errorf("Improper update, expected: %s, got: %s", "definitely pickles", s["to_be_overwritten"])
	}

	if _, exists := s["persistent_data"]; exists == true {
		t.Error("Update is not removing old values")
	}
}

func TestUpdateCreate(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sessionStore.Update("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	session, err := sessionStore.Get("test-id")
	if err != nil {
		if err != sess.ErrSessionNonexistant {
			t.Error("Received error on session retrieval")
		}

		return
	}

	s := session

	if s["test-data"] != "Wow" {
		t.Errorf("Session Data incorrect, got: %s, want: %s", s["test-data"], "Wow")
	}
}

func TestDestroy(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	testData := map[string]interface{}{
		"test-data": "Wow",
	}

	err := sessionStore.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation")
	}

	err = sessionStore.Destroy("test-id")
	if err != nil {
		t.Error("Received error on session destruction")
	}

	_, err = sessionStore.Get("test-id")
	if err != sess.ErrSessionNonexistant {
		t.Errorf("Session get should have returned ErrSessionNonexistant, instead got %+v", err)
	}
}

func TestDataTypes(t *testing.T) {
	var sessionStore sess.SessionStore = New(&redis.Options{
		Addr: "localhost:6379",
	}, 0, "session:")

	testData := map[string]interface{}{
		"test-int":   1,
		"test-float": 0.01,
		"test-int64": 1000000000000000000,
	}

	err := sessionStore.Create("test-id", &testData)
	if err != nil {
		t.Error("Received error on session creation", err)
	}

	session, err := sessionStore.Get("test-id")
	if err != nil {
		t.Error("Received error on session retrieval", err)
	}

	s := session

	if reflect.TypeOf(s["test-int"]).String() != "int64" {
		t.Errorf("Incorrect number type, expected int64, got %v", reflect.TypeOf(s["test-int"]))
	}

	if reflect.TypeOf(s["test-float"]).String() != "float64" {
		t.Errorf("Incorrect float type, expected float64, got %v", reflect.TypeOf(s["test-float"]))
	}
}
