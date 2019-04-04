package message

import "fmt"

func SendToTelegram(msg string) {
	// token := os.Getenv("TELEGRAM_BOT_TOKEN")
	// channel := os.Getenv("TELEGRAM_CHANNEL_NAME")
	// url := fmt.Sprintf(
	// 	"https://api.telegram.org/bot%s/sendMessage?chat_id=@%s&text=%s",
	// 	token, channel, msg,
	// )

	// http.Get(url)
	fmt.Println(msg)
}
