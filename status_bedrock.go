package mcstatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type rawBedrockStatus struct {
	Online      bool   `json:"online"`
	Host        string `json:"host"`
	Port        uint16 `json:"port"`
	IPAddress   string `json:"ip_address"`
	EULABlocked bool   `json:"eula_blocked"`
	RetrievedAt int64  `json:"retrieved_at"`
	ExpiresAt   int64  `json:"expires_at"`
	Version     *struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players *struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players"`
	MOTD struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		HTML  string `json:"html"`
	} `json:"motd"`
	Gamemode string `json:"gamemode"`
	ServerID string `json:"server_id"`
	Edition  string `json:"edition"`
}

type BedrockStatusResponse struct {
	Online      bool      `json:"online"`
	Host        string    `json:"host"`
	Port        uint16    `json:"port"`
	IPAddress   string    `json:"ip_address"`
	EULABlocked bool      `json:"eula_blocked"`
	RetrievedAt time.Time `json:"retrieved_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Version     *struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version,omitempty"`
	Players *struct {
		Online int `json:"online"`
		Max    int `json:"max"`
	} `json:"players,omitempty"`
	MOTD struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		HTML  string `json:"html"`
	} `json:"motd,omitempty"`
	Gamemode string `json:"gamemode,omitempty"`
	ServerID string `json:"server_id,omitempty"`
	Edition  string `json:"edition,omitempty"`
}

type BedrockStatusOptions struct {
	Timeout float64
}

func GetBedrockStatus(host string, port uint16, options ...BedrockStatusOptions) (*BedrockStatusResponse, error) {
	opts := BedrockStatusOptions{
		Timeout: 5.0,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	params := &url.Values{}
	params.Set("timeout", strconv.FormatFloat(opts.Timeout, 'f', 1, 64))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status/bedrock/%s:%s?%s", baseURL, url.PathEscape(host), url.PathEscape(strconv.FormatUint(uint64(port), 10)), params.Encode()), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("go-mcstatus %s", version))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mcstatus: unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var result rawBedrockStatus

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &BedrockStatusResponse{
		Online:      result.Online,
		Host:        result.Host,
		Port:        result.Port,
		IPAddress:   result.IPAddress,
		EULABlocked: result.EULABlocked,
		RetrievedAt: time.UnixMilli(result.RetrievedAt),
		ExpiresAt:   time.UnixMilli(result.ExpiresAt),
		Version:     result.Version,
		Players:     result.Players,
		MOTD:        result.MOTD,
		Gamemode:    result.Gamemode,
		ServerID:    result.ServerID,
		Edition:     result.Edition,
	}, nil
}
