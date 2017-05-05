package services

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"fmt"
)

func SendAccountVerificationEmail(username string, email string) {
	body := "<p><h3>Thank you for joining signing up with BTCPredictions! To verify your account and get started using the simulator, please click the link below.</h3></p><br><br><h4>http://localhost:4200/verifyemail/" + CreateToken(username) + "</h4>"
	from := mail.NewEmail("Kevin", "NoReply@BTCPredictions.com")
	subject := "Successfully reset password!"
	to := mail.NewEmail(username, email)
	content := mail.NewContent("text/plain", body)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest("", "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// func emailPassword(username string, email string, password string) {
// 	body := "Your password has been reset successfully. Your new password is: <br><br><h2>" + password + "</h2><br><br>"
// 	from := mail.NewEmail("Kevin", "NoReply@BTCPredictions.com")
// 	subject := "Successfully reset password!"
// 	to := mail.NewEmail(username, email)
// 	content := mail.NewContent("text/plain", body)
// 	m := mail.NewV3MailInit(from, subject, to, content)

// 	request := sendgrid.GetRequest("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", "/v3/mail/send", "https://api.sendgrid.com")
// 	request.Method = "POST"
// 	request.Body = mail.GetRequestBody(m)
// 	response, err := sendgrid.API(request)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}
// }






