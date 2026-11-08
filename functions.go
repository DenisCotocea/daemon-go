package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

func getClient() (*http.Client, error) {
	ctx := context.Background()
	oauthConfig := &clientcredentials.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		TokenURL:     tokenApiUrl,
	}

	return oauthConfig.Client(ctx), nil
}

func fetchDevices(client *http.Client, url string) (*TailscaleDevicesResponse, error) {
	var tailscaleDevices TailscaleDevicesResponse

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tailscaleDevices)
	if err != nil {
		return nil, err
	}

	return &tailscaleDevices, nil
}

func sendTelegramMessage(url string, text string) (*TelegramResponse, error) {
	var telegramResponseStatus TelegramResponse

	telegramPost := &TelegramPost{
		ChatId: os.Getenv("TELEGRAM_CHAT_ID"),
		Text:   text,
	}

	telegramJson, err := json.Marshal(telegramPost)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(telegramJson)
	resp, err := http.Post(url, "application/json; charset=utf-8", body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Telegram API Error (Status %s): %s\n", resp.Status, string(bodyBytes))
	} else {
		err = json.NewDecoder(resp.Body).Decode(&telegramResponseStatus)
		if err != nil {
			return nil, err
		}
	}

	return &telegramResponseStatus, nil
}
