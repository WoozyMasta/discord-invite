package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// InviteRequest represents the payload sent to Discord API to create an invite.
type InviteRequest struct {
	MaxAge  int  `json:"max_age"`
	MaxUses int  `json:"max_uses"`
	Unique  bool `json:"unique"`
}

// InviteResponse represents the response received from Discord API after creating an invite.
type InviteResponse struct {
	Code string `json:"code"`
}

// makeInvite generates a new Discord invite by sending a request to the Discord API.
// It returns the invite code or an error if the operation fails.
func (c *Config) makeInvite() (string, error) {
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/invites", c.ChannelID)

	inviteReq := InviteRequest{
		MaxAge:  c.MaxAge,
		MaxUses: c.MaxUses,
		Unique:  c.Unique,
	}

	bodyBytes, err := json.Marshal(inviteReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal invite request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bot "+c.BotToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("Error close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Discord API error: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	var inviteResp InviteResponse
	if err := json.NewDecoder(resp.Body).Decode(&inviteResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return inviteResp.Code, nil
}
