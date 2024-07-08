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
			Name: "buildpack-phases",
			Labels: map[string]string{
				"app.kubernetes.io/version": "0.1",
			},
			Annotations: map[string]string{
				"task.dev/pipelines.minVersion": "0.40.0",
				"task.dev/categories":           "Image Build",
				"task.dev/tags":                 "image-build",
				"task.dev/displayName":          "Buildpacks phases",
				"task.dev/platforms":            "linux/amd64",
			},
		},
		Spec: task.TaskSpec{
			Description: `The Buildpacks-Phases task builds source into a container image and pushes it to
    a registry, using Cloud Native Buildpacks. This task separately calls the aspects of the
    Cloud Native Buildpacks lifecycle, to provide increased security via container isolation.`,
			Workspaces: []task.Workspace{
				{Name: "source", Description: "Directory where application source is located."},
				{Name: "cache", Description: "Directory where cache is stored (when no cache image is provided).", Optional: true},
			},
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
			StepTemplate: task.StepTemplate{
				[]task.EnvVar{
					{Name: "CNB_PLATFORM_API", Value: "0.12"},
					{Name: "CNB_EXPERIMENTAL_MODE", Value: "$(params.CNB_EXPERIMENTAL_MODE)"},
					{Name: "HOME", Value: "$(params.USER_HOME)"},
				},
			},
			Volumes: []task.Volume{
				{Name: "empty-dir", EmptyDir: "{}"},
				{Name: "layers-dir", EmptyDir: "{}"},
				{Name: "kaniko-dir", EmptyDir: "{}"},
			},
			Steps: []task.Step{
				{
					Name:  "prepare",
					Image: "quay.io/swsmirror/bash",
					Args: []string{
						"--env-vars",
						"$(params.CNB_ENV_VARS[*])",
					},
					Script: util.StatikString("/create_folder_set_permissions.sh"),
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
						{Name: "$(params.PLATFORM_DIR)", MountPath: "/platform"},
					},
				},
				{
					Name:            "analyze",
					Image:           "$(params.CNB_BUILDER_IMAGE)",
					ImagePullPolicy: "Always",
					Command: []string{
						"/cnb/lifecycle/analyzer",
					},
					Args: []string{
						"-log-level=$(params.CNB_LOG_LEVEL)",
						"-layers=/layers",
						"-run-image=$(params.CNB_RUN_IMAGE)",
						"-cache-image=$(params.CACHE_IMAGE)",
						"-uid=$(params.CNB_USER_ID)",
						"-gid=$(params.CNB_GROUP_ID)",
						"$(params.APP_IMAGE)",
					},
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
					},
				},
				{
					Name:            "detect",
					Image:           "$(params.CNB_BUILDER_IMAGE)",
					ImagePullPolicy: "Always",
					Command: []string{
						"/cnb/lifecycle/detector",
					},
					Args: []string{
						"-log-level=$(params.CNB_LOG_LEVEL)",
						"-app=$(workspaces.source.path)/$(params.SOURCE_SUBPATH)",
						"-group=/layers/group.toml",
						"-plan=/layers/plan.toml",
					},
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
						{Name: "$(params.PLATFORM_DIR)", MountPath: "/platform"},
						{Name: "empty-dir", MountPath: "/tekton/home"},
					},
				},
				{
					Name:            "restore",
					Image:           "$(params.CNB_BUILDER_IMAGE)",
					ImagePullPolicy: "Always",
					Command: []string{
						"/cnb/lifecycle/restorer",
					},
					Args: []string{
						"-log-level=$(params.CNB_LOG_LEVEL)",
						"-build-image=$(params.CNB_BUILD_IMAGE)",
						"-group=/layers/group.toml",
						"-layers=/layers",
						"-cache-dir=$(workspaces.cache.path)",
						"-cache-image=$(params.CACHE_IMAGE)",
						"-uid=$(params.CNB_USER_ID)",
						"-gid=$(params.CNB_GROUP_ID)",
					},
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
						{Name: "kaniko-dir", MountPath: "/kaniko"},
					},
				},
				{
					Name:            "build",
					Image:           "$(params.CNB_BUILDER_IMAGE)",
					ImagePullPolicy: "Always",
					Command: []string{
						"/cnb/lifecycle/builder",
					},
					Args: []string{
						"-log-level=$(params.CNB_LOG_LEVEL)",
						"-app=$(workspaces.source.path)/$(params.SOURCE_SUBPATH)",
						"-layers=/layers",
						"-group=/layers/group.toml",
						"-plan=/layers/plan.toml",
					},
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
						{Name: "$(params.PLATFORM_DIR)", MountPath: "/platform"},
						{Name: "empty-dir", MountPath: "/tekton/home"},
					},
				},
				{
					Name:            "export",
					Image:           "$(params.CNB_BUILDER_IMAGE)",
					ImagePullPolicy: "Always",
					Command: []string{
						"/cnb/lifecycle/exporter",
					},
					Args: []string{
						"-log-level=$(params.CNB_LOG_LEVEL)",
						"-app=$(workspaces.source.path)/$(params.SOURCE_SUBPATH)",
						"-layers=/layers",
						"-group=/layers/group.toml",
						"-cache-dir=$(workspaces.cache.path)",
						"-cache-image=$(params.CACHE_IMAGE)",
						"-report=/layers/report.toml",
						"-process-type=$(params.PROCESS_TYPE)",
						"-uid=$(params.CNB_USER_ID)",
						"-gid=$(params.CNB_GROUP_ID)",
						"$(params.APP_IMAGE)",
					},
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
					},
				},
				{
					Name:   "results",
					Image:  "quay.io/swsmirror/bash",
					Script: util.StatikString("/image_url_digest.sh"),
					VolumeMounts: []task.VolumeMount{
						{Name: "layers-dir", MountPath: "/layers"},
					},
				},
			},
		},
	}
	return task
}
