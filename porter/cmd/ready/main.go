package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/porter-dev/porter/api/server/shared/config/envloader"
	lr "github.com/porter-dev/porter/pkg/logger"
)

func main() {
	logger := lr.NewConsole(true)

	envConf, err := envloader.FromEnv()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/livez", envConf.ServerConf.Port))

	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		os.Exit(1)
	}

	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/api/readyz", envConf.ServerConf.Port))

	if err != nil || resp.StatusCode >= http.StatusBadRequest {
		os.Exit(1)
	}
}
