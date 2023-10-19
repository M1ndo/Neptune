# Makefile
# Define the target platforms
PLATFORMS := windows linux
# Define the output directory and names for the executables
OUTPUT_DIR := builds
WINDOWS_BINARY := Neptune.exe
LINUX_BINARY := Neptune

FYNE_EXISTS := $(shell command -v fyne2 2> /dev/null)
BUILD_COMMAND := fyne package --icon icon.png --appBuild 2 --appVersion 1.0.1 --release

OPTIMIZE_GCC := CGO_CFLAGS="-flto -O2" CGO_LDFLAGS="-flto"
CLEAN_COMMAND := rm -rf $(OUTPUT_DIR)

LOCAL != test -d $(DESTDIR)/usr/local && echo -n "/local" || echo -n ""
LOCAL ?= $(shell test -d $(DESTDIR)/usr/local && echo "/local" || echo "")
PREFIX ?= /usr$(LOCAL)

Name := "Neptune"
Exec := "Neptune"
Icon := "Neptune.png"

# .PHONY: all clean $(PLATFORMS)

.PHONY: all
all: $(PLATFORMS) move

.PHONY: windows
windows:
	@echo "Building for windows"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ $(BUILD_COMMAND) --name $(WINDOWS_BINARY) -os windows

.PHONY: linux
linux:
	@echo "Building for linux"
	$(if $(FYNE_EXISTS), \
		CGO_ENABLED=1 GOOS=linux $(OPTIMIZE_GCC) $(BUILD_COMMAND) --name $(LINUX_BINARY) -os linux, \
		CGO_ENABLED=1 GOOS=linux $(OPTIMIZE_GCC) go build -o misc/usr/local/bin/ . \
	)

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
