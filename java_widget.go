package mcstatus

import (
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"strconv"
)

type JavaWidgetOptions struct {
	Dark    bool
	Rounded bool
	Timeout float64
}

func GetJavaWidget(host string, port uint16, options ...JavaWidgetOptions) (image.Image, error) {
	opts := JavaWidgetOptions{
		Dark:    true,
		Rounded: true,
		Timeout: 5.0,
	}

	if len(options) > 0 {
		opts = options[0]
	}

	params := &url.Values{}
	params.Set("dark", strconv.FormatBool(opts.Dark))
	params.Set("rounded", strconv.FormatBool(opts.Rounded))
	params.Set("timeout", strconv.FormatFloat(opts.Timeout, 'f', 1, 64))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/widget/java/%s:%s?%s", baseURL, url.PathEscape(host), url.PathEscape(strconv.FormatUint(uint64(port), 10)), params.Encode()), nil)

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
