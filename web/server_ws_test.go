package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/injoyai/tdx/protocol"
	"golang.org/x/net/websocket"
)

func TestHandleQuoteWebSocketRequiresCode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ws/quote", nil)
	rec := httptest.NewRecorder()

	handleQuoteWebSocket(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d want %d", rec.Code, http.StatusOK)
	}

	var resp Response
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Code != -1 {
		t.Fatalf("unexpected response code: got %d want -1", resp.Code)
	}
	if !strings.Contains(resp.Message, "股票代码不能为空") {
		t.Fatalf("unexpected message: %s", resp.Message)
	}
}

func TestHandleQuoteWebSocketStreamsQuotePayload(t *testing.T) {
	original := quoteFetcher
	quoteFetcher = func(codes ...string) (protocol.QuotesResp, error) {
		return protocol.QuotesResp{
			{
				Code:    codes[0],
				Active1: 1,
			},
		}, nil
	}
	defer func() {
		quoteFetcher = original
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/quote", handleQuoteWebSocket)
	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/quote?code=000001&interval=1"
	conn, err := websocket.Dial(wsURL, "", "http://localhost/")
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.Close()

	var resp Response
	if err := websocket.JSON.Receive(conn, &resp); err != nil {
		t.Fatalf("receive websocket payload: %v", err)
	}

	if resp.Code != 0 {
		t.Fatalf("unexpected response code: got %d want 0", resp.Code)
	}

	payloadMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected payload type: %T", resp.Data)
	}

	if payloadMap["type"] != "quote" {
		t.Fatalf("unexpected payload type field: %#v", payloadMap["type"])
	}

	codes, ok := payloadMap["codes"].([]interface{})
	if !ok || len(codes) != 1 || codes[0] != "000001" {
		t.Fatalf("unexpected codes payload: %#v", payloadMap["codes"])
	}

	quotes, ok := payloadMap["quotes"].([]interface{})
	if !ok || len(quotes) != 1 {
		t.Fatalf("unexpected quotes payload: %#v", payloadMap["quotes"])
	}
}
