# Makefile Specification

#### Version 0.1.0

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [BCP 14](https://tools.ietf.org/html/bcp14) [RFC2119](https://tools.ietf.org/html/rfc2119) [RFC8174](https://tools.ietf.org/html/rfc8174) when, and only when, they appear in all capitals, as shown here.

## Introduction

The Makefile Specification defines standard interfaces which MUST or SHOULD be implemented in a cloud native application's Makefile.

## Specification

### Overview

The Options in following sections are environment variables used to change the default behavior of Makefile.

The severity `REQUIRED` means that you MUST support the option. 

The severity `OPTIONAL` means that you MAY support the option.

### Usage for targets

The Makefile MUST print usage message for make targets when `$HELP == y`

#### Example

```bash
HELP=y make build
# Build code.
# make all == make build
#
# Args:
#   WHAT: Directory names to build.  If any of these directories has a 'main'
#     package, the build will produce executable files under bin/.
#     If not specified, "everything" will be built.
#   GOFLAGS: Extra flags to pass to 'go' when building.
#   GOLDFLAGS: Extra linking flags passed to 'go' when building.
#   GOGCFLAGS: Additional go compile flags passed to 'go' when building.
#
# Example:
#   make
#   make all or make build
#   make build WHAT=cmd/server GOFLAGS=-v
#   make all GOGCFLAGS="-N -l"
#     Note: Use the -N -l options to disable compiler optimizations an inlining.
#           Using these build options allows you to subsequently use source
#           debugging tools like delve.
```

### Color log

The Makefile MAY print colorized log when `$COLOR_LOG == true`. 

It is OPTIONAL. 

### Version

The valid version MUST be able to describe the current status of git tree.

A complete version information contains:

-   a git tag versioned by [Semantic Versioning 2.0.0](http://semver.org/spec/v2.0.0.html) (semver). e.g. v1.0.0 v0.1.0-alpha.1
-   commits number after latest git tag and latest commit hash
-   git tree status: 
    -   `clean` indicates no changes since the git commit id.
    -   `dirty` indicates source code changes after the git commit id
    -   `archive` indicates the tree was produced by `git archive`

#### Use Case

| version                 | docker tag              | description                              |
| ----------------------- | ----------------------- | ---------------------------------------- |
| v0.0.1                  | v0.0.1                  | a clean git tree with the latest tag `v0.0.1` |
| v0.0.1-dirty            | v0.0.1-dirty            | a dirty git tree with changes not staged |
| v0.0.1-2+1b4531e1acf800 | v0.0.1-2-1b4531e1acf800 | a clean git tree with a tag ` v0.0.1`.<br />There are two new commits after the tag and the latest commit is `1b4531e1acf800` |

If your git repo contains changes not staged, a `-dirty` suffix MUST always be appended to final version.

It is important that the second and third cases MUST only be used for debugging and testing to alert developer that you are working on a dirty git tree and you need to create a new git tag before release.

**For releasing a new version**:

1.  the final version MUST follow the Semantic Versioning 2.0.0 (semver). e.g. v1.0.0 v0.1.0-alpha.1
2.  the git tree MUST be clean
3.  the latest tag MUST be even with the latest commit id

#### Options

| Name    | Usage                       | Type   | Default | Severity |
| ------- | --------------------------- | ------ | ------- | -------- |
| VERSION | used to overwrite `git tag` | string |         | REQUIRED |

#### Example

```bash
VERSION=v0.0.1 make all
```

### Compile

The Makefile MUST support the following targets for compiling.

```bash
make all # equal to make build
make build
make build-local
make build-in-container
```

Rules for Makefile targets:

-   `make all` MUST be equal to `make build`
-   `make build` means to compile you project. The Makefile MUST support the ability to compile the project on localhost and in docker container.
-   For convenience, the Makefile SHOULD support the ability to build targets directories directly. For example:

```bash
make cmd/controller # equal to make build WHAT=cmd/controller
```

#### Options

| Name        | Usage                                    | Type   | Default | Severity |
| ----------- | ---------------------------------------- | ------ | ------- | -------- |
| LOCAL_BUILD | If set true, project will be built on local machine, otherwise, built in docker. | string | true    | REQUIRED |
| WHAT        | Directory names to build.  If any of these directories has a main package, the build will produce executable files under bin/. e.g `cmd/controller` | string |         | REQUIRED |

#### Example

The most basic rules:

```bash
# make all targets
make all
make build
make build LOCAL_BUILD=true

# make cmd/controller on localhost
make build WHAT=cmd/controller LOCAL-BUILD=true
# make cmd/controller in docker
make build WHAT=cmd/controller LOCAL-BUILD=false
```

the followings are extensions for the basic way

```bash
# make all targets on localhost
make build-local
# make all targets in docker contaienr
make build-in-container
# make cmd/controller
make cmd/controller

# more ...
make build-local WHAT=cmd/controller
make build-in-container WHAT=cmd/controller
make cmd/controller LOCAL_BUILD=false
```

#### Compile different language

For more information about different language compile specification:

-   [Golang Compile](#golangCompile)


### Unit Test

The Makefile MUST support the following targets for unit test

```bash
make unittest
```

#### Test different language 

For more information about different language unittest specification:

-   [Golang Unittest](#golangUnittest)

### Docker

The Makefile MUST support the following targets for docker build.

```bash
make container
make push
```

Rules for Makefile targets:

-   The Makefile MUST support the ability to build or push docker images for multiple registries, and modify docker image by adding prefix and suffix.
-   For convenience, the Makefile SHOULD support the ability to build targets directories which contain Dockerfiles directly. For example:

```bash
make build/controller # equal to make container WHAT=build/controller
```

#### Options

| Name                 | Usage                                    | Type   | Default | Severity |
| -------------------- | ---------------------------------------- | ------ | ------- | -------- |
| DOCKER_REGISTRIES    | Docker registries to push                | array  |         | REQUIRED |
| DOCKER_IMAGE_PREFIX  | Docker image prefix.                     | string |         | REQUIRED |
| DOCKER_IMAGE_SUFFIX  | Docker image suffix.                     | string |         | REQUIRED |
| DOCKER_FORCE_PUSH    | Force pushing to override images in remote registries | string | false  | OPTIONAL |
| DOCKER_BUILD_TARGETS | All pre-defined directory names of targets for docker build. e.g `build/controller` | array  |         | REQUIRED |
| WHAT                 | Directories containing Dockerfile        | string |         | REQUIRED |

The difference between `DOCKER_BUILD_TARGETS` and `WHAT` is that `DOCKER_BUILD_TARGETS` means the all pre-defined docker build targets and `WHAT` means the target in this build. Makefile will build all targets in `DOCKER_BUILD_TARGETS` if you don't specify any targets in `WHAT`.

#### Example

```bash
# build all docker images
make container
# build build/controller's Dockerfile
make build/controller
make container WHAT=build/controller
# push all docker images
make push

# build and push build/controller docker iamge
make container push WHAT=build/controller
```

## Language Specification
### Golang

#### <a name="golangCompile"></a>Compile

For compiling golang, the Makefile MUST support the ability to build multiple main packages for multiple platforms in one project. And It MUST has the ability to enable or disable the use of cgo.

##### Options

| Name                | Usage                                    | Type   | Default                  | Severity |
| ------------------- | ---------------------------------------- | ------ | ------------------------ | -------- |
| GO_ONBUILD_IMAGE    | Porject will be built in the image when `${LOCAL_BUILD} != true` | string | golang:1.9.2-alpine3.6   | REQUIRED |
| GO_BUILD_PLATFORMS  | The project will be built for these platforms. | array  | linux/amd64 darwin/amd64 | REQUIRED |
| GOFLAGS             | Extra flags to pass to 'go' when building. | string |                          | REQUIRED |
| GOLDFLAGS           | Extra linking flags passed to 'go' when building. | string |                          | REQUIRED |
| GOGCFLAGS           | Additional go compile flags passed to 'go' when building. | string |                          | REQUIRED |
| GO_BUILD_TARGETS    | All pre-defined directory names of targets for go build. e.g `cmd/controller` | array  | user defined             | REQUIRED |
| GO_STATIC_LIBRARIES | Determine which go build targets should use `CGO_ENABLED=0`. | array  | user defined             | REQUIRED |

##### Example

```bash
# make all targets for linux/amd64 
make all GO_BUILD_PLATFORMS=linux/amd64

# make cmd/controller for linux/amd64
make all WHAT=cmd/controller GO_BUILD_PLATFORMS=linux/amd64
make cmd/controller GO_BUILD_PLATFORMS=linux/amd64

# make cmd/controller for linux/amd64 in docker container
make build-in-container WHAT=cmd/controller GO_BUILD_PLATFORMS=linux/amd64
make build WHAT=cmd/controller LOCAL_BUILD=false  GO_BUILD_PLATFORMS=linux/amd64
make cmd/controller LOCAL_BUILD=false GO_BUILD_PLATFORMS=linux/amd64

# make cmd/controller for linux/amd64 on local
make build-local WHAT=cmd/controller GO_BUILD_PLATFORMS=linux/amd64
make build LOCAL_BUILD=true WHAT=cmd/controller GO_BUILD_PLATFORMS=linux/amd64
make cmd/controller LOCAL_BUILD=true GO_BUILD_PLATFORMS=linux/amd64
```

#### <a name="golangUnittest"></a>Unittest

The Makefile MUST support the ability to skip some directories which are unexpect for unittest. By the way, the packages under `vendor test tests scripts hack` MUST be always ignored in unittest.

##### Options

| Name               | Usage                                    | Type  | Default | Severity |
| ------------------ | ---------------------------------------- | ----- | ------- | -------- |
| GO_TEST_EXCEPTIONS | Go test will ignore the pkg under exceptions dirs.<br />` vendor test tests scripts hack` are always be skipped. | array |         | REQUIRED |

##### Example

```bash
# skip third_party dir
make unittest GO_TEST_EXCEPTIONS=third_party
```

