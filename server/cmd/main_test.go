package main

import (
	"CloudScapes/pkg/logger"
	server "CloudScapes/server/internal"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestServer should run all nodejs integration tests for the server application
func TestServer(t *testing.T) {
	logFile := "server integration testing"
	if err := logger.InitLogger(true, &logFile); err != nil {
		panic(err)
	}
	defer logger.Flush()

	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Could not get working directory")
	}

	testFolder := fmt.Sprintf("%s/tests", cwd)

	testFiles := make([]string, 0)
	if err := filepath.Walk(testFolder, func(path string, f os.FileInfo, err error) error {
		// test files must end end with .js and not start with _
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".js") && !strings.HasPrefix(f.Name(), "_") {
			testFiles = append(testFiles, path)
		}
		return err
	}); err != nil {
		t.Fatal("Could not traverse test directory")
	}

	runTest := func(fullFilePath string) (string, error) {
		start := time.Now()
		filename := path.Base(fullFilePath)
		fmt.Printf("Testing '%s'\n", filename)

		cmd := exec.Command("node", fullFilePath)
		cmd.Dir = cwd
		cmd.Env = nil
		output, err := cmd.CombinedOutput()

		if err != nil {
			return "", fmt.Errorf("FAILED %s:\n\nError:\n\t%s\nResult:\n\t%s\n\n", filename, err, string(output))
		}

		return fmt.Sprintf("Finished %s - %s", filename, time.Since(start)), nil
	}

	for _, testFile := range testFiles {
		res, err := runTest(testFile)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res)
	}
}
