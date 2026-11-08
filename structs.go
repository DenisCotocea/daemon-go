package main

type TailscaleDevicesResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Addresses []string `json:"addresses"`
	User      string   `json:"user"`
	Name      string   `json:"name"`
	HostName  string   `json:"hostname"`
	OS        string   `json:"os"`
	IsActive  bool     `json:"connectedToControl"`
	LastSeen  string   `json:"lastSeen"`
}

type TelegramPost struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramResponse struct {
	OK bool `json:"ok"`
}
