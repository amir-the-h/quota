package quota

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

// GOB64 is a serializer that uses gob and base64.

// serialize serializes the given data.
func serialize(i interface{}) (string, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(i)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// SerializeCandle serializes the given candle.
func SerializeCandle(c *Candle) (data string, err error) {
	data, err = serialize(c)
	if err != nil {
		return "", fmt.Errorf("failed to encode the candle: %s", err)
	}

	return data, err
}

// DeserializeCandle deserializes the given data.
func DeserializeCandle(data string) (c *Candle, err error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the candle: %s", err)
	}

	bReader := bytes.NewReader(b)
	d := gob.NewDecoder(bReader)
	err = d.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the candle: %s", err)
	}

	return c, err
}

// SerializeQuota serializes the given quota.
func SerilizeQuota(q *Quota) (data string, err error) {
	data, err = serialize(q)
	if err != nil {
		return "", fmt.Errorf("failed to encode the quota: %s", err)
	}

	return data, err
}

// DeserializeQuota deserializes the given data.
func DeserializeQuota(data string) (q *Quota, err error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the quota: %s", err)
	}

	bReader := bytes.NewReader(b)
	d := gob.NewDecoder(bReader)
	err = d.Decode(q)
	if err != nil {
		return nil, fmt.Errorf("failed to decode the quota: %s", err)
	}

	return q, err
}
