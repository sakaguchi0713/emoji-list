package main

import (
	"encoding/json"
	"github.com/lob-inc/rssp/server/shared/logger"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	logger.Infof("Start big-stamp server.")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("SLASHCOMMAND")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}
	channelID := r.Form.Get("channel_id")

	sendMsgUrl := "https://slack.com/api/chat.postMessage"
	emojiMsg := EmojiList(w, token)

	var msgs []string
	for k, _ := range emojiMsg {
		msgs = append(msgs, "["+k+": :"+k+":]")
	}
	sendMsgUrlOption := url.Values{}
	sendMsgUrlOption.Add("token", token)
	sendMsgUrlOption.Add("channel", channelID)
	sendMsgUrlOption.Add("text", strings.Join(msgs, " "))

	_, err := http.Post(sendMsgUrl+"?"+sendMsgUrlOption.Encode(), "", nil)
	if err != nil {
		http.Error(w, "Error parse json.", http.StatusBadRequest)
	}
}

func EmojiList(w http.ResponseWriter, token string) (map[string]interface{}) {
	emojiListUrlOption := url.Values{}
	emojiListUrlOption.Add("token", token)

	emojiGetUrl := "https://slack.com/api/emoji.list?" + emojiListUrlOption.Encode()
	resp, err := http.Post(emojiGetUrl, "", nil)
	if err != nil {
		http.Error(w, "Error get emoji.list.", http.StatusBadRequest)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error don't response read.", http.StatusBadRequest)
	}

	var result interface{}
	if err = json.Unmarshal(b, &result); err != nil {
		http.Error(w, "Error parse json.", http.StatusBadRequest)
	}
	msg := result.(map[string]interface{})

	var emoji interface{}
	eb, err := json.Marshal(msg["emoji"])
	if err != nil {
		logger.Errorf("Error emoji marshal err: %v", err)
		http.Error(w, "Error get emoji.list.", http.StatusBadRequest)
	}

	if eb == nil {
		logger.Error("not set emoji.")
	}

	json.Unmarshal(eb, &emoji)
	emojiMsg := emoji.(map[string]interface{})

	return emojiMsg
}
