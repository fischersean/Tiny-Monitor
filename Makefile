.PHONY: build clean run dist bundle

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

bundle: 
	@make build
	rm -rf dist/
	mkdir -p $(EXEPATH)
	mkdir -p $(RESOURCEPATH)
	mv $(DISTNAME) $(EXEPATH)/$(DISTNAME)
	cp Info.plist $(DISTPATH)
	cp resources/icon.icns $(RESOURCEPATH)

dist:
	@make bundle
	codesign  --deep -s seanwfischer@gmail.com dist/$(BUNDLENAME)
	codesign --verify --verbose dist/$(BUNDLENAME)

run:
	./$(DISTNAME)

clean:
	rm -f $(DISTNAME)
	rm -rf dist/
