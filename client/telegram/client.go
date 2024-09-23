package telegram

import (
	"BotMixology/events/telegram/buttons"
	"BotMixology/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}

}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap("cant get updates", err)
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap("cant unmarshal update response", err)
	}

	return res.Result, nil

}

func (c *Client) SendMessage(chatID int, text string, replyMarkUp buttons.ReplyMarkUp) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	replyMarkUpJSON, err := json.Marshal(replyMarkUp)
	if err != nil {
		return e.Wrap("cant marshal reply markup", err)
	}
	q.Add("reply_markup", string(replyMarkUpJSON))

	_, err = c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "Cant do request"
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}
