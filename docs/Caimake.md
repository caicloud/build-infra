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
go get -u github.com/caicloud/build-infra/...
```

### Requirements

-   `make` v4.2.1 or higher
-   `sed` GNU sed, v4.4 or higher
-   `git` v2.10.1 or higher
-   `docker` for building and pushing container
-   `golang` for compiling golang project, v1.8.0 or higher

### Initialization

Getting started from a git repo, use `caimake init <language>` to initialize it. 
The following languages are supported now:
- go

`caimake init` will scan your project:
-   treat project's name as docker images prefix

for `golang`, it will:
-   find all subdirectories under `cmd` as targets of `make build` 
-   find all subdirectories under `build` as targets of `make container` 

Then, caimake generates a `Makefile` and put `.caimake` under the project root.

Now, you are able to use `make` instead of `caimake`.

Use `caimake [command] —help` for more information about teh command.

See [Makefile Specification](Specification.md) for more information about detail.

### Expansibility

For expansibility, user can define their own custom targets and override ENV variables in `Makefile.expansion`. The original `Makefile` will include the expansion automatically. 

User SHOULD put all custom changes in `Makefile.expansion` to override the options in `Makefile`.

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

