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

// Run ...
func Run(code string, msgs chan<- map[string]string) {
	defer close(msgs)

	msgs <- map[string]string{
		"type":    msgTypeInfo,
		"message": "Creating temporary workspace.",
	}

	workspace, err := setupTempWorkspace(code)
	if err != nil {
		msgs <- map[string]string{
			"type":    msgTypeError,
			"message": err.Error(),
		}

		return
	}
	defer func() {
		err := os.RemoveAll(workspace)
		if err != nil {
			log.Printf("Error while trying to remove workspace: %v", err)
		}
	}()

	msgs <- map[string]string{
		"type":    msgTypeInfo,
		"message": "Transpiling the code from SystemVerilog to C++.",
	}

	if ok := transpileSVToCPP(workspace, msgs); !ok {
		return
	}
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

func transpileSVToCPP(workspace string, msgs chan<- map[string]string) bool {
	cmd := exec.Command("make", "obj_dir")
	cmd.Dir = workspace

	cmd.Stdin = nil

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		msgs <- map[string]string{
			"type":    msgTypeError,
			"message": err.Error(),
		}

		return false
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		msgs <- map[string]string{
			"type":    msgTypeError,
			"message": err.Error(),
		}

		return false
	}

	if err := cmd.Start(); err != nil {
		msgs <- map[string]string{
			"type":    msgTypeError,
			"message": err.Error(),
		}

		return false
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			msgs <- map[string]string{
				"type":    msgTypeStdout,
				"message": scanner.Text(),
			}
		}
		if scanner.Err() != nil {
			msgs <- map[string]string{
				"type":    msgTypeWarning,
				"message": scanner.Err().Error(),
			}
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			msgs <- map[string]string{
				"type":    msgTypeStderr,
				"message": scanner.Text(),
			}
		}
		if scanner.Err() != nil {
			msgs <- map[string]string{
				"type":    msgTypeWarning,
				"message": scanner.Err().Error(),
			}
		}

		wg.Done()
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if errors.Is(err, &exec.ExitError{}) {
			msgs <- map[string]string{
				"type":    "exit",
				"message": err.Error(),
			}

			return false
		}

		msgs <- map[string]string{
			"type":    msgTypeWarning,
			"message": fmt.Sprintf("Waiting for the command: %v", err),
		}
	}

	return true
}
