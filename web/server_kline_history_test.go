package main

import (
	"testing"
	"time"

	"github.com/injoyai/tdx/protocol"
)

func TestFilterKlineHistoryByDateRange(t *testing.T) {
	resp := &protocol.KlineResp{
		Count: 3,
		List: []*protocol.Kline{
			{Time: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), Close: 1000},
			{Time: time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), Close: 2000},
			{Time: time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC), Close: 3000},
		},
	}

	start := time.Date(2024, 11, 2, 0, 0, 0, 0, time.Local)
	end := time.Date(2024, 11, 3, 0, 0, 0, 0, time.Local)

	filtered := filterKlineHistoryByDateRange(resp, start, end)

	if filtered.Count != 2 {
		t.Fatalf("unexpected count: got %d want 2", filtered.Count)
	}
	if len(filtered.List) != 2 {
		t.Fatalf("unexpected list length: got %d want 2", len(filtered.List))
	}
	if got := filtered.List[0].Time.Format("20060102"); got != "20241102" {
		t.Fatalf("unexpected first item date: %s", got)
	}
	if got := filtered.List[1].Time.Format("20060102"); got != "20241103" {
		t.Fatalf("unexpected second item date: %s", got)
	}
}
