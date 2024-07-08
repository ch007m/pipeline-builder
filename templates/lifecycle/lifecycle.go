package lifecycle

import (
	"github.com/ch007m/pipeline-builder/model/common"
	"github.com/ch007m/pipeline-builder/model/pipeline"
)

// See: https://raw.githubusercontent.com/redhat-buildpacks/catalog/main/tekton/pipeline/rhtap/01/pipeline-buildpacks.yaml
func CreatePipeline() pipeline.Pipeline {
	pipeline := pipeline.Pipeline{
		APIVersion: "tekton.dev" + common.TEKTON_API_VERSION,
		Kind:       "Pipeline",
		Metadata: common.Metadata{
			Name: "pipeline-konflux",
		},
		Spec: pipeline.Spec{
			Workspaces: []pipeline.Workspace{
				{Name: "workspace"},
				{Name: "git-auth", Optional: true},
			},
			Params: []pipeline.Param{
				{Name: "git-url", Type: "string", Description: "Source Repository URL"},
				{Name: "revision", Type: "string", Description: "Revision of the Source Repository", Default: ""},
				{Name: "output-image", Type: "string", Description: "Fully Qualified Output Image"},
				{Name: "path-context", Type: "string", Description: "The path to your source code", Default: "."},
				{Name: "dockerfile", Type: "string", Description: "Path to the Dockerfile", Default: "Dockerfile"},
				{Name: "rebuild", Type: "string", Description: "Force rebuild image", Default: "false"},
				{Name: "enable-sbom", Type: "string", Description: "A boolean indicating whether we would like to execute a step", Default: "true"},
				{Name: "grype-sbom-format", Type: "string", Description: "Format to be used to export/show the SBOM. Format available for grype are 'json table cyclonedx cyclonedx-json sarif template'", Default: "table"},
				{Name: "skip-checks", Type: "string", Description: "Skip checks against built image", Default: "false"},
				{Name: "skip-optional", Type: "string", Description: "Skip optional checks, set false if you want to run optional checks", Default: "true"},
				{Name: "hermetic", Type: "string", Description: "Execute the build with network isolation", Default: "false"},
				{Name: "prefetch-input", Type: "string", Description: "Build dependencies to be prefetched by Cachi2", Default: ""},
				{Name: "java", Type: "string", Description: "Java build", Default: "false"},
				{Name: "snyk-secret", Type: "string", Description: "Snyk Token Secret Name", Default: ""},
				{Name: "image-expires-after", Description: "Image tag expiration time, time values could be something like 1h, 2d, 3w for hours, days, and weeks, respectively.", Default: ""},
				{Name: "sourceSubPath", Type: "string", Description: "Subpath of the git cloned project where code should be used", Default: "."},
				{Name: "cnbBuilderImage", Type: "string", Description: "Buildpacks Builder image to be used to build the runtime", Default: ""},
				{Name: "cnbLifecycleImage", Type: "string", Description: "TODO", Default: ""},
				{Name: "cnbBuildImage", Type: "string", Description: "TODO", Default: ""},
				{Name: "cnbRunImage", Type: "string", Description: "TODO", Default: ""},
				// TODO: Check how the YAML could be rendered like : default: [""]
				{Name: "cnbBuildEnvVars", Type: "array", Description: "Environment variables to set during _build-time_.", Default: []string{""}},
				{Name: "cnbLogLevel", Type: "string", Description: "TODO", Default: "info"},
				{Name: "cnbExperimentalMode", Type: "string", Description: "TODO", Default: "true"},
				{Name: "AppImage", Type: "string", Description: "TODO", Default: ""},
			},
			Results: []pipeline.Result{
				{Name: "IMAGE_URL", Description: "The URL of the built `APPLICATION_IMAGE`", Value: "$(tasks.build-container.results.IMAGE_URL)"},
				{Name: "IMAGE_DIGEST", Description: "The digest of the built `APPLICATION_IMAGE`", Value: "$(tasks.build-container.results.IMAGE_DIGEST)"},
				{Name: "CHAINS-GIT_URL", Description: "", Value: "$(tasks.clone-repository.results.url)"},
				{Name: "CHAINS-GIT_COMMIT", Description: "", Value: "$(tasks.clone-repository.results.commit)"},
			},
			Finally: []pipeline.Finally{
				{
					Name: "show-sbom",
					When: []pipeline.When{
						{Input: "$(params.enable-sbom)", Operator: "in", Values: []string{"true"}},
					},
					Params: []pipeline.Param{
						{Name: "GRYPE_IMAGE", Value: "anchore/grype:v0.65.1"},
						{Name: "ARGS", Value: []string{
							"$(tasks.build-container.results.IMAGE_URL)",
							"-o $(params.grype-sbom-format)",
						}},
					},
					TaskRef: pipeline.TaskRef{
						Resolver: "git",
						Params: []pipeline.Param{
							{Name: "url", Value: "https://github.com/tektoncd/catalog.git"},
							{Name: "revision", Value: "main"},
							{Name: "pathInRepo", Value: "task/grype/0.1/grype.yaml"},
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
						{Name: "pipelinerun-name", Value: "$(context.pipelineRun.name)"},
						{Name: "git-url", Value: "$(tasks.clone-repository.results.url)?rev=$(tasks.clone-repository.results.commit)"},
						{Name: "image-url", Value: "$(params.output-image)"},
						{Name: "build-task-status", Value: "$(tasks.build-container.status)"},
					},
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Value: "quay.io/redhat-appstudio-tekton-catalog/task-summary:0.1@sha256:e69f53a3991d7088d8aa2827365ab761ab7524d4269f296b4a78b0f085789d30"},
							{Name: "name", Value: "summary"},
							{Name: "kind", Value: "Task"},
						},
					},
				},
			},
			Tasks: []pipeline.Task{
				{
					Name: "init",
					Params: []pipeline.Param{
						{Name: "image-url", Value: "$(params.output-image)"},
						{Name: "rebuild", Value: "$(params.rebuild)"},
						{Name: "skip-checks", Value: "$(params.skip-checks)"},
						{Name: "skip-optional", Value: "$(params.skip-optional)"},
						{Name: "pipelinerun-name", Value: "$(context.pipelineRun.name)"},
						{Name: "pipelinerun-uid", Value: "$(context.pipelineRun.uid)"},
					},
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Value: "quay.io/redhat-appstudio-tekton-catalog/init:0.1"},
							{Name: "name", Value: "init"},
						},
					},
				},
				{
					Name: "clone-repository",
					Params: []pipeline.Param{
						{Name: "url", Value: "$(params.git-url)"},
						{Name: "revision", Value: "$(params.revision)"},
					},
					RunAfter: []string{"init"},
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Value: "quay.io/redhat-appstudio-tekton-catalog/task-git-clone:0.1@sha256:1f84973a21aabea38434b1f663abc4cb2d86565a9c7aae1f90decb43a8fa48eb"},
							{Name: "name", Value: "git-clone"},
						},
					},
					Workspaces: []pipeline.WorkspaceBinding{
						{Name: "output", Workspace: "workspace"},
					},
				},
				{
					Name:     "fetch-sources",
					RunAfter: []string{"init"},
					Params: []pipeline.Param{
						{Name: "url", Value: "$(params.git-url)"},
						{Name: "revision", Value: "$(params.revision)"},
						{Name: "subdirectory", Value: "$(params.sourceSubPath)"},
						{Name: "depth", Value: "0"},
					},
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Value: "quay.io/redhat-appstudio-tekton-catalog/git-clone:0.1"},
							{Name: "name", Value: "git-clone"},
							{Name: "kind", Value: "Task"},
						},
					},
					Workspaces: []pipeline.WorkspaceBinding{
						{Name: "output", Workspace: "workspace"},
					},
				},
				{
					Name:     "build-container",
					RunAfter: []string{"fetch-sources"},
					Params: []pipeline.Param{
						{Name: "IMAGE", Value: "$(params.output-image)"},
						{Name: "DOCKERFILE", Value: "$(params.dockerfile)"},
						{Name: "CONTEXT", Value: "$(params.path-context)"},
						{Name: "PREFETCH_INPUT", Value: "$(params.prefetch-input)"},
						{Name: "BUILDER_IMAGE", Value: "$(params.cnbBuilderImage)"},
						{Name: "LIFECYCLE_IMAGE", Value: "$(params.cnbLifecycleImage)"},
						{Name: "RUN_IMAGE", Value: "$(params.cnbRunImage)"},
						{Name: "PLATFORM_API_VERSION", Value: "0.8"},
						{Name: "BUILDPACKS", Value: "$(params.cnbBuildImage)"},
						{Name: "BUILDER_ENV_VARS", Value: "$(params.cnbBuildEnvVars)"},
						{Name: "EXPERIMENTAL", Value: "$(params.cnbExperimentalMode)"},
						{Name: "LOG_LEVEL", Value: "$(params.cnbLogLevel)"},
					},
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Value: "quay.io/redhat-appstudio-tekton-catalog/build:0.1"},
							{Name: "name", Value: "build"},
							{Name: "kind", Value: "Task"},
						},
					},
					Workspaces: []pipeline.WorkspaceBinding{
						{Name: "source", Workspace: "workspace"},
					},
				},
			},
		},
	}
	return pipeline
}
