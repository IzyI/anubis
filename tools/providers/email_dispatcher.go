package providers

import "fmt"

func SenderEmail(s string) (string, string, error) {
	//TODO.MD: написать отправку sms
	//TODO.MD: понять как можно сделать защиту от большого количества отправки смс
	fmt.Printf("Send email %s \n", s)
	return s + "IdSend", "myEmail_EmailService", nil

}
