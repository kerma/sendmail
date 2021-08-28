package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"os/user"
	"path"

	"github.com/kerma/sendmail"
)

const (
	defaultConfig string = "sendmail/config.json"
)

var (
	from     = flag.String("from", "", "From address")
	subject  = flag.String("subject", "", "Email subject")
	to       = flag.String("to", "", "To address(es)")
	cc       = flag.String("cc", "", "CC address(es)")
	dryrun   = flag.Bool("dryrun", false, "Testing mode")
	confPath = flag.String("conf", "", "Config file path")
	attach   = flag.String("attach", "", "Attachment path")
	html     = flag.Bool("html", false, "Send as html")
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

func readLines(s *bufio.Scanner) string {
	var b bytes.Buffer

	for s.Scan() {
		fmt.Fprintln(&b, s.Text())
	}
	if err := s.Err(); err != nil {
		fatal("Failed reading standard input:", err)
	}
	return b.String()
}

func readLine(o string, s *bufio.Scanner) string {
	var b bytes.Buffer

	fmt.Printf(o)
	s.Scan()
	fmt.Fprint(&b, s.Text())
	if err := s.Err(); err != nil {
		fatal("Failed reading standard input:", err)
	}
	return b.String()

}

func main() {

	log.SetFlags(0)
	log.SetPrefix("=> ")

	var c sendmail.Config
	flag.StringVar(&c.Server, "server", "", "SMTP server host")
	flag.IntVar(&c.Port, "port", 0, "SMTP server port")
	flag.StringVar(&c.User, "user", "", "SMTP server username")
	flag.StringVar(&c.Password, "password", "", "SMTP server password")
	flag.Parse()

	file, _ := os.Open(getConfPath(confPath))
	defer file.Close()

	config, _ := sendmail.NewConfig(file)
	config.Update(&c)

	if config.Server == "" {
		fatal("Invalid config, missing server value")
	}

	scanner := bufio.NewScanner(os.Stdin)

	if *from == "" {
		from = getFrom()
	}

	if *to == "" {
		*to = readLine("To: ", scanner)
	}
	toAddr, err := mail.ParseAddressList(*to)
	if err != nil {
		log.Println(err)
		fatal("Invalid email addresses:", *to)
	}

	var ccAddr []*mail.Address
	if *cc != "" {
		ccAddr, err = mail.ParseAddressList(*cc)
		if err != nil {
			fatal("Invalid email addresses:", *cc)
		}
	}

	if *subject == "" {
		*subject = readLine("Subject: ", scanner)
	}

	log.Printf("Email body (ctrl-d to send):\n\n")
	body := readLines(scanner)

	m := sendmail.NewMail(from, toAddr, ccAddr, *subject, body, *attach, *html)
	log.Println("Sending...")
	if *dryrun == false {
		if err := sendmail.Send(m, config); err != nil {
			fatal("Failed to send:", err)
		}
	}
	log.Println("Done.")
}

func getConfPath(p *string) string {
	if *p != "" {
		return *p
	}

	config := os.Getenv("XDG_CONFIG_HOME")
	if config != "" {
		return path.Join(config, defaultConfig)
	}

	home := os.Getenv("HOME")
	return path.Join(home, ".config", defaultConfig)
}
