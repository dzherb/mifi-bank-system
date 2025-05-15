package cbr

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

const dailyInfoUrl = "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"

func dailyInfo(body []byte) ([]byte, error) {
	r, err := http.Post(
		dailyInfoUrl,
		"application/soap+xml",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close() //nolint:errcheck

	if r.StatusCode != http.StatusOK {
		return nil, errors.New("status " + r.Status)
	}

	return io.ReadAll(r.Body)
}
