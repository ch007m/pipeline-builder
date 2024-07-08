package lifecycle

import (
	"github.com/ch007m/pipeline-builder/model/common"
	"github.com/ch007m/pipeline-builder/model/pipeline"
)

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
					TaskRef: pipeline.TaskRef{
						Resolver: "git",
						Params: []pipeline.Param{
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
					TaskRef: pipeline.TaskRef{
						Resolver: "bundles",
						Params: []pipeline.Param{
							{Name: "bundle", Type: "string", Default: "quay.io/redhat-appstudio-tekton-catalog/task-summary:0.1@sha256:e69f53a3991d7088d8aa2827365ab761ab7524d4269f296b4a78b0f085789d30"},
							{Name: "name", Type: "string", Default: "summary"},
							{Name: "kind", Type: "string", Default: "Task"},
						},
					},
				},
			},
			Tasks: []pipeline.Task{
				// TODO : Add Tasks
				{
					Name: "",
				},
			},
		},
	}
	return pipeline
}
