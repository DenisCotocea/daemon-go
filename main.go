package main

import (
	"fmt"
	"os"
	"time"
)

const tokenApiUrl = "https://api.tailscale.com/api/v2/oauth/token"
const tailscaleApiUrl = "https://api.tailscale.com/api/v2/tailnet/"
const telegramApiUrl = "https://api.telegram.org/bot"

func main() {
	getDeviceUrl := tailscaleApiUrl + os.Getenv("TAILNET_ID") + "/devices"
	getTelegramUrl := telegramApiUrl + os.Getenv("TELEGRAM_BOT_SECRET") + "/sendMessage"

	client, err := getClient()
	if err != nil {
		fmt.Println("Error getting client:", err)
		return
	}

	devices, err := fetchDevices(client, getDeviceUrl)
	if err != nil {
		fmt.Println("Error fetching devices:", err)
		return
	}

	devicesMap := make(map[string]bool)
	for _, value := range devices.Devices {
		devicesMap[value.HostName] = value.IsActive
	}
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			devices, err := fetchDevices(client, getDeviceUrl)
			if err != nil {
				fmt.Println("Error fetching devices:", err)
				continue
			}
			for _, value := range devices.Devices {
				oldStatus := devicesMap[value.HostName]
				newStatus := value.IsActive
				if oldStatus != newStatus {
					var msg string
					if newStatus {
						msg = fmt.Sprintf("🟢 TAILSCALE ALERT\n\nHost: %s\nStatus: ONLINE\nUser: %s", value.HostName, value.User)
					} else {
						msg = fmt.Sprintf("🔴 TAILSCALE ALERT\n\nHost: %s\nStatus: OFFLINE\nUser: %s", value.HostName, value.User)
					}

					devicesMap[value.HostName] = newStatus

					_, err := sendTelegramMessage(getTelegramUrl, msg)
					if err != nil {
						fmt.Println("Error sending telegram message:", err)
					}
				}
			}
		}
	}()

	select {}
}
