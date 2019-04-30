package message

import (
	"fmt"
	"github.com/sfreiberg/gotwilio"
	"os"
)

func MakeVoiceCall() {
	accountSid := os.Getenv("TWILLIO_SID")
	authToken := os.Getenv("TWILLIO_TOKEN")

	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := os.Getenv("TWILLIO_FROM")
	to := os.Getenv("TWILLIO_TO")

	callbackParams := gotwilio.NewCallbackParameters("https://demo.twilio.com/welcome/voice/")
	_, exp, err := twilio.CallWithUrlCallbacks(from, to, callbackParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	if exp != nil {
		fmt.Println(exp.Message)
		return
	}
}
