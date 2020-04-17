package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/manifoldco/promptui"
	yaml "gopkg.in/yaml.v2"
)

// Server config structure
type Server struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// config file structure
type YamlServers struct {
	Servers []Server `yaml:"servers"`
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "f", "", "Select server list")
	flag.Parse()

	if fileName == "" {

		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
			return
		}

		if _, err := os.Stat(fmt.Sprintf("%s/ssh.default.yaml", usr.HomeDir)); os.IsNotExist(err) {
			fmt.Println("Please provide yaml file by using -f option")
			return
		}

		fileName = fmt.Sprintf("%s/ssh.default.yaml", usr.HomeDir)
	}

	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var Config YamlServers
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `{{ "\u2713" }} {{ .Name | cyan }} ({{ .Username }}{{ "@" }}{{ .Host }})`,
		Inactive: `  {{ .Name | cyan }} ({{ .Username }}{{ "@" }}{{ .Host }})`,
		Selected: `Connecting to {{ .Username }}{{ "@" }}{{ .Host }} via SSH:`,
		Details: `
    --------- Server ----------
    {{ "Name:" | faint }}	{{ .Name }}
    {{ "Host:" | faint }}	{{ .Host }}
    {{ "Username:" | faint }}	{{ .Username }}
    {{ "Password:" | faint }}	{{ .Password }}`,
	}

	searcher := func(input string, index int) bool {
		server := Config.Servers[index]
		name := strings.Replace(strings.ToLower(server.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Connect to specific server, please select one to enter SSH.",
		Items:     Config.Servers,
		Templates: templates,
		Size:      15,
		Searcher:  searcher,
	}

	index, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	selected := Config.Servers[index]
	cmd := exec.Command("sshpass", "-p", selected.Password, "ssh", fmt.Sprintf("%s@%s", selected.Username, selected.Host))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}
