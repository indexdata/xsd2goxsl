GO ?= go
GIT ?= git
GOFMT ?= gofmt "-s"
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
GOTOOLS=tools.go
MAIN_BINARY=xsd2goxsl
MAIN_PACKAGE=.
COVERAGE=coverage.out
XSD2GOXSL=xsd2go.xsl
XSD_TEST_DIR=./xsd
XSD_FILES := $(wildcard $(XSD_TEST_DIR)/*.xsd)

.PHONY: check run fmt fmt-check vet clean

all: $(MAIN_BINARY)

$(MAIN_BINARY):  $(GOFILES)
	$(GO) build -v -o $(MAIN_BINARY) ./$(MAIN_PACKAGE)

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

define nsImports
namespaceImports=http://docs.oasis-open.org/ns/search-ws/diagnostic=github.com/indexdata/xsd2goxsl/xsd/diagnostic.xsd_test,\
http://docs.oasis-open.org/ns/search-ws/facetedResults=github.com/indexdata/xsd2goxsl/xsd/facetedResults.xsd_test,\
http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis=github.com/indexdata/xsd2goxsl/xsd/searchResultAnalysis.xsd_test,\
http://docs.oasis-open.org/ns/search-ws/xcql=github.com/indexdata/xsd2goxsl/xsd/xcql.xsd_test,\
http://docs.oasis-open.org/ns/search-ws/scan=github.com/indexdata/xsd2goxsl/xsd/scan.xsd_test
endef

define qAttr
"qAttrImport=utils \"github.com/indexdata/go-utils/utils\"" qAttrType=utils.PrefixAttr
endef

check: $(XSD2GOXSL) $(GOFILES) $(XSD_FILES)
	$(foreach file, $(XSD_FILES), rm -rf $(file)_test;)
	$(foreach file, $(XSD_FILES), mkdir $(file)_test;)
	$(foreach file, $(XSD_FILES), $(GO) run $(MAIN_PACKAGE) $(file) $(file)_test/schema.go "$(nsImports)" ${qAttr} buildtag=checkxsd;)
	$(foreach file, $(XSD_FILES), diff -u $(file)_test/schema.go $(file).out.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) vet -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) build -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), rm -rf $(file)_test;)
	$(foreach file, $(XSD_FILES), mkdir $(file)_test;)
	$(foreach file, $(XSD_FILES), $(GO) run $(MAIN_PACKAGE) $(file) $(file)_test/schema.go "$(nsImports)" ${qAttr} buildtag=checkxsd json=yes;)
	$(foreach file, $(XSD_FILES), diff -u $(file)_test/schema.go $(file).out.json.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) vet -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) build -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), rm -rf $(file)_test;)
	$(foreach file, $(XSD_FILES), mkdir $(file)_test;)
	$(foreach file, $(XSD_FILES), $(GO) run $(MAIN_PACKAGE) $(file) $(file)_test/schema.go "$(nsImports)" ${qAttr} buildtag=checkxsd namespaced=yes;)
	$(foreach file, $(XSD_FILES), diff -u $(file)_test/schema.go $(file).out.ns.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) vet -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) build -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), rm -rf $(file)_test;)
	$(foreach file, $(XSD_FILES), mkdir $(file)_test;)
	$(foreach file, $(XSD_FILES), $(GO) run $(MAIN_PACKAGE) $(file) $(file)_test/schema.go "$(nsImports)" ${qAttr} buildtag=checkxsd validate=yes;)
	$(foreach file, $(XSD_FILES), diff -u $(file)_test/schema.go $(file).out.validate.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) vet -tags checkxsd $(file)_test/schema.go || exit;)
	$(foreach file, $(XSD_FILES), $(GO) build -tags checkxsd $(file)_test/schema.go || exit;)
