package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

const tokenApiUrl = "https://api.tailscale.com/api/v2/oauth/token"
const tailscaleApiUrl = "https://api.tailscale.com/api/v2/tailnet/"
const telegramApiUrl = "https://api.telegram.org/bot"

var (
	devicesMap = make(map[string]bool)
	mu         sync.RWMutex
)

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

	mu.Lock()
	for _, value := range devices.Devices {
		devicesMap[value.Name] = value.IsActive
	}
	mu.Unlock()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			devices, err := fetchDevices(client, getDeviceUrl)
			if err != nil {
				fmt.Println("Error fetching devices:", err)
				continue
			}
			for _, value := range devices.Devices {
				mu.RLock()
				oldStatus := devicesMap[value.Name]
				mu.RUnlock()
				newStatus := value.IsActive
				if oldStatus != newStatus {
					var msg string
					if newStatus {
						msg = fmt.Sprintf("🟢 TAILSCALE ALERT\n\nHost: %s\nName: %s\nOS: %s\nStatus: ONLINE\nUser: %s", value.HostName, value.Name, value.OS, value.User)
					} else {
						msg = fmt.Sprintf("🔴 TAILSCALE ALERT\n\nHost: %s\nName: %s\nOS: %s\nStatus: OFFLINE\nUser: %s", value.HostName, value.Name, value.OS, value.User)
					}

					mu.Lock()
					devicesMap[value.Name] = newStatus
					mu.Unlock()

					_, err := sendTelegramMessage(getTelegramUrl, msg)
					if err != nil {
						fmt.Println("Error sending telegram message:", err)
					}
				}
			}
		}
	}()

	http.HandleFunc("/status", EndPointHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
