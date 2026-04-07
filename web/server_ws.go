package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/injoyai/tdx/protocol"
	"golang.org/x/net/websocket"
)

var quoteFetcher = func(codes ...string) (protocol.QuotesResp, error) {
	return client.GetQuote(codes...)
}

func handleQuoteWebSocket(w http.ResponseWriter, r *http.Request) {
	codes := splitCodes(r.URL.Query().Get("code"))
	if len(codes) == 0 {
		errorResponse(w, "股票代码不能为空")
		return
	}

	interval, err := parseWSInterval(r.URL.Query().Get("interval"))
	if err != nil {
		errorResponse(w, err.Error())
		return
	}

	websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()

		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		if err := sendQuoteFrame(conn, codes, interval); err != nil {
			return
		}

		for range ticker.C {
			if err := sendQuoteFrame(conn, codes, interval); err != nil {
				return
			}
		}
	}).ServeHTTP(w, r)
}

func parseWSInterval(raw string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return 3, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("interval 参数无效，必须为 1-60 的整数秒")
	}
	if value > 60 {
		return 60, nil
	}
	return value, nil
}

func sendQuoteFrame(conn *websocket.Conn, codes []string, interval int) error {
	quotes, err := quoteFetcher(codes...)
	if err != nil {
		return websocket.JSON.Send(conn, Response{
			Code:    -1,
			Message: fmt.Sprintf("获取行情失败: %v", err),
			Data: map[string]interface{}{
				"type":      "quote",
				"codes":     codes,
				"interval":  interval,
				"timestamp": time.Now().Format(time.RFC3339),
			},
		})
	}

	return websocket.JSON.Send(conn, Response{
		Code:    0,
		Message: "success",
		Data: map[string]interface{}{
			"type":      "quote",
			"codes":     codes,
			"interval":  interval,
			"timestamp": time.Now().Format(time.RFC3339),
			"quotes":    quotes,
		},
	})
}
