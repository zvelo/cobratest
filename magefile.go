// +build mage

package main

import (
	"context"
	"fmt"
	"go/build"

	"github.com/magefile/mage/mg"

	"zvelo.io/zmage"
)

const (
	Exe         = "./cobratest"
	ExeDir      = "./cmd/cobratest"
)

// Default is the default mage target
var Default = Build

// Build all executables
func Build(ctx context.Context) error {
	mg.CtxDeps(ctx, Cobratest)
	return nil
}

// Dep ensures `Gopkg.lock` and `vendor/` in sync with `Gopkg.toml`
func Dep(ctx context.Context) error {
	return zmage.Dep(ctx)
}

// Cobratest builds the `cobratest` binary
func Cobratest(ctx context.Context) error {
	mg.CtxDeps(ctx, Dep)
	return zmage.BuildExe(build.Default, ExeDir, Exe)
}

// Install installs all the executables
func Install(ctx context.Context) error {
	mg.CtxDeps(ctx, Dep)
	return zmage.Install(build.Default, ExeDir)
}

// GoVet runs `go vet`
func GoVet(ctx context.Context) error {
	return zmage.GoVet()
}

// Clean up after yourself
func Clean(ctx context.Context) error {
	return zmage.Clean(Exe)
}

// GoPackages all the non-vendor packages in the repository
func GoPackages(ctx context.Context) error {
	pkgs, err := zmage.GoPackages(build.Default)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg.ImportPath, pkg.Name)
	}

	return nil
}
