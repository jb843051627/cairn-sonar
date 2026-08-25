package store

import (
	"encoding/json"
	"time"
)

func marshal(value any) (string, error)            { b, err := json.Marshal(value); return string(b), err }
func unmarshal[T any](raw string, target *T) error { return json.Unmarshal([]byte(raw), target) }
func timeText(t time.Time) string                  { return t.UTC().Format(time.RFC3339Nano) }
func parseTime(raw string) time.Time               { t, _ := time.Parse(time.RFC3339Nano, raw); return t }
func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
func intBool(v int) bool { return v != 0 }
