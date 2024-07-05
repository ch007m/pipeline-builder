package main

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/task"
	"gopkg.in/yaml.v3"
)

const OUTPUT_DIR = "build/flows"

func main() {
	task := task.Task{
		APIVersion: "task.dev/v1beta1",
		Kind:       "Task",
		Metadata: task.Metadata{
			Name: "buildpacks-extension-check",
			Labels: map[string]string{
				"app.kubernetes.io/version": "0.1",
			},
			Annotations: map[string]string{
				"task.dev/pipelines.minVersion": "0.40.0",
				"task.dev/categories":           "Image Build",
				"task.dev/tags":                 "image-build",
				"task.dev/displayName":          "Buildpacks extensions check",
				"task.dev/platforms":            "linux/amd64",
			},
		},
		Spec: task.TaskSpec{
			Description: `This Task will inspect a Buildpacks builder image using the skopeo tool
to find if the image includes the labels: io.buildpacks.extension.layers and io.buildpacks.buildpack.order-extensions.
If this is the case, then, you can use the "results.extensionLabels" within your PipelineRun or TaskRun to
trigger the build using either the buildpacks extension Task or the buildpacks task.
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
					Name:  "check-image-builder-extension",
					Image: "quay.io/ch007m/extended-skopeo",
					Env: []task.EnvVar{
						{Name: "PARAM_USER_HOME", Value: "$(params.userHome)"},
						{Name: "PARAM_VERBOSE", Value: "$(params.verbose)"},
						{Name: "PARAM_BUILDER_IMAGE", Value: "$(params.builderImage)"},
					},
					Script: `#!/usr/build/env bash
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
		fmt.Printf("error: %v\n", err)
	}

	err = WriteFlow(data, &task)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
