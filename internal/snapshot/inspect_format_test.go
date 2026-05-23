package snapshot

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestInspectReport_JSONRoundtrip(t *testing.T) {
	orig := InspectReport{
		Path:  "/tmp/snap",
		Total: 1,
		Blocks: []BlockInspection{
			{
				ULID:       "01HXZ0000000000000000000",
				MinTime:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				MaxTime:    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
				Duration:   "24h0m0s",
				Series:     10,
				Samples:    500,
				Chunks:     5,
				Compaction: 1,
			},
		},
	}

	b, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got InspectReport
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Path != orig.Path {
		t.Errorf("path mismatch: %s != %s", got.Path, orig.Path)
	}
	if got.Total != orig.Total {
		t.Errorf("total mismatch: %d != %d", got.Total, orig.Total)
	}
	if len(got.Blocks) != 1 || got.Blocks[0].ULID != orig.Blocks[0].ULID {
		t.Error("block ULID mismatch after roundtrip")
	}
}

func TestInspectReport_JSONKeys(t *testing.T) {
	r := InspectReport{Path: "/p", Total: 0, Blocks: []BlockInspection{}}
	b, _ := json.Marshal(r)
	s := string(b)
	for _, key := range []string{"path", "total_blocks", "blocks"} {
		if !strings.Contains(s, key) {
			t.Errorf("missing JSON key: %s", key)
		}
	}
}

func TestBlockInspection_LabelsOmittedWhenNil(t *testing.T) {
	insp := BlockInspection{
		ULID:    "X",
		Series:  1,
		Labels:  nil,
	}
	b, _ := json.Marshal(insp)
	if strings.Contains(string(b), "labels") {
		t.Error("labels should be omitted when nil")
	}
}
