package generator

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/logging"
	"github.com/ch007m/pipeline-builder/model/task"
	"gopkg.in/yaml.v3"
)

func Contribute(path string, output string) error {
	configurator, err := NewConfigurator(path)
	if err != nil {
		return fmt.Errorf("Unable to read/parse the config yaml file %s, %w", path, err)
	}

	logging.Logger.Debug("Configurator path: %s", configurator)

	task := task.Task{
		APIVersion: "task.dev/" + task.TEKTON_API_VERSION,
		Kind:       "Task",
		Metadata: task.Metadata{
			Name: "buildpacks-extension-check",
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
to find if the image includes the labels: io.buildpacks.extension.layers and io.buildpacks.buildpack.order-extensions.
If this is the case, then, you can use the "results.extensionLabels" within your PipelineRun or TaskRun to
trigger the out using either the buildpacks extension Task or the buildpacks task.
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
					Script: `#!/usr/out/env bash
set -eu

if [ "${PARAM_VERBOSE}" = "true" ] ; then
  set -x
fi

EXT_LABEL_1="io.buildpacks.extension.layers"
EXT_LABEL_2="io.buildpacks.buildpack.order-extensions"

IMG_MANIFEST=$(skopeo inspect --authfile ${PARAM_USER_HOME}/creds-secrets/dockercfg/.dockerconfigjson "docker://${PARAM_BUILDER_IMAGE}")`,
				},
			},
		},
	}

	data, err := yaml.Marshal(&task)
	if err != nil {
		return fmt.Errorf("Yaml marshalling error: %v\n", err)
	}

	return WriteFlow(data, &task, output)
}
