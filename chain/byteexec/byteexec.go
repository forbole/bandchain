package byteexec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
)

var (
	fileMode = os.FileMode(0744)
)

func writeFile(executable []byte) (string, string, error) {
	dir, err := ioutil.TempDir("/tmp", "temp")
	if err != nil {
		return "", "", err
	}
	filename := filepath.Join(dir, "exec")
	err = ioutil.WriteFile(filename, executable, fileMode)
	if err != nil {
		return "", "", err
	}
	filename, err = filepath.Abs(filename)
	if err != nil {
		return "", "", err
	}
	return dir, filename, nil
}

// RunOnLocal spawns a new subprocess and runs the given executable. NOT SAFE!
func RunOnLocal(executable []byte, timeOut time.Duration, arg string) ([]byte, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return nil, err
	}

	dir, filename, err := writeFile(executable)
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dir) // clean up

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)

	defer cancel()

	return exec.CommandContext(ctx, filename, args...).Output()
}

// RunOnDocker runs the given executable in a new docker container.
func RunOnDocker(executable []byte, sandboxMode bool, timeOut time.Duration, arg string) ([]byte, error) {
	args, err := shellwords.Parse(arg)
	if err != nil {
		return nil, err
	}

	dir, filename, err := writeFile(executable)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir) // clean up

	commands := []string{"run"}
	if sandboxMode {
		commands = append(commands, "--runtime=runsc")
	}
	commands = append(
		commands, "-d", "--rm", "band-provider", "sleep", fmt.Sprintf("%d", int(timeOut.Seconds())),
	)
	rawID, err := exec.Command("docker", commands...).Output()

	if err != nil {
		return nil, err
	}

	containerID := strings.TrimSpace(string(rawID))
	defer exec.Command("docker", "stop", containerID).Output()

	_, err = exec.Command(
		"docker", "cp", filename, fmt.Sprintf("%s:/exec", containerID),
	).Output()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	newArgs := append([]string{"exec", containerID, "./exec"}, args...)

	return exec.CommandContext(ctx, "docker", newArgs...).Output()
}

// RunOnAWSLambda runs the given executable on AMS Lambda platform.
func RunOnAWSLambda(executable []byte, timeOut time.Duration, arg string) ([]byte, error) {
	requestBody, err := json.Marshal(map[string]string{
		"executable": string(executable),
		"calldata":   arg,
	})

	fmt.Println("--------- 1")

	request, err := http.NewRequest("POST", "https://dmptasv4j8.execute-api.ap-southeast-1.amazonaws.com/bash-execute", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	fmt.Println("--------- 2")

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: timeOut,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("--------- 3")

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("--------- 4")

	type result struct {
		Returncode int    `json:"returncode"`
		Stdout     string `json:"stdout"`
		Stderr     string `json:"stderr"`
	}

	r := result{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}

	fmt.Println(r)

	return []byte(r.Stdout), nil
}
