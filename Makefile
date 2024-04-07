PROJECT_NAME := whispergo
OUTPUT := build/bin/$(PROJECT_NAME)
UPX_EXECUTABLE := upx

CGO_LDFLAGS=-L$(CUDA_PATH)/stubs -lcuda

CC ?= gcc

INCLUDE_PATH := $(abspath external/whisper.cpp)
LIBRARY_PATH := $(abspath external/whisper.cpp)


.PHONY: whisper release compressed dev clean

whisper:
	@echo Build whisper
	@${MAKE} -C external/whisper.cpp libwhisper.a

release: whisper
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} wails build
	@echo $(OUTPUT)

compressed: release
	$(UPX_EXECUTABLE) $(UPXFLAGS) $(OUTPUT);

dev: whisper
	@C_INCLUDE_PATH=${INCLUDE_PATH} LIBRARY_PATH=${LIBRARY_PATH} wails dev

clean:
	rm -rf build/bin
	@${MAKE} -C external/whisper.cpp clean

.DEFAULT_GOAL := release