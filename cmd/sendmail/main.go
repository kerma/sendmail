package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/kerma/sendmail"
)

var (
	toAddr []string
	ccAddr []string

	from     = flag.String("from", "", "From address")
	subject  = flag.String("subject", "", "Email subject")
	to       = flag.String("to", "", "To address(es)")
	cc       = flag.String("cc", "", "CC address(es)")
	confPath = flag.String("conf", "/etc/sendmail/config.json", "Config file path")
)

func fatal(s ...interface{}) {
	log.SetPrefix("ERROR: ")
	log.Println(s...)
	os.Exit(1)
}

func getFrom() *string {
	host, err := os.Hostname()
	if err != nil {
		fatal(err)
	}

	u, err := user.Current()
	if err != nil {
		fatal(err)
	}

	s := u.Username + "@" + host
	return &s
}

func readStdin() string {
	var b bytes.Buffer

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Fprintln(&b, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fatal("Failed reading standard input:", err)
	}
	return b.String()
}

func updateConf(c *sendmail.Config, server *string, port *int, username, password *string) {
	if *server != "" {
		c.Server = *server
	}
	if *port != 0 {
		c.Port = *port
	}
	if *username != "" {
		c.User = *username
		c.Password = *password
	}
}

func main() {

	log.SetFlags(0)
	log.SetPrefix("=> ")

	// Add config override flags
	var c sendmail.Config
	flag.StringVar(&c.Server, "server", "", "SMTP server host")
	flag.IntVar(&c.Port, "port", 0, "SMTP server port")
	flag.StringVar(&c.User, "user", "", "SMTP server username")
	flag.StringVar(&c.Password, "password", "", "SMTP server password")
	flag.Parse()

	file, _ := os.Open(*confPath)
	defer file.Close()

	config, err := sendmail.NewConfig(file)
	if err != nil {
		fatal(err)
	}
	config.Update(&c)

	if config.Server == "" {
		fatal("Invalid config, missing server value")
	}

	if *from == "" {
		from = getFrom()
	}

	if *to == "" {
		fatal("Missing -to address")
	} else {
		toAddr = strings.Split(*to, ",")
	}

	if *cc != "" {
		ccAddr = strings.Split(*cc, ",")
	}

	log.Println("Email body (ctrl-d to send):\n")
	body := readStdin()

	m := sendmail.NewMail(from, toAddr, ccAddr, subject, &body, false)
	log.Println("Sending...")
	sendmail.Send(m, config)
	log.Println("Done.")
}
