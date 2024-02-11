package mcstatus

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type rawJavaStatus struct {
	Online      bool   `json:"online"`
	Host        string `json:"host"`
	Port        uint16 `json:"port"`
	IPAddress   string `json:"ip_address"`
	EULABlocked bool   `json:"eula_blocked"`
	RetrievedAt int64  `json:"retrieved_at"`
	ExpiresAt   int64  `json:"expires_at"`
	Version     *struct {
		NameRaw   string `json:"name_raw"`
		NameClean string `json:"name_clean"`
		NameHTML  string `json:"name_html"`
		Protocol  int    `json:"protocol"`
	} `json:"version"`
	Players *struct {
		Online int `json:"online"`
		Max    int `json:"max"`
		List   []struct {
			UUID      string `json:"uuid"`
			NameRaw   string `json:"name_raw"`
			NameClean string `json:"name_clean"`
			NameHTML  string `json:"name_html"`
		} `json:"list"`
	} `json:"players"`
	MOTD struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		HTML  string `json:"html"`
	} `json:"motd"`
	Icon *string `json:"icon"`
	Mods []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"mods"`
	Software *string `json:"software"`
	Plugins  []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"plugins"`
	SRVRecord *struct {
		Host string `json:"host"`
		Port uint16 `json:"port"`
	} `json:"srv_record"`
}

type JavaStatusResponse struct {
	Online      bool      `json:"online"`
	Host        string    `json:"host"`
	Port        uint16    `json:"port"`
	IPAddress   string    `json:"ip_address"`
	EULABlocked bool      `json:"eula_blocked"`
	RetrievedAt time.Time `json:"retrieved_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	Version     *struct {
		NameRaw   string `json:"name_raw"`
		NameClean string `json:"name_clean"`
		NameHTML  string `json:"name_html"`
		Protocol  int    `json:"protocol"`
	} `json:"version,omitempty"`
	Players *struct {
		Online int `json:"online"`
		Max    int `json:"max"`
		List   []struct {
			UUID      string `json:"uuid"`
			NameRaw   string `json:"name_raw"`
			NameClean string `json:"name_clean"`
			NameHTML  string `json:"name_html"`
		} `json:"list"`
	} `json:"players,omitempty"`
	MOTD struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		HTML  string `json:"html"`
	} `json:"motd,omitempty"`
	Icon image.Image `json:"icon,omitempty"`
	Mods []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"mods,omitempty"`
	Software *string `json:"software"`
	Plugins  []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"plugins,omitempty"`
	SRVRecord *struct {
		Host string `json:"host"`
		Port uint16 `json:"port"`
	} `json:"srv_record,omitempty"`
}

type JavaStatusOptions struct {
	Query   bool
	Timeout float64
}

func GetJavaStatus(host string, port uint16, options ...JavaStatusOptions) (*JavaStatusResponse, error) {
	opts := JavaStatusOptions{
		Query:   true,
		Timeout: 5.0,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	params := &url.Values{}
	params.Set("query", strconv.FormatBool(opts.Query))
	params.Set("timeout", strconv.FormatFloat(opts.Timeout, 'f', 1, 64))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status/java/%s:%s?%s", baseURL, url.PathEscape(host), url.PathEscape(strconv.FormatUint(uint64(port), 10)), params.Encode()), nil)

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

	var result rawJavaStatus

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var icon image.Image = nil

	if result.Icon != nil && strings.HasPrefix(*result.Icon, "data:image/png;base64,") {
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(*result.Icon, "data:image/png;base64,"))

		if err != nil {
			return nil, err
		}

		if icon, err = png.Decode(bytes.NewReader(data)); err != nil {
			return nil, err
		}
	}

	return &JavaStatusResponse{
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
		Icon:        icon,
		Mods:        result.Mods,
		Software:    result.Software,
		Plugins:     result.Plugins,
		SRVRecord:   result.SRVRecord,
	}, nil
}
