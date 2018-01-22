#!/bin/bash

# Exit on error. Append "|| true" if you expect an error.
set -o errexit
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch the error in pipeline.
set -o pipefail

MAKR_RULES_ROOT=$(dirname "${BASH_SOURCE}")/..
VERBOSE="${VERBOSE:-1}"
source "${MAKR_RULES_ROOT}/lib/init.sh"

build() {
	docker::build_images "$@"
}

push() {
	docker::push_images "$@"
}

usage() {
	log::usage_from_stdin <<EOF
usage: $cmd <commands> [TARGETS]

Available Commands:
    build      build docker image
    push       push docker image to registries
EOF
}

subcommand=${1-}
case $subcommand in
	"" | "-h" | "--help")
		usage
		;;
	*)
		shift
		${subcommand} $@
		;;
esac
