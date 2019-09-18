// The idea of this package is:
// pull-push images from any docker repo. Sends an email if something goes wrong
// Usage:
//
// ./deemon your-registry.domain.tld/yourimage:0.1

package main

import (
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	gomail "gopkg.in/gomail.v2"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err)
	}

	_, err = cli.ImagePull(ctx, os.Args[1], types.ImagePullOptions{})
	if err != nil {
		log.Println(err)
		// we should send the error through email
		sendEmail(err.Error())
	} else {
		log.Printf("your-registry.domain.tld is ok")
	}

}

func sendEmail(msg string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "youremail@domain.tld")
	m.SetHeader("To", "youremail@domain.tld")
	m.SetHeader("Subject", "[Alert] your-registry.domain.tld")
	m.SetBody("text/plain", msg)

	d := gomail.Dialer{Host: "your.smtp.dev", Port: 25}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
