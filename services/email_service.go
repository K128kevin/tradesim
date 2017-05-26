package services

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"fmt"
	"net/url"
	"os"
)

func SendAccountVerificationEmail(username string, email string) {
	token := CreateToken(username)
	fmt.Printf("\nToken: %s", token)
	encoded := url.QueryEscape(token)
	fmt.Printf("\nEncoded username: %s", encoded)
	body := "<p><h3>Thank you for joining signing up with DemoInvestor! To verify your account and get started using the simulator, please click the link below.</h3></p><br><br><h4>http://" + os.Getenv("TRADESIM_HOST") + "/verify/" + encoded + "</h4>"
	from := mail.NewEmail("DemoInvestor", "NoReply@DemoInvestor.com")
	subject := "Please Verify DemoInvestor Account"
	to := mail.NewEmail(username, email)
	content := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_TOKEN"), "/v3/mail/send", "https://api.sendgrid.com")
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

func SendNewPasswordEmail(username string, email string, newPassword string) {
	body := "<p><h3>Your DemoInvestor password has been successfully reset. Your new password is:</h3></p><br><br><h4>" + newPassword + "</h4><br><br><p><h3>If you believe this is a mistake, please email me directly at DemoInvestor@gmailcom.</h3></p>"
	from := mail.NewEmail("DemoInvestor", "NoReply@DemoInvestor.com")
	subject := "Successfully Reset DemoInvestor Password!"
	to := mail.NewEmail(username, email)
	content := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_TOKEN"), "/v3/mail/send", "https://api.sendgrid.com")
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

func SendResetPasswordLink(username string, email string, token string) {
	fmt.Printf("\nToken: %s", token)
	encoded := url.QueryEscape(token)
	fmt.Printf("\nEncoded username: %s", encoded)
	link := "http://" + os.Getenv("TRADESIM_HOST") + "/resetPassword/" + encoded
	body := "<p><h3>Please click the following link to reset your password: " + link + "<br><br>If you believe this is a mistake, please email me directly at DemoInvestor@gmailcom.</h3></p>"
	from := mail.NewEmail("DemoInvestor", "NoReply@DemoInvestor.com")
	subject := "Reset Password Request"
	to := mail.NewEmail(username, email)
	content := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, to, content)

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_TOKEN"), "/v3/mail/send", "https://api.sendgrid.com")
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






