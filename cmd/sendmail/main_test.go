package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

var (
	tmpfile    *os.File
	binaryName = "sendmail"
)

// Test setup: build binary and create temp config file
func TestMain(m *testing.M) {

	// build binary
	make := exec.Command("go", "build", "github.com/kerma/sendmail/cmd/sendmail")
	err := make.Run()
	if err != nil {
		fmt.Printf("could not make binary for %s: %v\n", binaryName, err)
		os.Exit(1)
	}

	// set up conf file
	tmpfile, err = ioutil.TempFile("", "conf.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte("{\"server\": \"smtp.example.com\"}")
	if _, err := tmpfile.Write(content); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := tmpfile.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// run tests
	os.Exit(m.Run())
}

// Run parameterized integration tests
func TestRunner(t *testing.T) {

	tests := []struct {
		name  string
		stdin string
		args  []string
		error string
	}{
		{
			"Happy path",
			"email body",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-to", "hello@test.com",
				"-cc", "copy@test.com, second@test.com",
				"-subject", "test"},
			"",
		},
		{
			"To via prompt",
			"hello@test.com\nemail body",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-subject", "test"},
			"",
		},
		{
			"Invalid To address",
			"not_an_email.com",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-subject", "test"},
			"exit status 1",
		},
		{
			"Invalid CC address",
			"email body",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-to", "hello@test.com",
				"-cc", "is@email.com,not_an_email.com",
				"-subject", "test"},
			"exit status 1",
		},
		{
			"Server via flag",
			"email body",
			[]string{
				"-dryrun",
				"-server", "smtp.example.com",
				"-to", "hello@test.com",
				"-subject", "test"},
			"",
		},
		{
			"To, subject and body via stdin",
			"hello@example.com\nHello subject\nThis is email",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name()},
			"",
		},
		{
			"Attachment",
			"email body",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-to", "hello@test.com",
				"-subject", "test",
				"-attach", "filepath.txt"},
			"",
		},
		{
			"HTML body",
			"<h1>Hello</h1><p>world</p>",
			[]string{
				"-dryrun",
				"-conf", tmpfile.Name(),
				"-to", "hello@test.com",
				"-subject", "test",
				"-html"},
			"",
		},
	}

	d, _ := os.Getwd()
	command := path.Join(d, binaryName)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(command, tt.args...)

			stdin, err := cmd.StdinPipe()
			if err != nil {
				t.Fatal(err)
			}

			if tt.stdin != "" {
				go func() {
					defer stdin.Close()
					io.WriteString(stdin, tt.stdin)
				}()
			}

			out, err := cmd.CombinedOutput()
			if err != nil && err.Error() != tt.error {
				fmt.Printf("Error: %s\n", err.Error())
				fmt.Printf("Expected error: %s\n", tt.error)
				fmt.Printf("Output: %s\n", out)
				t.Fatal("Unexpected error")
			}
		})
	}
}
