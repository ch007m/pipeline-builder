package generator

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/model/task"
	"github.com/disiqueira/gotree"
	"os"
	"path/filepath"
)

func WriteFlow(content []byte, task *task.Task, output_dir string) error {
	t := gotree.New("pipeline")
	t.Add(string(content))

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return fmt.Errorf("Unable to get the current path !")
	}
	file := filepath.Join(path, output_dir, task.Metadata.Name+".yaml")

	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return fmt.Errorf("unable to create %s\n%w", filepath.Dir(file), err)
	}

	if err := os.WriteFile(file, content, 0644); err != nil {
		return fmt.Errorf("unable to write %s\n%w", file, err)
	}

	fmt.Println(t.Print())
	return nil
}
