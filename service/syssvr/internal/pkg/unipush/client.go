package unipush

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

type Client struct {
	cfg    config.UniPushConf
	client *http.Client
}

type SendReq struct {
	PushClientIDs     []string
	Title             string
	Content           string
	Payload           map[string]any
	ForceNotification *bool
	RequestID         string
}

type sendBody struct {
	PushClientIDs     []string       `json:"push_clientids"`
	Title             string         `json:"title"`
	Content           string         `json:"content"`
	Payload           map[string]any `json:"payload"`
	ForceNotification bool           `json:"force_notification"`
	RequestID         string         `json:"request_id,omitempty"`
}

type sendResp struct {
	ErrCode int            `json:"errCode"`
	ErrMsg  string         `json:"errMsg"`
	Data    map[string]any `json:"data"`
}

func NewClient(cfg config.UniPushConf) *Client {
	timeout := time.Duration(cfg.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.cfg.Enabled && c.cfg.HttpUrl != "" && c.cfg.Secret != ""
}

func (c *Client) Send(ctx context.Context, req SendReq) error {
	if !c.Enabled() {
		return nil
	}
	if len(req.PushClientIDs) == 0 {
		return nil
	}
	body := sendBody{
		PushClientIDs:     req.PushClientIDs,
		Title:             req.Title,
		Content:           req.Content,
		Payload:           req.Payload,
		ForceNotification: resolveForceNotification(c.cfg.ForceNotification, req.ForceNotification),
		RequestID:         req.RequestID,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.HttpUrl, bytes.NewReader(raw))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Push-Secret", c.cfg.Secret)

	res, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var resp sendResp
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return fmt.Errorf("unipush invalid json: %s", string(respBytes))
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("unipush errCode=%d errMsg=%s", resp.ErrCode, resp.ErrMsg)
	}
	if err := checkPushBatchResults(resp.Data); err != nil {
		return err
	}
	logx.WithContext(ctx).Infof("unipush sent clients=%d requestID=%s", len(req.PushClientIDs), req.RequestID)
	return nil
}

func resolveForceNotification(defaultValue bool, reqValue *bool) bool {
	if reqValue != nil {
		return *reqValue
	}
	return defaultValue
}

func checkPushBatchResults(data map[string]any) error {
	if data == nil {
		return nil
	}
	raw, ok := data["results"]
	if !ok {
		return nil
	}
	list, ok := raw.([]any)
	if !ok {
		return nil
	}
	for i, item := range list {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		code := m["errCode"]
		switch v := code.(type) {
		case float64:
			if v != 0 {
				return fmt.Errorf("unipush batch[%d] errCode=%v errMsg=%v", i, code, m["errMsg"])
			}
		case string:
			if v != "" && v != "0" {
				return fmt.Errorf("unipush batch[%d] errCode=%s errMsg=%v", i, v, m["errMsg"])
			}
		}
	}
	return nil
}
