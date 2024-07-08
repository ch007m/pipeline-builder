package util

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/model/pipeline"
	"github.com/disiqueira/gotree"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rakyll/statik/fs"
)

func WriteFlow(content []byte, pipeline *pipeline.Pipeline, output_dir string) error {
	t := gotree.New("pipeline")
	t.Add(string(content))

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("Unable to get the current path !")
	}
	file := filepath.Join(path, output_dir, pipeline.Metadata.Name+".yaml")

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return fmt.Errorf("unable to create %s\n%w", filepath.Dir(file), err)
	}

	if err := os.WriteFile(file, content, 0644); err != nil {
		return fmt.Errorf("unable to write %s\n%w", file, err)
	}

	fmt.Println(t.Print())
	return nil
}

func StatikString(path string) string {
	statik, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := statik.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}
