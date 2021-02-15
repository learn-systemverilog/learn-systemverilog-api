package transpiler

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
)

// Transpile Transpile your code from SystemVerilog to C++.
func Transpile(code string, logs chan<- interface{}) (string, error) {
	defer close(logs)

	logs <- logInternal("Creating temporary workspace.", logInternalSeverityInfo)

	workspace, err := setupTempWorkspace(code)
	if err != nil {
		logs <- logInternal(err.Error(), logInternalSeverityError)

		return "", err
	}
	defer func() {
		err := os.RemoveAll(workspace)
		if err != nil {
			logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityDebug)
		}
	}()

	logs <- logInternalWorkspace(
		"Transpiling the code from SystemVerilog to C++.",
		workspace,
		logInternalSeverityInfo,
	)

	if err := transpileSVToCPP(workspace, logs); err != nil {
		return "", fmt.Errorf("transpiling from sv to cpp: %w", err)
	}

	if err := transpileCPPToJS(workspace, logs); err != nil {
		return "", fmt.Errorf("transpiling from cpp to js: %w", err)
	}

	output, err := getOutput(workspace)
	if err != nil {
		return "", fmt.Errorf("reading output: %w", err)
	}

	return string(output), nil
}

func setupTempWorkspace(code string) (workspace string, err error) {
	log.Println("Creating temp workspace...")
	workspace, err = ioutil.TempDir("", "lsv_api_transpiler_workspace_")
	if err != nil {
		return "", fmt.Errorf("creating temp workspace: %w", err)
	}

	defer func() {
		if err != nil {
			err := os.RemoveAll(workspace)
			if err != nil {
				log.Printf("Error while trying to remove workspace: %v", err)
			}
		}
	}()

	log.Println("Temp workspace created:", workspace)

	log.Println("Copying template workspace.")
	err = copyDir(workSpaceTemplatePath, workspace)
	if err != nil {
		return "", fmt.Errorf("copying template workspace: %w", err)
	}

	log.Println("Template workspace copied successfully.")

	topSV := path.Join(workspace, mainSVFileName)

	log.Println("Creating main user code file:", topSV)
	err = ioutil.WriteFile(topSV, []byte(code), 0600)
	if err != nil {
		return "", fmt.Errorf("creating main user code file: %w", err)
	}

	log.Println("Main user code file created successfully.")

	return
}

func transpileSVToCPP(workspace string, logs chan<- interface{}) error {
	logs <- logInternal("Transpiling from SystemVerilog to C++.", logInternalSeverityInfo)

	return runMakeTarget("obj_dir", workspace, logs)
}

func transpileCPPToJS(workspace string, logs chan<- interface{}) error {
	logs <- logInternal("Transpiling from C++ to JavaScript.", logInternalSeverityInfo)

	return runMakeTarget("simulator.js", workspace, logs)
}

func getOutput(workspace string) ([]byte, error) {
	outputPath := filepath.Join(workspace, outputJSFileName)

	output, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func runMakeTarget(target, workspace string, logs chan<- interface{}) error {
	cmd := exec.Command("make", target)
	cmd.Dir = workspace

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityError)

		return fmt.Errorf("stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityError)

		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityError)

		return fmt.Errorf("starting command: %w", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logs <- newLogStdout(scanner.Text())
		}
		if scanner.Err() != nil {
			logs <- logInternalWorkspace(scanner.Err().Error(), workspace, logInternalSeverityWarn)
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			logs <- newLogStderr(scanner.Text())

		}
		if scanner.Err() != nil {
			logs <- logInternalWorkspace(scanner.Err().Error(), workspace, logInternalSeverityWarn)
		}

		wg.Done()
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityError)

			return fmt.Errorf("waiting command: %w", err)
		}

		// Do we need to handle this error?
		logs <- logInternalWorkspace(err.Error(), workspace, logInternalSeverityWarn)
	}

	return nil
}

func logInternal(msg, severity string) LogInternal {
	log.Println(msg)

	return newLogInternal(msg, severity)
}

func logInternalWorkspace(msg, workspace, severity string) LogInternal {
	log.Println(workspace+":", msg)

	return newLogInternal(msg, severity)
}

// copyDir Copy the contents of a src directory to a dst one.
// Obs. 1: It does not copy recursively.
// Obs. 2: The method of copying individual files is a bit expensive for larges ones.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return errors.New("src must be a directory")
	}

	dstInfo, err := os.Stat(dst)
	if err != nil {
		return err
	}

	if !dstInfo.IsDir() {
		return errors.New("dst must be a directory")
	}

	err = filepath.Walk(src, func(srcFilename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if os.SameFile(srcInfo, info) {
			return nil
		}

		if info.IsDir() {
			return fmt.Errorf("recursive copying is not supported: %s", info.Name())
		}

		fileData, err := ioutil.ReadFile(srcFilename)
		if err != nil {
			return err
		}

		dstFilename := path.Join(dst, info.Name())

		return ioutil.WriteFile(dstFilename, fileData, 0600)
	})

	return err
}
