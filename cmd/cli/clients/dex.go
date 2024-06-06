package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/NorskHelsenett/ror-cli/cmd/cli/config"
	"github.com/NorskHelsenett/ror-cli/cmd/cli/responses"

	"github.com/NorskHelsenett/ror/pkg/rlog"

	"github.com/spf13/viper"
)

var ErrHttpResponse400 = errors.New("api returned http status 400")
var ErrUnhandledResponseCase = errors.New("api returned unhandled error code")

func main() {
}

func FetchDeviceCode() (error, *responses.CodeResponse) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	payload := url.Values{}

	payload.Set("client_id", "ror-cli")
	payload.Set("scope", "openid email profile groups offline_access")

	encodedPayload := payload.Encode()

	req, err := http.NewRequest(
		http.MethodPost,
		viper.GetString(config.ApiDex)+"/device/code",
		strings.NewReader(encodedPayload))
	if err != nil {
		return err, nil
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedPayload)))

	response, err := client.Do(req)
	if err != nil {
		return err, nil
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err, nil
	}

	if response.StatusCode == 200 {
		codeResponse := &responses.CodeResponse{}

		err = json.Unmarshal(bodyBytes, codeResponse)

		if err != nil {
			return err, nil
		}

		return nil, codeResponse
	} else {
		var formatted bytes.Buffer

		err := json.Indent(&formatted, bodyBytes, "", "\t")

		if err != nil {
			return err, nil
		}

		rlog.Fatal("[Contact maintainer] Handle unknown error: ", err, rlog.String("formatted bytes", string(formatted.Bytes())))
	}

	return nil, nil
}

func FetchJWToken(code *responses.CodeResponse) (*responses.TokenResponse, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for {
		// curl http://localhost:5556/dex/token -d device_code=$1 -d grant_type=urn:ietf:params:oauth:grant-type:device_code -d client_id=ror.sky.test.nhn.no
		payload := url.Values{}
		payload.Set("device_code", code.DeviceCode)
		payload.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
		payload.Set("client_id", "ror-cli")

		encodedPayload := payload.Encode()
		req, err := http.NewRequest(
			http.MethodPost,
			viper.GetString(config.ApiDex)+"/token",
			strings.NewReader(encodedPayload))
		if err != nil {
			return nil, err
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(encodedPayload)))

		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if response.StatusCode == 200 {
			if err != nil {
				return nil, err
			}

			tokenResponse := &responses.TokenResponse{}
			err = json.Unmarshal(bodyBytes, tokenResponse)
			if err != nil {
				return nil, err
			}

			return tokenResponse, nil
		} else {
			badRequest := &responses.BadRequestResponse{}
			err = json.Unmarshal(bodyBytes, badRequest)
			if err != nil {
				return nil, err
			}

			if badRequest.Error == "authorization_pending" {
				// Ignore, this is expected behavior
			} else if badRequest.Error == "access_denided" {
				//	shouldBreak = true
				rlog.Fatal("Sorry nothing we can do, access was denied by dex", fmt.Errorf("access denied"))
			} else if badRequest.Error == "expired_token" {
				//	shouldBreak = true
				rlog.Fatal("Device code is over 5 minutes and therefor ruled invalid by dex.", fmt.Errorf("expired token"))
			} else if badRequest.Error == "slow_down" {
				// Ignore, this shouldn't happen but it dosen't do any damage.
			} else {
				//	shouldBreak = true
				rlog.Fatal("[Contact maintainer] Handle unknown error: ", fmt.Errorf("unhandled error"), rlog.String("response", badRequest.Error))
			}
		}

		time.Sleep(time.Duration(code.Interval) * time.Second)
	}
}

func Ping(url string) (int, error) {
	var client = http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	_ = resp.Body.Close()
	return resp.StatusCode, nil
}
