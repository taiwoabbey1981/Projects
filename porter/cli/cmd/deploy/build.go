package deploy

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	api "github.com/porter-dev/porter/api/client"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/cli/cmd/docker"
	"github.com/porter-dev/porter/cli/cmd/pack"
)

// BuildAgent builds a new Docker container image for a new version of an application
type BuildAgent struct {
	*SharedOpts

	APIClient   api.Client
	ImageRepo   string
	Env         map[string]string
	ImageExists bool
}

// BuildDocker uses the local Docker daemon to build the image
func (b *BuildAgent) BuildDocker(
	ctx context.Context,
	dockerAgent *docker.Agent,
	basePath,
	buildCtx,
	dockerfilePath,
	tag string,
	currentTag string,
) error {
	buildCtx, dockerfilePath, isDockerfileInCtx, err := ResolveDockerPaths(
		basePath,
		buildCtx,
		dockerfilePath,
	)
	if err != nil {
		return err
	}

	opts := &docker.BuildOpts{
		ImageRepo:         b.ImageRepo,
		Tag:               tag,
		CurrentTag:        currentTag,
		BuildContext:      buildCtx,
		Env:               b.Env,
		DockerfilePath:    dockerfilePath,
		IsDockerfileInCtx: isDockerfileInCtx,
		UseCache:          b.UseCache,
	}

	return dockerAgent.BuildLocal(
		ctx,
		opts,
	)
}

// BuildPack uses the cloud-native buildpack client to build a container image
func (b *BuildAgent) BuildPack(ctx context.Context, dockerAgent *docker.Agent, dst, tag, prevTag string, buildConfig *types.BuildConfig) error {
	// retag the image with "pack-cache" tag so that it doesn't re-pull from the registry
	if b.ImageExists {
		err := dockerAgent.TagImage(
			ctx,
			fmt.Sprintf("%s:%s", b.ImageRepo, prevTag),
			fmt.Sprintf("%s:%s", b.ImageRepo, "pack-cache"),
		)
		if err != nil {
			return err
		}
	}

	// create pack agent and build opts
	packAgent := &pack.Agent{}

	opts := &docker.BuildOpts{
		ImageRepo:    b.ImageRepo,
		Tag:          tag,
		BuildContext: dst,
		Env:          b.Env,
		UseCache:     b.UseCache,
	}

	// call builder
	return packAgent.Build(ctx, opts, buildConfig, fmt.Sprintf("%s:%s", b.ImageRepo, "pack-cache"))
}

// ResolveDockerPaths returns a path to the dockerfile that is either relative or absolute, and a path
// to the build context that is absolute.
//
// The return value will be relative if the dockerfile exists within the build context, absolute
// otherwise. The second return value is true if the dockerfile exists within the build context,
// false otherwise.
func ResolveDockerPaths(
	basePath string,
	buildContextPath string,
	dockerfilePath string,
) (
	resBuildCtxPath string,
	resDockerfilePath string,
	isDockerfileRelative bool,
	err error,
) {
	resBuildCtxPath, err = filepath.Abs(buildContextPath)
	resDockerfilePath = dockerfilePath

	// determine if the given dockerfile path is relative
	if !filepath.IsAbs(dockerfilePath) {
		// if path is relative, join basepath with path
		resDockerfilePath = filepath.Join(basePath, dockerfilePath)
	}

	// compare the path to the dockerfile with the build context
	pathComp, err := filepath.Rel(resBuildCtxPath, resDockerfilePath)
	if err != nil {
		return "", "", false, err
	}

	if !strings.HasPrefix(pathComp, ".."+string(os.PathSeparator)) {
		// return the relative path to the dockerfile
		return resBuildCtxPath, pathComp, true, nil
	}

	resDockerfilePath, err = filepath.Abs(resDockerfilePath)

	if err != nil {
		return "", "", false, err
	}

	return resBuildCtxPath, resDockerfilePath, false, nil
}
