//go:build ee

package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/helm"
	"github.com/porter-dev/porter/pkg/logger"
	"github.com/stefanmcshane/helm/pkg/release"
)

type RetryHelmAgent struct {
	form          *helm.Form
	l             *logger.Logger
	agent         *helm.Agent
	retryCount    uint
	retryInterval time.Duration
}

func NewRetryHelmAgent(
	ctx context.Context,
	form *helm.Form,
	l *logger.Logger,
	retryCount uint,
	retryInterval time.Duration,
) (*RetryHelmAgent, error) {
	if l == nil {
		l = logger.New(true, os.Stdout)
	}

	helmAgent, err := helm.GetAgentOutOfClusterConfig(ctx, form, l)
	if err != nil {
		return nil, err
	}

	return &RetryHelmAgent{
		form, l, helmAgent, retryCount, retryInterval,
	}, nil
}

func (a *RetryHelmAgent) ListReleases(
	ctx context.Context,
	namespace string,
	filter *types.ReleaseListFilter,
) ([]*release.Release, error) {
	for i := uint(0); i < a.retryCount; i++ {
		releases, err := a.agent.ListReleases(ctx, namespace, filter)

		if err == nil {
			return releases, nil
		} else {
			log.Printf("recreating helm agent for retrying ListReleases. Error: %v", err)

			a.agent, err = helm.GetAgentOutOfClusterConfig(ctx, a.form, a.l)

			if err != nil {
				return nil, fmt.Errorf("error recreating helm agent for retrying ListReleases: %w", err)
			}
		}

		time.Sleep(a.retryInterval)
	}

	return nil, fmt.Errorf("maxiumum number of retries (%d) reached for ListReleases", a.retryCount)
}

func (a *RetryHelmAgent) GetReleaseHistory(
	ctx context.Context,
	name string,
) ([]*release.Release, error) {
	for i := uint(0); i < a.retryCount; i++ {
		releases, err := a.agent.GetReleaseHistory(ctx, name)

		if err == nil {
			return releases, nil
		} else {
			log.Printf("recreating helm agent for retrying GetReleaseHistory. Error: %v", err)

			a.agent, err = helm.GetAgentOutOfClusterConfig(ctx, a.form, a.l)

			if err != nil {
				return nil, fmt.Errorf("error recreating helm agent for retrying GetReleaseHistory: %w", err)
			}
		}

		time.Sleep(a.retryInterval)
	}

	return nil, fmt.Errorf("maxiumum number of retries (%d) reached for GetReleaseHistory", a.retryCount)
}

func (a *RetryHelmAgent) DeleteReleaseRevision(
	ctx context.Context,
	name string,
	version int,
) error {
	for i := uint(0); i < a.retryCount; i++ {
		err := a.agent.DeleteReleaseRevision(ctx, name, version)

		if err == nil {
			return nil
		} else {
			log.Printf("recreating helm agent for retrying DeleteReleaseRevision. Error: %v", err)

			a.agent, err = helm.GetAgentOutOfClusterConfig(ctx, a.form, a.l)

			if err != nil {
				return fmt.Errorf("error recreating helm agent for retrying DeleteReleaseRevision: %w", err)
			}
		}

		time.Sleep(a.retryInterval)
	}

	return fmt.Errorf("maxiumum number of retries (%d) reached for DeleteReleaseRevision", a.retryCount)
}
