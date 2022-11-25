package marvel

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type MarvelClient struct {
	BaseURL    string
	PubKey     string
	PrivKey    string
	HttpClient *http.Client
}

type character struct {
	//ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type characterResponse struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Data            struct {
		Offset  int         `json:"offset"`
		Limit   int         `json:"limit"`
		Total   int         `json:"total"`
		Count   int         `json:"count"`
		Results []character `json:"results"`
	} `json:"data"`
}

func (c *MarvelClient) md5Hash(ts int64) string {
	tsForHash := strconv.Itoa(int(ts))
	hash := md5.Sum([]byte(tsForHash + c.PrivKey + c.PubKey))
	return hex.EncodeToString(hash[:])
}

func (c *MarvelClient) signURL(url string) string {
	ts := time.Now().Unix()
	hash := c.md5Hash(ts)
	return fmt.Sprintf("%s&ts=%d&apikey=%s&hash=%s", url, ts, c.PubKey, hash)

}

func (c *MarvelClient) GetCharacters(limit int) ([]character, error) {
	url := c.BaseURL + fmt.Sprintf("/characters?limit=%d", limit)
	url = c.signURL(url)
	res, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var characterResponse characterResponse
	if err := json.NewDecoder(res.Body).Decode(&characterResponse); err != nil {
		return nil, err
	}

	return characterResponse.Data.Results, nil
}
