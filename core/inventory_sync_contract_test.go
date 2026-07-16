package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeInventorySyncV1(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("testdata", "inventory_sync.v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	hostvars, hosts, groups, err := decodeInventorySync(payload)
	if err != nil {
		t.Fatalf("decode canonical inventory sync: %v", err)
	}
	if len(hostvars) != 2 || len(hosts) != 2 {
		t.Fatalf("decoded hostvars=%d hosts=%d, want 2 each", len(hostvars), len(hosts))
	}
	if len(groups) != 1 || len(groups["web"]) != 1 || groups["web"][0] != "web-1" {
		t.Fatalf("decoded groups = %#v", groups)
	}
}
