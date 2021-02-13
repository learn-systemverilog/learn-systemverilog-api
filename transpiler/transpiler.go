package transpiler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// Run ...
func Run(code string) error {
	log.Println("Creating temp workspace...")
	workspace, err := ioutil.TempDir("", "lsv_api_transpiler_workspace_")
	if err != nil {
		return fmt.Errorf("creating temp workspace: %w", err)
	}

	log.Println("Temp workspace created:", workspace)

	log.Println("Copying template workspace.")
	err = copyDir(workSpaceTemplatePath, workspace)
	if err != nil {
		return fmt.Errorf("copying template workspace: %w", err)
	}

	log.Println("Template workspace copied successfully.")

	topSV := path.Join(workspace, mainSVFileName)

	log.Println("Creating main user code file:", topSV)
	err = ioutil.WriteFile(topSV, []byte(code), 0600)
	if err != nil {
		return fmt.Errorf("creating main user code file: %w", err)
	}

	log.Println("Main user code file created successfully.")

	return nil
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

		if err := ioutil.WriteFile(dstFilename, fileData, 0600); err != nil {
			return err
		}

		return nil
	})

	return err
}
