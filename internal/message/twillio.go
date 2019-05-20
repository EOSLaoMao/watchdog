package message

import (
	"fmt"
	"os"
	"strings"

	"github.com/sfreiberg/gotwilio"
)

func MakeVoiceCall() {
	accountSid := os.Getenv("TWILLIO_SID")
	authToken := os.Getenv("TWILLIO_TOKEN")

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := os.Getenv("TWILLIO_FROM")
	tos := os.Getenv("TWILLIO_TO")

	for _, to := range strings.Split(tos, ",") {
		callbackParams := gotwilio.NewCallbackParameters("https://demo.twilio.com/welcome/voice/")
		_, exp, err := twilio.CallWithUrlCallbacks(from, to, callbackParams)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exp != nil {
			fmt.Println(exp.Message)
			continue
		}
	}
}
