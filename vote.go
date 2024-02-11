package mcstatus

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type SendVoteOptions struct {
	Version     uint
	Host        string
	Port        uint16
	Timeout     time.Duration
	Username    string
	ServiceName string
	Timestamp   time.Time

	// Votifier 1
	PublicKey string
	IP        string

	// Votifier 2
	Token string
	UUID  string
}

func SendVote(host string, port uint16, options ...SendVoteOptions) error {
	opts := SendVoteOptions{
		Version:     2,
		Timeout:     time.Second * 5,
		ServiceName: "mcstatus.io",
		Timestamp:   time.Now(),
	}

	if len(options) > 0 {
		opts = options[0]
	}

	params := &url.Values{}
	params.Set("version", strconv.FormatUint(uint64(opts.Version), 10))
	params.Set("host", opts.Host)
	params.Set("port", strconv.FormatUint(uint64(opts.Port), 10))
	params.Set("timeout", strconv.FormatFloat(opts.Timeout.Seconds(), 'f', 1, 64))
	params.Set("username", opts.Username)
	params.Set("serviceName", opts.ServiceName)
	params.Set("timestamp", opts.Timestamp.Format(time.RFC3339))

	if opts.Version == 1 {
		params.Set("publickey", opts.PublicKey)
		params.Set("ip", opts.IP)
	} else if opts.Version == 2 {
		params.Set("token", opts.Token)
		params.Set("uuid", opts.UUID)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vote?%s", baseURL, params.Encode()), nil)

	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", fmt.Sprintf("go-mcstatus %s", version))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mcstatus: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
