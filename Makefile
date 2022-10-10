# ====================================================================================
# Setup Project

PROVIDER_VERSION ?= v0.31.0

# set the shell to bash always
SHELL := /bin/bash

# ====================================================================================
# Setup directories and paths

# a working directory that holds all temporary or working items generated
# during the build. The items will be discarded on a clean build and they
# will never be cached.
ifeq ($(origin WORK_DIR), undefined)
WORK_DIR := $(ROOT_DIR).work
endif

# ====================================================================================
# Colors

BLACK        := $(shell printf "\033[30m")
BLACK_BOLD   := $(shell printf "\033[30;1m")
RED          := $(shell printf "\033[31m")
RED_BOLD     := $(shell printf "\033[31;1m")
GREEN        := $(shell printf "\033[32m")
GREEN_BOLD   := $(shell printf "\033[32;1m")
YELLOW       := $(shell printf "\033[33m")
YELLOW_BOLD  := $(shell printf "\033[33;1m")
BLUE         := $(shell printf "\033[34m")
BLUE_BOLD    := $(shell printf "\033[34;1m")
MAGENTA      := $(shell printf "\033[35m")
MAGENTA_BOLD := $(shell printf "\033[35;1m")
CYAN         := $(shell printf "\033[36m")
CYAN_BOLD    := $(shell printf "\033[36;1m")
WHITE        := $(shell printf "\033[37m")
WHITE_BOLD   := $(shell printf "\033[37;1m")
CNone        := $(shell printf "\033[0m")

# ====================================================================================
# Logger

TIME_LONG  = `date +%Y-%m-%d' '%H:%M:%S`
TIME_SHORT = `date +%H:%M:%S`
TIME       = $(TIME_SHORT)

INFO = echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN = echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR  = echo ${TIME} ${RED}[FAIL]${CNone}
OK   = echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL = (echo ${TIME} ${RED}[FAIL]${CNone} && false)

# ====================================================================================
# Commands

all: fetch generate

fetch:
	@$(INFO) Fetch crossplane provider-aws GitRepo
	@mkdir -p ${WORK_DIR}
	@if [ ! -d "${WORK_DIR}/provider-aws" ]; then \
		cd ${WORK_DIR} && git clone "https://github.com/crossplane/provider-aws.git"; \
	fi
	@cd ${WORK_DIR}/provider-aws && git fetch origin && git checkout $(PROVIDER_VERSION)
	@$(OK) Fetch crossplane provider-aws GitRepo

generate:
	@$(INFO) Generating CRDs
	@go run pkg/main.go .
	@$(OK) Generating CRDs
