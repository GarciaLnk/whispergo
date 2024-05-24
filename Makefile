PROJECT_NAME ?= whispergo
OUTPUT := build/bin/$(PROJECT_NAME)
UPX_EXECUTABLE := upx
MODEL_FILE := ggml-tiny-q5_1.bin
MODEL_URL := https://huggingface.co/ggerganov/whisper.cpp/resolve/main/$(MODEL_FILE)

ifeq ($(WHISPER_CUDA), 1)
	CUDA_PATH ?= /usr/local/cuda
	CUDA_LIBPATH ?= $(CUDA_PATH)/lib64
	CGO_LDFLAGS += -lcublas -lcudart -lcuda -L$(CUDA_LIBPATH) -L$(CUDA_LIBPATH)/stubs
endif

CC ?= gcc

INCLUDE_PATH := $(abspath external/whisper.cpp)
LIBRARY_PATH := $(abspath external/whisper.cpp)


.PHONY: whisper release compressed dev clean

whisper:
	echo Build whisper
	@${MAKE} CC=$(CC) -C external/whisper.cpp libwhisper.a

models/$(MODEL_FILE):
	echo Download tiny model
	@wget -q -nc $(MODEL_URL) -P models

release: whisper models/$(MODEL_FILE)
	CGO_LDFLAGS="$(CGO_LDFLAGS)" C_INCLUDE_PATH=$(INCLUDE_PATH) LIBRARY_PATH=$(LIBRARY_PATH) wails build $(BUILD_ARGS) -o $(PROJECT_NAME)

dev: whisper models/$(MODEL_FILE)
	C_INCLUDE_PATH=$(INCLUDE_PATH) LIBRARY_PATH=$(LIBRARY_PATH) wails dev

clean:
	rm -rf build/bin
	@${MAKE} -C external/whisper.cpp clean

.DEFAULT_GOAL := release