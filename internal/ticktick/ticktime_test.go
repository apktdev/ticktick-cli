package ticktick

import (
	"encoding/json"
	"testing"
)

func TestChecklistItemCompletedTimeAcceptsUnixMillis(t *testing.T) {
	var item ChecklistItem
	err := json.Unmarshal([]byte(`{"completedTime":1764053112000,"title":"x"}`), &item)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if item.CompletedTime == "" {
		t.Fatalf("completedTime should be populated")
	}
}

func TestChecklistItemCompletedTimeAcceptsNull(t *testing.T) {
	var item ChecklistItem
	err := json.Unmarshal([]byte(`{"completedTime":null,"title":"x"}`), &item)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if item.CompletedTime != "" {
		t.Fatalf("expected empty completedTime for null, got %q", item.CompletedTime)
	}
}
