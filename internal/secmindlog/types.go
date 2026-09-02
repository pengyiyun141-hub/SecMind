package secmindlog

import (
	"time"
)

type FetchLog struct {
    Source     string    `json:"source"`
    NewCount   int       `json:"new_count"`
    TotalCount int       `json:"total_count"`
    FetchAt    time.Time `json:"fetch_at"`
}