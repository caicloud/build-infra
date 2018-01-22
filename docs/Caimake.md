# Caimake

## About the Porject

caimake helps you to build your project appropriately.

It follows the [Makefile Specification](Specification.md).

## Getting Started

### Installation

#### 1. Download from Github release

Download directly from GIthub [release](directly)

#### 2. Use `go get`

```shell
go get -u https://github.com/caicloud/build-infra/...
```

### Requirements

-   `make`
-   `sed`
-   `git`
-   `docker` for build and push container
-   `golang` for compile golang project

### Usage

Getting started from a git repo, use `caimake init` to initialize it. 

```shell
caimake init
```

`caimake init` will scan your project:

-   find all subdirectories under `cmd` as targets of `make build` 
-   find all subdirectories under `build` as targets of `make docker` 
-   treat project's name as docker images prefix

Then, it generates a `Makefile` and `.caimake` under the project root.

Now, you are able to use `make` instead of `caimake`.

Use `caimake [command] —help` for more information about a command.

See [Makefile Specification](Specification.md) for more information about detail.

### Expansibility

For expansibility, user can define their own custom targets and override ENV variables in `Makefile.expansion`. The original `Makefile` will include the expansion automatically. 

#### pre-build hook

A `pre-build` target will be triggered before `make build` if it is defined.

### Update

#### Self update

Use `caimake update` to update `caimake` binary, it will download latest release from Github releases

#### Manually

```shell
git clone https://github.com/caicloud/build-infra/ 
cd build-infra
make build·
```

