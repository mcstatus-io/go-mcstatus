package mcstatus

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"strconv"
)

var (
	//go:embed icon.png
	defaultIconBytes []byte
	defaultIcon      image.Image = nil
)

func init() {
	var err error

	if defaultIcon, err = png.Decode(bytes.NewReader(defaultIconBytes)); err != nil {
		panic(err)
	}
}

type IconOptions struct {
	Timeout float64
}

func GetIcon(host string, port uint16, options ...IconOptions) (image.Image, error) {
	opts := IconOptions{
		Timeout: 5.0,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	params := &url.Values{}
	params.Set("timeout", strconv.FormatFloat(opts.Timeout, 'f', 1, 64))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/icon/%s:%s?%s", baseURL, url.PathEscape(host), url.PathEscape(strconv.FormatUint(uint64(port), 10)), params.Encode()), nil)

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

	return png.Decode(resp.Body)
}

func GetDefaultIcon() image.Image {
	return defaultIcon
}
