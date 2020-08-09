.PHONY: build clean run dist bundle test

GO = go
BUILDFLAGS = GO11MODULES=1 GOOS=darwin
MODNAME = tinymonitor
DISTNAME = Tiny\ Monitor
BUNDLENAME = $(DISTNAME).app
DISTPATH = dist/$(BUNDLENAME)/Contents
EXEPATH = $(DISTPATH)/MacOS
RESOURCEPATH = $(DISTPATH)/Resources

build:
	$(BUILDFLAGS) $(GO) build -o $(DISTNAME) cmd/$(MODNAME)/main.go 

test:
	$(GO) test ./...

bundle: 
	@make build
	@make test 
	rm -rf dist/
	mkdir -p $(EXEPATH)
	mkdir -p $(RESOURCEPATH)
	mv $(DISTNAME) $(EXEPATH)/$(DISTNAME)
	cp Info.plist $(DISTPATH)
	cp resources/icon.icns $(RESOURCEPATH)

dist:
	@make bundle
	# Make sure the identity used is the "Developer ID Application" cert
	codesign  --deep -s Sean\ Fischer dist/$(BUNDLENAME)
	codesign --verify --verbose dist/$(BUNDLENAME)
	create-dmg dist/$(BUNDLENAME) dist/

run:
	./$(DISTNAME)

clean:
	rm -f $(DISTNAME)
	rm -rf dist/
