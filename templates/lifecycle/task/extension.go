package lifecycle

import (
	"github.com/ch007m/pipeline-builder/model/common"
	"github.com/ch007m/pipeline-builder/model/task"
	_ "github.com/ch007m/pipeline-builder/templates/lifecycle/task/statik"
	"github.com/ch007m/pipeline-builder/util"
)

//go:generate statik -src . -include *.sh

func CreateExtensionTask() task.Task {
	task := task.Task{
		APIVersion: "task.dev/" + common.TEKTON_API_VERSION,
		Kind:       "Task",
		Metadata: common.Metadata{
			Name: "buildpack-extension-check",
			Labels: map[string]string{
				"app.kubernetes.io/version": "0.1",
			},
			Annotations: map[string]string{
				"task.dev/pipelines.minVersion": "0.40.0",
				"task.dev/categories":           "Image Build",
				"task.dev/tags":                 "image-out",
				"task.dev/displayName":          "Buildpacks extensions check",
				"task.dev/platforms":            "linux/amd64",
			},
		},
		Spec: task.TaskSpec{
			Description: `This Task will inspect a Buildpacks generator image using the skopeo tool
to find if the image includes the labels: io.buildpack.extension.layers and io.buildpack.templates.order-extensions.
If this is the case, then, you can use the "results.extensionLabels" within your PipelineRun or TaskRun to
trigger the out using either the buildpack extension Task or the buildpack task.
Additionally, the CNB USER ID and CNB GROUP ID of the image will be exported as results.`,
			Params: []task.Param{
				{
					Name:        "userHome",
					Description: "Absolute path to the user's home directory.",
					Type:        "string",
					Default:     "/tekton/home",
				},
				{
					Name:        "verbose",
					Description: "Log the commands that are executed during `git-clone`'s operation.",
					Type:        "string",
					Default:     "false",
				},
				{
					Name:        "builderImage",
					Description: "Builder image to be scanned",
					Type:        "string",
				},
			},
			Results: []task.Result{
				{
					Name:        "uid",
					Description: "UID of the user specified in the image",
				},
				{
					Name:        "gid",
					Description: "GID of the user specified in the image",
				},
				{
					Name:        "extensionLabels",
					Description: "Extensions labels defined in the image",
				},
			},
			Steps: []task.Step{
				{
					Name:  "check-image-generator-extension",
					Image: "quay.io/ch007m/extended-skopeo",
					Env: []task.EnvVar{
						{Name: "PARAM_USER_HOME", Value: "$(params.userHome)"},
						{Name: "PARAM_VERBOSE", Value: "$(params.verbose)"},
						{Name: "PARAM_BUILDER_IMAGE", Value: "$(params.builderImage)"},
					},
					Script: util.StatikString("/extension.sh"),
				},
			},
		},
	}
	return task
}
