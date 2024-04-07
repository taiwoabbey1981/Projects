package pack

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/client-go/util/homedir"
)

type BuildpackNameTestResult struct {
	name string
	err  error
}

func setupPackTest(tb testing.TB) func(tb testing.TB) {
	porterHome := filepath.Join(homedir.HomeDir(), ".porter")
	if err := os.MkdirAll(porterHome, 0o750); err != nil {
		tb.Errorf("unable to initialize porter home folder for tests: %s", err.Error())
	}

	return func(tb testing.TB) {
		dstFiles, err := os.ReadDir(porterHome)
		if err != nil {
			return
		}

		for _, info := range dstFiles {
			if !info.Type().IsDir() {
				continue
			}

			if err := os.RemoveAll(filepath.Join(porterHome, info.Name())); err != nil {
				tb.Errorf("unable to delete porter home subfolder for tests: %s", err.Error())
			}
		}
	}
}

func TestGetBuildpackName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected BuildpackNameTestResult
	}{
		{
			"empty buildpack name",
			"",
			BuildpackNameTestResult{"", errors.New("please specify a buildpack name")},
		},
		{
			"cnb short name",
			"heroku/nodejs",
			BuildpackNameTestResult{"heroku/nodejs", nil},
		},
		{
			"cnb urn",
			"urn:cnb:registry:heroku/nodejs",
			BuildpackNameTestResult{"urn:cnb:registry:heroku/nodejs", nil},
		},
		{
			"cnb shim",
			"https://cnb-shim.herokuapp.com/v1/heroku/nodejs?version=0.0.0&name=Node.js",
			BuildpackNameTestResult{"https://cnb-shim.herokuapp.com/v1/heroku/nodejs?version=0.0.0&name=Node.js", nil},
		},
		{
			"invalid tgz",
			"https://example.com/tar.tgz",
			BuildpackNameTestResult{"", errors.New("please provide either a github.com URL or a ZIP file URL")},
		},
		{
			"github repo",
			"https://github.com/heroku/buildpacks-nodejs/archive/fa2dc153e4683181608307ecb3922eaaeb43d92c.zip",
			BuildpackNameTestResult{filepath.Join(homedir.HomeDir(), ".porter", "buildpacks-nodejs-fa2dc153e4683181608307ecb3922eaaeb43d92c"), nil},
		},
		{
			"github repo zip",
			"https://github.com/heroku/buildpacks-nodejs/archive/refs/tags/v1.1.6.zip",
			BuildpackNameTestResult{filepath.Join(homedir.HomeDir(), ".porter", "buildpacks-nodejs-1.1.6"), nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardownPackTest := setupPackTest(t)
			defer teardownPackTest(t)

			ctx := context.Background()
			actual, err := getBuildpackName(ctx, tt.input)
			if actual != tt.expected.name {
				t.Errorf("got %s, want %s", actual, tt.expected.name)
			}

			if err != nil && tt.expected.err == nil {
				t.Errorf("got unexpected error: %s", err.Error())
			}

			if err == nil && tt.expected.err != nil {
				t.Errorf("missing expected error %s", tt.expected.err)
			}

			if err != nil && tt.expected.err != nil {
				if err.Error() != tt.expected.err.Error() {
					t.Errorf("wrong error: %v", err)
				}
			}
		})
	}
}
