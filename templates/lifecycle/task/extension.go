package lifecycle

import (
	"github.com/ch007m/pipeline-builder/model/task"
	_ "github.com/ch007m/pipeline-builder/templates/lifecycle/task/statik"
)

//go:generate statik -src . -include *.sh

func CreateExtensionTask() task.Task {
	return nil
}
