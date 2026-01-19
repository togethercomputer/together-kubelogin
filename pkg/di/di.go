//go:build wireinject
// +build wireinject

// Package di provides dependency injection.
package di

import (
	"github.com/google/wire"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/cmd"
	credentialpluginreader "github.com/pahluwalia-tcloud/together-kubelogin/pkg/credentialplugin/reader"
	credentialpluginwriter "github.com/pahluwalia-tcloud/together-kubelogin/pkg/credentialplugin/writer"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/browser"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/clock"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/logger"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/reader"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/infrastructure/stdio"
	kubeconfigLoader "github.com/pahluwalia-tcloud/together-kubelogin/pkg/kubeconfig/loader"
	kubeconfigWriter "github.com/pahluwalia-tcloud/together-kubelogin/pkg/kubeconfig/writer"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/oidc/client"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/tlsclientconfig/loader"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/tokencache/repository"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/usecases/authentication"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/usecases/clean"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/usecases/credentialplugin"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/usecases/setup"
	"github.com/pahluwalia-tcloud/together-kubelogin/pkg/usecases/standalone"
)

// NewCmd returns an instance of infrastructure.Cmd.
func NewCmd() cmd.Interface {
	wire.Build(
		NewCmdForHeadless,

		// dependencies for production
		clock.Set,
		stdio.Set,
		logger.Set,
		browser.Set,
	)
	return nil
}

// NewCmdForHeadless returns an instance of infrastructure.Cmd for headless testing.
func NewCmdForHeadless(clock.Interface, stdio.Stdin, stdio.Stdout, logger.Interface, browser.Interface) cmd.Interface {
	wire.Build(
		// use-cases
		authentication.Set,
		standalone.Set,
		credentialplugin.Set,
		setup.Set,
		clean.Set,

		// infrastructure
		cmd.Set,
		reader.Set,
		kubeconfigLoader.Set,
		kubeconfigWriter.Set,
		repository.Set,
		client.Set,
		loader.Set,
		credentialpluginreader.Set,
		credentialpluginwriter.Set,
	)
	return nil
}
