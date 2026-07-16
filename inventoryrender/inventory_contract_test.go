package inventoryrender

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/praetordev/models"
)

func TestBuildMatchesRenderedInventoryV1(t *testing.T) {
	want, err := os.ReadFile(filepath.Join("testdata", "inventory_rendered.v1.ini"))
	if err != nil {
		t.Fatal(err)
	}
	hosts := []models.Host{
		{ID: 1, Name: "web-1", Variables: json.RawMessage(`{"ansible_port":22,"ansible_host":"10.0.0.10"}`)},
		{ID: 2, Name: "db-1", Variables: json.RawMessage(`{"ansible_host":"10.0.0.20"}`)},
	}
	groups := []models.Group{{ID: 10, Name: "web"}}
	got := build(hosts, groups, map[int64][]int64{10: {1}})
	if got != string(want) {
		t.Fatalf("rendered inventory contract drifted:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
