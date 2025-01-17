package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/grafana/grafana/pkg/build/config"
	"github.com/grafana/grafana/pkg/build/errutil"
	"github.com/grafana/grafana/pkg/build/plugins"
	"github.com/grafana/grafana/pkg/build/syncutil"
	"github.com/urfave/cli/v2"
)

func BuildInternalPlugins(c *cli.Context) error {
	cfg := config.Config{
		NumWorkers: c.Int("jobs"),
	}

	const grafanaDir = "."
	metadata, err := config.GetMetadata(filepath.Join("dist", "version.json"))
	if err != nil {
		return err
	}
	verMode, err := config.GetVersion(metadata.ReleaseMode)
	if err != nil {
		return err
	}

	log.Println("Building internal Grafana plug-ins...")

	ctx := context.Background()

	p := syncutil.NewWorkerPool(cfg.NumWorkers)
	defer p.Close()

	var g *errutil.Group
	g, ctx = errutil.GroupWithContext(ctx)
	if err := plugins.Build(ctx, grafanaDir, p, g, verMode); err != nil {
		return cli.Exit(err.Error(), 1)
	}
	if err := g.Wait(); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	if err := plugins.Download(ctx, grafanaDir, p); err != nil {
		return cli.Exit(err.Error(), 1)
	}

	log.Println("Successfully built Grafana plug-ins!")

	return nil
}
