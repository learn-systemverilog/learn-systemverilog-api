package transpiler

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
)

// Run ...
func Run(code string) error {
	log.Println("Creating temp workspace...")
	workspace, err := ioutil.TempDir("", "lsv_api_transpiler_workspace_")
	if err != nil {
		return fmt.Errorf("creating temp workspace: %w", err)
	}

	log.Println("Temp workspace created:", workspace)

	topSV := path.Join(workspace, mainSVFileName)

	log.Println("Creating main user code file:", topSV)
	err = ioutil.WriteFile(topSV, []byte(code), 0600)
	if err != nil {
		return fmt.Errorf("creating main user code file: %w", err)
	}

	log.Println("Main user code file created successfully.")

	return nil
}
