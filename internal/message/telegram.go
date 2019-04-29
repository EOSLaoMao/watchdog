package message

import (
	"fmt"
	"io/ioutil"
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

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

}
