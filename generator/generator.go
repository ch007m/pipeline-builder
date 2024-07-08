package generator

import (
	"fmt"
	"github.com/ch007m/pipeline-builder/logging"
	"github.com/ch007m/pipeline-builder/model/common"
	"github.com/ch007m/pipeline-builder/model/pipeline"
	"github.com/ch007m/pipeline-builder/model/task"

	"gopkg.in/yaml.v3"
)

func Contribute(path string, output string) error {
	configurator, err := NewConfigurator(path)
	if err != nil {
		return fmt.Errorf("Unable to read/parse the config yaml file %s, %w", path, err)
	}

	logging.Logger.Debug("Configurator path: %s", configurator)

	pipeline := pipeline.Pipeline{
		APIVersion: "tekton.dev/v1beta1",
		Kind:       "Pipeline",
		Metadata: common.Metadata{
			Name: "pipeline-rhtap",
		},
		Spec: pipeline.Spec{
			Workspaces: []pipeline.Workspace{
				{Name: "workspace"},
				{Name: "git-auth", Optional: true},
			},
			Params: []pipeline.Param{
				{Description: "Source Repository URL", Name: "git-url", Type: "string"},
				{Description: "Revision of the Source Repository", Name: "revision", Type: "string", Default: ""},
				{Description: "Fully Qualified Output Image", Name: "output-image", Type: "string"},
				{Description: "The path to your source code", Name: "path-context", Type: "string", Default: "."},
				{Description: "Path to the Dockerfile", Name: "dockerfile", Type: "string", Default: "Dockerfile"},
				{Description: "Force rebuild image", Name: "rebuild", Type: "string", Default: "false"},
				{Description: "A boolean indicating whether we would like to execute a step", Name: "enable-sbom", Type: "string", Default: "true"},
				{Description: "Format to be used to export/show the SBOM. Format available for grype are 'json table cyclonedx cyclonedx-json sarif template'", Name: "grype-sbom-format", Type: "string", Default: "table"},
				{Description: "Skip checks against built image", Name: "skip-checks", Type: "string", Default: "false"},
				{Description: "Skip optional checks, set false if you want to run optional checks", Name: "skip-optional", Type: "string", Default: "true"},
				{Description: "Execute the build with network isolation", Name: "hermetic", Type: "string", Default: "false"},
				{Description: "Build dependencies to be prefetched by Cachi2", Name: "prefetch-input", Type: "string", Default: ""},
				{Description: "Java build", Name: "java", Type: "string", Default: "false"},
				{Description: "Snyk Token Secret Name", Name: "snyk-secret", Type: "string", Default: ""},
				{Description: "Image tag expiration time, time values could be something like 1h, 2d, 3w for hours, days, and weeks, respectively.", Name: "image-expires-after", Default: ""},
				{Description: "Subpath of the git cloned project where code should be used", Name: "sourceSubPath", Type: "string", Default: "."},
				{Description: "Buildpacks Builder image to be used to build the runtime", Name: "cnbBuilderImage", Type: "string", Default: ""},
				{Description: "TODO", Name: "cnbLifecycleImage", Type: "string", Default: ""},
				{Description: "TODO", Name: "cnbBuildImage", Type: "string", Default: ""},
				{Description: "TODO", Name: "cnbRunImage", Type: "string", Default: ""},
				{Description: "Environment variables to set during _build-time_.", Name: "cnbBuildEnvVars", Type: "array", Default: `[""]`},
				{Description: "TODO", Name: "cnbLogLevel", Type: "string", Default: "info"},
				{Description: "TODO", Name: "cnbExperimentalMode", Type: "string", Default: "true"},
				{Description: "TODO", Name: "AppImage", Type: "string", Default: ""},
			},
			Results: []pipeline.Result{
				{Description: "The URL of the built `APPLICATION_IMAGE`", Name: "IMAGE_URL", Value: "$(tasks.build-container.results.IMAGE_URL)"},
				{Description: "The digest of the built `APPLICATION_IMAGE`", Name: "IMAGE_DIGEST", Value: "$(tasks.build-container.results.IMAGE_DIGEST)"},
				{Description: "", Name: "CHAINS-GIT_URL", Value: "$(tasks.clone-repository.results.url)"},
				{Description: "", Name: "CHAINS-GIT_COMMIT", Value: "$(tasks.clone-repository.results.commit)"},
			},
			Finally: []pipeline.Finally{
				{
					Name: "show-sbom",
					When: []pipeline.When{
						{Input: "$(params.enable-sbom)", Operator: "in", Values: []string{"true"}},
					},
					Params: []pipeline.Param{
						{Name: "GRYPE_IMAGE", Type: "string", Default: "anchore/grype:v0.65.1"},
						{Name: "ARGS", Type: "string", Default: "$(tasks.build-container.results.IMAGE_URL), -o $(params.grype-sbom-format)"},
					},
					TaskRef: task.TaskRef{
						Resolver: "git",
						Params: []task.Param{
							{Name: "url", Type: "string", Default: "https://github.com/tektoncd/catalog.git"},
							{Name: "revision", Type: "string", Default: "main"},
							{Name: "pathInRepo", Type: "string", Default: "task/grype/0.1/grype.yaml"},
						},
					},
					Workspaces: []pipeline.WorkspaceBinding{
						{Workspace: "workspace", Name: "source-dir"},
					},
				},
				{
					Name: "show-summary",
					When: []pipeline.When{
						{Input: "$(params.enable-sbom)", Operator: "in", Values: []string{"true"}},
					},
					Params: []pipeline.Param{
						{Name: "pipelinerun-name", Type: "string", Default: "$(context.pipelineRun.name)"},
						{Name: "git-url", Type: "string", Default: "$(tasks.clone-repository.results.url)?rev=$(tasks.clone-repository.results.commit)"},
						{Name: "image-url", Type: "string", Default: "$(params.output-image)"},
						{Name: "build-task-status", Type: "string", Default: "$(tasks.build-container.status)"},
					},
					TaskRef: task.TaskRef{
						Resolver: "bundles",
						Params: []task.Param{
							{Name: "bundle", Type: "string", Default: "quay.io/redhat-appstudio-tekton-catalog/task-summary:0.1@sha256:e69f53a3991d7088d8aa2827365ab761ab7524d4269f296b4a78b0f085789d30"},
							{Name: "name", Type: "string", Default: "summary"},
							{Name: "kind", Type: "string", Default: "Task"},
						},
					},
				},
			},
		},
	}

	task := task.Task{
		APIVersion: "task.dev/" + common.TEKTON_API_VERSION,
		Kind:       "Task",
		Metadata: common.Metadata{
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

	data, err := yaml.Marshal(&pipeline)
	if err != nil {
		return fmt.Errorf("Yaml marshalling error: %v\n", err)
	}

	return WriteFlow(data, &task, output)
}
