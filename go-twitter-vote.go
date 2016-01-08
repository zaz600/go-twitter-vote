/*
Пример скрипта для голосования в твиттере без API
*/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
)

func getToken(s string) string {
	tokenRE := regexp.MustCompile(`value="(.+?)" name="authenticity_token"`)
	res := tokenRE.FindStringSubmatch(s)
	return res[1]
}

func main() {

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	resp, err := client.Get("https://twitter.com")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	token := getToken(string(body))
	resp.Body.Close()

	data := url.Values{}
	data.Add("session[username_or_email]", "twiuser")
	data.Add("session[password]", "twipassword")
	data.Add("remember_me", "1")
	data.Add("return_to_ssl", "true")
	data.Add("scribe_log", "")
	data.Add("redirect_after_login", "/")
	data.Add("authenticity_token", token)

	resp2, err := client.PostForm("https://twitter.com/sessions", data)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp2.Body.Close()
	body, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		log.Fatalln(err)
	}

	tweetID := "685225554665013248"
	cardID := "685225554270765056"
	cardName := "poll4choice_text_only"
	voteURL := "https://twitter.com/i/cards/api/v1.json?tweet_id=%s&card_name=%s&authenticity_token=%s&forward=false&capi_uri=capi://passthrough/1"
	voteURL = fmt.Sprintf(voteURL, tweetID, cardName, token)

	s := fmt.Sprintf(`{"twitter:string:card_uri":"card://%s","twitter:long:original_tweet_id":"%s","twitter:string:selected_choice":"1"}`, cardID, tweetID)
	var jsonStr = []byte(s)

	req, err := http.NewRequest("POST", voteURL, bytes.NewBuffer(jsonStr))
	//req.Header.Set("Content-Type", "application/json")

	resp3, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()
	log.Println(resp3.Status)
}
