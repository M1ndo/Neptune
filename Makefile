# Makefile
# Define the target platforms
PLATFORMS := windows linux
# Define the output directory and names for the executables
OUTPUT_DIR := builds
WINDOWS_BINARY := Neptune.exe
LINUX_BINARY := Neptune

BUILD_COMMAND := fyne package --icon icon.png --appBuild 2 --appVersion 1.0.1 --release
OPTIMIZE_GCC := CGO_CFLAGS="-flto -O2" CGO_LDFLAGS="-flto"
CLEAN_COMMAND := rm -rf $(OUTPUT_DIR)

.PHONY: all clean $(PLATFORMS)

all: $(PLATFORMS) move

windows:
	@echo "Building for windows"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(BUILD_COMMAND) --name $(WINDOWS_BINARY) -os windows

linux:
	@echo "Building for linux"
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ $(OPTIMIZE_GCC) $(BUILD_COMMAND) --name $(LINUX_BINARY) -os linux

move: 
	mv Neptune.tar.xz Neptune.exe builds/

update-pkg-cache:
	GOPROXY=https://proxy.golang.org GO111MODULE=on \
	go get github.com/m1ndo/$(PACKAGE)@v$(VERSION)

clean:
	@echo "Cleaning..."
	$(CLEAN_COMMAND)
