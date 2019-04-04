package message

import (
	"fmt"
	"net/http"
	"os"
)

func SendToTelegram(msg string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	channel := os.Getenv("TELEGRAM_CHANNEL_NAME")
	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=@%s&text=%s&parse_mode=HTML",
		token, channel, msg,
	)

	http.Get(url)
}
