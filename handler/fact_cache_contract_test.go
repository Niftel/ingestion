package handler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFactCacheUploadV1Decodes(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("testdata", "fact_cache_upload.v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	var upload FactCacheUpload
	if err := json.Unmarshal(payload, &upload); err != nil {
		t.Fatalf("decode canonical fact upload: %v", err)
	}
	if len(upload.Facts) != 2 {
		t.Fatalf("decoded %d host fact sets, want 2", len(upload.Facts))
	}
	for host, facts := range upload.Facts {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(facts, &object); err != nil {
			t.Fatalf("facts for %s are not an object: %v", host, err)
		}
	}
}

func TestInventoryFactsV1IsHostKeyedObject(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("testdata", "inventory_facts.v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	var facts map[string]json.RawMessage
	if err := json.Unmarshal(payload, &facts); err != nil {
		t.Fatalf("decode canonical inventory facts: %v", err)
	}
	if len(facts) != 2 {
		t.Fatalf("decoded %d host fact sets, want 2", len(facts))
	}
}
