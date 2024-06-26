GO ?= go
GIT ?= git
GOFMT ?= gofmt "-s"
XSLT ?= xsltproc
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
GOTOOLS=tools.go
MAIN_BINARY=xsd2goxsl
MAIN_PACKAGE=.
COVERAGE=coverage.out
XSD2GOXSL=xsd2go.xsl
XSD_TEST_DIR=./xsd

.PHONY: check run fmt fmt-check vet clean

all: $(MAIN_BINARY)

$(MAIN_BINARY):  $(GOFILES)
	$(GO) build -v -o $(MAIN_BINARY) ./$(MAIN_PACKAGE)

check: generate
	$(GO) test -v -cover -coverpkg=./... -coverprofile=$(COVERAGE) ./...

run: all
	$(GO) run -buildvcs=true ./$(MAIN_PACKAGE)

fmt:
	$(GOFMT) -w $(GOFILES)

fmt-check:
	$(GOFMT) -d $(GOFILES)

vet:
	$(GO) vet $(PACKAGES)

clean:
	rm -f $(MAIN_BINARY)
	rm -f $(COVERAGE)
	rm -rf $(XSD_TEST_DIR)/*_test/

check-xsd: $(XSD2GOXSL) $(XSD_TEST_DIR)/*.xsd
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), rm -rf $(file)_test;)
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), mkdir $(file)_test;)
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), $(XSLT) --stringparam qAttrImport 'utils "github.com/indexdata/go-utils/utils"' \
	--stringparam qAttrType "utils.PrefixAttr" --stringparam buildtag checkxsd $(XSD2GOXSL) $(file) > $(file)_test/schema.go;)
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), diff $(file)_test/schema.go $(file).out || exit;)
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), $(GO) vet -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(wildcard $(XSD_TEST_DIR)/*.xsd), $(GO) build -tags checkxsd $(file)_test/schema.go || exit;)
