package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"time"
)

const (
	namespace               = "test-ns"
	releaseName             = "test-release"
	helmInstallationTimeout = time.Second * 10
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(pprof.New())

	app.Post("/charts/install", func(ctx *fiber.Ctx) error {
		r, err := installChart()
		if err != nil {
			return err
		}

		return ctx.JSON(r)
	})

	panic(app.Listen(":8080"))
}

func installChart() (*release.Release, error) {
	actionConfig := new(action.Configuration)

	opts := genericclioptions.ConfigFlags{}

	err := actionConfig.Init(
		&opts,
		namespace,
		"",
		func(format string, v ...interface{}) {
			log.Printf(format, v)
		},
	)
	if err != nil {
		return nil, err
	}

	installAction := action.NewInstall(actionConfig)
	installAction.Namespace = namespace
	installAction.ReleaseName = releaseName
	installAction.Atomic = true
	installAction.Timeout = helmInstallationTimeout
	installAction.CreateNamespace = true

	loadedChart, err := loader.Load("busybox")
	if err != nil {
		return nil, err
	}

	return installAction.Run(loadedChart, nil)
}
