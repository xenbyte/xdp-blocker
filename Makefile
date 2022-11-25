CLANG ?= clang 
CFLAGS ?= -O2  -target bpf

MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(dir $(MKFILE_PATH))
HEADERS := $(CURRENT_DIR)/headers

GOCMD := go 
GOBUILD := $(GOCMD) build
GO_SOURCE := main.go
GO_BINARY := main 

EBPF_SOURCE := bpf/xdp_drop.c
EBPF_BINARY := bpf/xdp_drop.elf


all: 
	$(CLANG) -I $(HEADERS) $(CFLAGS) -c $(EBPF_SOURCE) -o $(EBPF_BINARY)

run: 
	$(GOBUILD)
