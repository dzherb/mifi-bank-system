package cbr

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
)

type KeyRate struct {
	Val  float64
	Date time.Time
}

func CurrentKeyRate() (KeyRate, error) {
	body := []byte(`<?xml version="1.0" encoding="utf-8"?>
	<soap12:Envelope
	  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	  xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	  xmlns:soap12="http://www.w3.org/2003/05/soap-envelope"
	>
	  <soap12:Body>
		<AllDataInfoXML xmlns="http://web.cbr.ru/" />
	  </soap12:Body>
	</soap12:Envelope>`,
	)

	data, err := dailyInfo(body)
	if err != nil {
		return KeyRate{}, err
	}

	return parseKeyRate(data)
}

func parseKeyRate(xml []byte) (KeyRate, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xml); err != nil {
		return KeyRate{}, err
	}

	el := doc.FindElement("//KEY_RATE")
	if el == nil {
		return KeyRate{}, errors.New("no KEY_RATE found")
	}

	kr := KeyRate{}

	for _, attr := range el.Attr {
		key := strings.ToLower(attr.Key)
		if key == "val" {
			val, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return KeyRate{}, err
			}

			kr.Val = val

			continue
		}

		if key == "date" {
			date, err := time.Parse("02.01.2006", attr.Value)
			if err != nil {
				return KeyRate{}, err
			}

			kr.Date = date
		}
	}

	return kr, nil
}
