// The idea of this package is:
// pull-push images from any docker repo. Sends an email if something goes wrong
// Usage:
//
// ./deemon -image your-repo.domain.tld/yourimage:0.1

package main

import (
	"log"
	"os"

	"flag"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	gomail "gopkg.in/gomail.v2"
)

func main() {
	image := flag.String("image", os.Args[1], "docker image you want to pull")
	flag.Parse()
	// TODO: we should be able to use different smtp hosts
	testPull(*image)
}

func testPull(image string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Println(err)
	}

	_, err = cli.ImagePull(ctx, image, types.ImagePullOptions{})
	log.Println("Pulling image", image)
	if err != nil {
		log.Println(err)
		// we should send the error through email
		sendEmail(err.Error())
	} else {
		log.Printf("%s is ok", image)
	}

}

func sendEmail(msg string) {

	m := gomail.NewMessage()
	m.SetHeader("From", "youremail@domain.tld")
	m.SetHeader("To", "youremail@domain.tld")
	m.SetHeader("Subject", "[Alert] yourregistry.domain.tld")
	m.SetBody("text/plain", msg)

	d := gomail.Dialer{Host: "smtp.domain.tld", Port: 25}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
