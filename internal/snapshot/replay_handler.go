package snapshot

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type apiResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

func writeAPIError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiResponse{Status: "error", Error: msg})
}

func (s *ReplayServer) handleQueryRange(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	start, err := parseUnixOrRFC3339(q.Get("start"))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid start: "+err.Error())
		return
	}
	end, err := parseUnixOrRFC3339(q.Get("end"))
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid end: "+err.Error())
		return
	}

	matchers, err := ParseMatchers(q["match[]"])
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid matchers: "+err.Error())
		return
	}

	opts := QueryOptions{
		Matchers: matchers,
		Start:    start,
		End:      end,
	}

	series, err := QuerySeries(s.db, opts)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(apiResponse{Status: "success", Data: series})
}

func parseUnixOrRFC3339(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	if ts, err := strconv.ParseFloat(s, 64); err == nil {
		return time.Unix(int64(ts), 0).UTC(), nil
	}
	return time.Parse(time.RFC3339, s)
}
