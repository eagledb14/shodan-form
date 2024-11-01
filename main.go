package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/eagledb14/form-scanner/types"
)

func main() {
	loadEnvVars()

	auto := flag.Bool("auto", false, "run in automatic mode")
	flag.Parse()

	if *auto {
		autoCreateEventFiles()
	} else {
		state := types.NewState()
		go openBrowser("localhost:8080")
		serv(":8080", &state)
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Println("Error opening browser:", err)
	}
}

func loadEnvVars() {
	file, err := os.Open("./resources/key.env")
	if err != nil {
		panic("error opening file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			panic("invalid line format")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := os.Setenv(key, value); err != nil {
			panic("error setting env var")
		}
	}

	if err := scanner.Err(); err != nil {
		panic("error reading file")
	}
}
