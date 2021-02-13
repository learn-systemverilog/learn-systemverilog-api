package transpiler

import (
	"io/ioutil"
	"log"
	"path"
)

// Run ...
func Run(code string) error {
	log.Println("Creating temp workspace...")
	workspace, err := ioutil.TempDir("", "lsv_api_transpiler_workspace_")
	if err != nil {
		return err
	}

	log.Println("Temp workspace created:", workspace)

	topSV := path.Join(workspace, mainSVFileName)

	log.Println("Creating "+mainSVFileName+":", topSV)
	err = ioutil.WriteFile(topSV, []byte(code), 0600)
	if err != nil {
		return err
	}

	log.Println("top.sv created successfully.")

	return nil
}
