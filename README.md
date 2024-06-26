# xsd2go.xsl

xsd2go.xsl is an XSLT stylesheet that converts XML Schema to Go type definitions. Not all XSD elements are supported, but the most
popular are:

*  global and local `element` definitions and declarations, local definitions are hoisted to the package level
* `attribute`
* `simpleType`, `complexType` and `simpleContent`, `complexContent`
* `sequence` and `choice`
* `extension`, `restriction` and `enumeration`
* `annotation` and `documentation`

# Use

The stylesheet uses XSLT 1.0 plus a `tokenize` function from EXSLT as supported by `xsltproc` (libxml2)

```
xsltproc xsd2go.xsl some.xsd > some.go
```

The following parameters are supported (`--stringparam` in `xsltproc`):

* `indent` indent string, default 2 spaces
* `break` line break char, default `&#10;` or CR (carriage return)
* `debug` write out schema types in comments, default 'no'
* `buildtag` //go:build tags for the generated file, default empty
* `targetNamespace` defaults to `/xs:schema/@targetNamespace`
* `package` Go package name, defaults to `str:tokenize(str:tokenize($targetNamespace, '/')[last()],'.')[1]`
* `omitempty` whether to set _omitempty_ modifier on field tags for optional and repeating elements, default `yes`
* `xmlns_xsd` defaults to `http://www.w3.org/2001/XMLSchema`
* `attributeForm` _qualified_ or _unqualified_, defaults to `xs:schema/@attributeFormDefault`
* `qAttrType` data type to use for qualified attributes, defaults to `string` which is non-qualified, see [PrefixAttr](https://github.com/indexdata/edge-slnp/blob/xsd-2-go-xsl/utils/xml.go)
* `qAttrImport` import for the qualified attributes type package, e.g `xmlext "github.com/indexdata/edge-slnp/utils"`. Note that this import is applied only when `attributeForm` is _qualified_
* `dateTimeType` data type to use for XSD dateTime, defaults to `string`
* `timeType` data type to use for XSD time, defaults to `string`
* `dateType` data type to use for XSD date, defaults to `string`
* `decimalType` data type to use for XSD decimal, defaults to `float64`
* `typeImports` comma-separated list of imports for additional type definitions used for date/time and decimal types

You can also see rendered Go structs in the browser by prepending:

```
<?xml-stylesheet type="text/xsl" href="xsd2go.xsl"?>
```

to the XSD and opening the file in the browser, best by serving it with a local HTTP server like `python3 -m http.server` to avoid local security constraints. This works in Firefox but fails silently in Chrome, likely because of missing EXSLT support.

# Use in a Go build

This repo includes a simple Go wrapper over `xsltproc` which you can run with:

```
go run github.com/indexdata/xsd2goxsl <in.xsd> <out.go> <param-name=param-value>,...
```

e.g

```
go run github.com/indexdata/xsd2goxsl xsd/ncip_v2_02.xsd ncip/schema.go "qAttrImport=utils \"github.com/indexdata/go-utils/utils\"" qAttrType=utils.PrefixAttr dateTimeType=utils.XSDDateTime
```

This allows using it in a Go project during the build with `go generate`. E.g by adding a Go `xsd-gen.go` file to the project with:

```
package ncip

//go:generate go  run github.com/indexdata/xsd2goxsl xsd/ncip_v2_02.xsd ncip/schema.go "qAttrImport=utils \"github.com/indexdata/go-utils/utils\"" qAttrType=utils.PrefixAttr dateTimeType=utils.XSDDateTime decimalType=utils.XSDDecimal
```

and running:

```
go generate
```

Additionally, you can add a disabled source file with an import for this project to force Go handling it as a dependency:

```
//go:build tools

package tools

//build-time toolchain dependencies
import (
	_ "github.com/indexdata/xsd2goxsl"
)
```

# Test

There are example XSDs and corresponding generated Go models under [./xsd](xsd/).
