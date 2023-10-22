# Makefile
# Define the target platforms
PLATFORMS := windows linux darwin
# Define the output directory and names for the executables
OUTPUT_DIR := builds
WINDOWS_BINARY := Neptune.exe
LINUX_BINARY := Neptune

FYNE_EXISTS := $(shell command -v fyne 2> /dev/null)
BUILD_COMMAND := fyne package --icon icon.png --appBuild 2 --appVersion 1.0.1 --release

OPTIMIZE_GCC := CGO_CFLAGS="-flto -O2" CGO_LDFLAGS="-flto"
CLEAN_COMMAND := rm -rf $(OUTPUT_DIR)
SYSTRAY := -tags systray

LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

Name := "Neptune"
Exec := "Neptune"
Icon := "Neptune.png"

.PHONY: all
all: $(PLATFORMS) move

.PHONY: windows
windows: BUILD_COMMAND := fyne package --icon icon.png --release
windows: GOOS := windows
windows: CC ?= x86_64-w64-mingw32-gcc
windows: CXX ?= x86_64-w64-mingw32-g++
windows: GOARCH ?= amd64
windows:
	@echo "Building for windows"
	$(if $(and $(PKG), $(FYNE_EXISTS)), \
		CGO_ENABLED=1 GOOS=$(GOOS) CC=$(CC) CXX=$(CXX) $(BUILD_COMMAND) --name $(WINDOWS_BINARY) -os $(GOOS), \
		CGO_ENABLED=1 GOOS=$(GOOS) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) go build -ldflags="-s -w -H=windowsgui" $(if $(TAGS),,$(SYSTRAY)) \
	)

.PHONY: linux
linux: GOOS := linux
linux: CC ?= gcc
linux: CXX ?= g++
linux: GOARCH ?= amd64
linux:
	@echo "Building for linux"
	$(if $(and $(PKG), $(FYNE_EXISTS)), \
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) $(BUILD_COMMAND) --name $(LINUX_BINARY) -os $(GOOS), \
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) go build -ldflags="-s -w" $(if $(TAGS),,$(SYSTRAY)) -o misc/usr/local/bin/ \
	)

.PHONY: linux-cli
linux-cli: GOOS := linux
linux-cli: CC ?= gcc
linux-cli: CXX ?= g++
linux-cli: GOARCH ?= amd64
linux-cli:
	@echo "Building CLI for linux"
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) go build -ldflags="-s -w" $(if $(TAGS),,$(SYSTRAY)) -o misc/usr/local/bin/$(LINUX_BINARY) cmd/Neptune-Cli/main.go

.PHONY: darwin
darwin: GOOS := darwin
darwin: CC ?= gcc
darwin: CXX ?= g++
darwin: GOARCH ?= amd64
darwin:
	@echo "Building for macos"
	$(if $(and $(PKG), $(FYNE_EXISTS)), \
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) $(BUILD_COMMAND) --name $(LINUX_BINARY) -os $(GOOS), \
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) go build -ldflags="-s -w" $(if $(TAGS),,$(SYSTRAY)) \
	)

.PHONY: darwin-cli
darwin-cli: GOOS := darwin
darwin: CC ?= gcc
darwin: CXX ?= g++
darwin: GOARCH ?= amd64
darwin-cli:
	@echo "Building CLI for macos"
		CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) CC=$(CC) CXX=$(CXX) $(OPTIMIZE_GCC) go build -ldflags="-s -w" $(if $(TAGS),,$(SYSTRAY)) -o $(LINUX_BINARY) cmd/Neptune-Cli/main.go

.PHONY: move
move:
	mv Neptune.tar.xz Neptune.exe builds/

.PHONY: update-pkg-cache
update-pkg-cache:
	GOPROXY=https://proxy.golang.org GO111MODULE=on \
	go get github.com/m1ndo/$(PACKAGE)@v$(VERSION)

.PHONY: install
install:
	install -Dm00644 misc/usr/local/share/applications/$(Name).desktop $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	install -Dm00755 misc/usr/local/bin/$(Exec) $(DESTDIR)$(PREFIX)/bin/$(Exec)
	install -Dm00644 misc/usr/local/share/pixmaps/$(Icon) $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)

.PHONY: uninstall
uninstall:
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(Name).desktop
	-rm $(DESTDIR)$(PREFIX)/bin/$(Exec)
	-rm $(DESTDIR)$(PREFIX)/share/pixmaps/$(Icon)

.PHONY: user-install
user-install:
	install -Dm00644 misc/usr/local/share/applications/$(Name).desktop $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop
	install -Dm00755 misc/usr/local/bin/$(Exec) $(DESTDIR)$(HOME)/.local/bin/$(Exec)
	install -Dm00644 misc/usr/local/share/pixmaps/$(Icon) $(DESTDIR)$(HOME)/.local/share/icons/$(Icon)
	sed -i -e "s,Exec=$(Exec),Exec=$(DESTDIR)$(HOME)/.local/bin/$(Exec),g" $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop

.PHONY: user-uninstall
user-uninstall:
	-rm $(DESTDIR)$(HOME)/.local/share/applications/$(Name).desktop
	-rm $(DESTDIR)$(HOME)/.local/bin/$(Exec)
	-rm $(DESTDIR)$(HOME)/.local/share/icons/$(Icon)

.PHONY: clean
clean:
	@echo "Cleaning..."
	$(CLEAN_COMMAND)
