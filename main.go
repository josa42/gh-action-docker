package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	tagExp = regexp.MustCompile(`^refs/tags/(.*)$`)
)

type placeholder struct {
	Branch, Tag string
}

func main() {

	image := getInput("image")

	ref := os.Getenv("GITHUB_REF")

	label := ""

	gitTags := tagExp.FindStringSubmatch(ref)
	if len(gitTags) == 2 {
		label = gitTags[1]
	} else if ref == "refs/heads/master" {
		label = "latest"
	}

	log.Printf("ref: %s", ref)

	if label != "" {
		tag := fmt.Sprintf("%s:%s", image, label)

		dockerLogin()
		docker("build", ".", "--tag", tag)
		docker("push", tag)
	}
}

func docker(args ...string) {
	fmt.Println("$ docker", strings.Join(args, " "))

	cmd := exec.Command("docker", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Panicln(err)
	}
}

func dockerLogin() error {
	user := getInput("docker_username")
	pass := getInput("docker_password")

	cmd := exec.Command("docker", "login", "--username", user, "--password-stdin")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	io.WriteString(stdin, pass)
	stdin.Close()

	return cmd.Run()
}

func getInput(name string, defaultValues ...string) string {
	if value := os.Getenv(fmt.Sprintf("INPUT_%s", strings.ToUpper(name))); value != "" {
		return value
	}

	if len(defaultValues) == 1 {
		return defaultValues[0]
	}

	log.Panicf("Error: Input '%s' is required!\n", name)

	return ""
}
