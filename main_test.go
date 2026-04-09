package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateTags(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "validate.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="http://example.com/validate">
  <xs:simpleType name="state">
    <xs:restriction base="xs:string">
      <xs:enumeration value="NEW"/>
      <xs:enumeration value="DONE"/>
    </xs:restriction>
  </xs:simpleType>
  <xs:simpleType name="label">
    <xs:restriction base="xs:string">
      <xs:enumeration value="loaned items"/>
      <xs:enumeration value="requested items"/>
    </xs:restriction>
  </xs:simpleType>
  <xs:complexType name="child">
    <xs:sequence>
      <xs:element name="code" type="xs:string"/>
      <xs:element name="count" type="xs:int"/>
      <xs:element name="enabled" type="xs:boolean"/>
    </xs:sequence>
  </xs:complexType>
  <xs:element name="root">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="title" type="xs:string"/>
        <xs:element name="nickname" type="xs:string" minOccurs="0"/>
        <xs:element name="status" type="state" minOccurs="0"/>
        <xs:element name="states" type="state" minOccurs="0" maxOccurs="3"/>
        <xs:element name="requiredStates" type="state" maxOccurs="3"/>
        <xs:element name="labels" type="label" minOccurs="0" maxOccurs="3"/>
        <xs:element name="child" type="child"/>
      </xs:sequence>
      <xs:attribute name="id" type="xs:string" use="required"/>
    </xs:complexType>
  </xs:element>
  <xs:element name="choiceRoot">
    <xs:complexType>
      <xs:choice>
        <xs:element name="email" type="xs:string"/>
        <xs:element name="phone" type="xs:string"/>
      </xs:choice>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{schemaFile, outputFile, "validate=yes"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `Title string `+"`"+`xml:"title" validate:"required"`+"`")
	assertContains(t, out, `Nickname string `+"`"+`xml:"nickname,omitempty"`+"`")
	assertContains(t, out, `Status *State `+"`"+`xml:"status,omitempty" validate:"omitempty,oneof=NEW DONE"`+"`")
	assertContains(t, out, `States []State `+"`"+`xml:"states,omitempty" validate:"omitempty,max=3,dive,oneof=NEW DONE"`+"`")
	assertContains(t, out, `RequiredStates []State `+"`"+`xml:"requiredStates,omitempty" validate:"min=1,max=3,dive,oneof=NEW DONE"`+"`")
	assertContains(t, out, `Labels []Label `+"`"+`xml:"labels,omitempty" validate:"omitempty,max=3,dive,oneof='loaned items' 'requested items'"`+"`")
	assertContains(t, out, `Child Child `+"`"+`xml:"child"`+"`")
	assertContains(t, out, `Id string `+"`"+`xml:"id,attr"`+"`")
	assertContains(t, out, `Count int `+"`"+`xml:"count"`+"`")
	assertContains(t, out, `Enabled bool `+"`"+`xml:"enabled"`+"`")
	assertContains(t, out, `Code string `+"`"+`xml:"code" validate:"required"`+"`")
	assertNotContains(t, out, `Child Child `+"`"+`xml:"child" validate:"required"`+"`")
	assertNotContains(t, out, `Id string `+"`"+`xml:"id,attr" validate:"required"`+"`")
	assertNotContains(t, out, `Count int `+"`"+`xml:"count" validate:"required"`+"`")
	assertNotContains(t, out, `Enabled bool `+"`"+`xml:"enabled" validate:"required"`+"`")
	assertContains(t, out, `Email string `+"`"+`xml:"email,omitempty"`+"`")
	assertContains(t, out, `Phone string `+"`"+`xml:"phone,omitempty"`+"`")
	assertNotContains(t, out, `Email string `+"`"+`xml:"email,omitempty" validate:"required"`+"`")
	assertNotContains(t, out, `Phone string `+"`"+`xml:"phone,omitempty" validate:"required"`+"`")
	assertNotContains(t, out, `Email string `+"`"+`xml:"email,omitempty" validate:"min=1"`+"`")
	assertNotContains(t, out, `Phone string `+"`"+`xml:"phone,omitempty" validate:"min=1"`+"`")
}

func TestValidateTagsMergeDuplicateChoiceFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "duplicate-choice-fields.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="http://example.com/choice-merge">
  <xs:element name="ItemId" type="xs:string"/>
  <xs:element name="BibliographicId" type="xs:string"/>
  <xs:element name="requestItem">
    <xs:complexType>
      <xs:sequence>
        <xs:choice>
          <xs:element ref="ItemId" maxOccurs="unbounded"/>
          <xs:sequence>
            <xs:element ref="BibliographicId" maxOccurs="unbounded"/>
            <xs:element ref="ItemId" minOccurs="0" maxOccurs="unbounded"/>
          </xs:sequence>
        </xs:choice>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{schemaFile, outputFile, "validate=yes"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `ItemId []string `+"`"+`xml:"ItemId,omitempty"`+"`")
	assertContains(t, out, `BibliographicId []string `+"`"+`xml:"BibliographicId,omitempty"`+"`")
	assertNotContains(t, out, `BibliographicId []string `+"`"+`xml:"BibliographicId,omitempty" validate:"min=1"`+"`")
	if strings.Count(out, `ItemId []string `+"`"+`xml:"ItemId,omitempty"`) != 1 {
		t.Fatalf("expected exactly one merged ItemId field, got output:\n%s", out)
	}
}

func TestValidateTagsChoiceBranchFieldsBecomeOptional(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "choice-branch-fields.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="http://example.com/choice-branch-fields">
  <xs:element name="root">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="header" type="xs:string"/>
        <xs:choice>
          <xs:element name="itemId" type="xs:string"/>
          <xs:sequence>
            <xs:element name="bibliographicId" type="xs:string"/>
            <xs:element name="requestId" type="xs:string"/>
          </xs:sequence>
        </xs:choice>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{schemaFile, outputFile, "validate=yes"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `Header string `+"`"+`xml:"header" validate:"required"`+"`")
	assertContains(t, out, `ItemId string `+"`"+`xml:"itemId,omitempty"`+"`")
	assertContains(t, out, `BibliographicId string `+"`"+`xml:"bibliographicId,omitempty"`+"`")
	assertContains(t, out, `RequestId string `+"`"+`xml:"requestId,omitempty"`+"`")
	assertNotContains(t, out, `ItemId string `+"`"+`xml:"itemId,omitempty" validate:"required"`+"`")
	assertNotContains(t, out, `BibliographicId string `+"`"+`xml:"bibliographicId,omitempty" validate:"required"`+"`")
	assertNotContains(t, out, `RequestId string `+"`"+`xml:"requestId,omitempty" validate:"required"`+"`")
}

func TestNamespaceImports(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "namespace-imports.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  xmlns:imp="http://example.com/imported"
  targetNamespace="http://example.com/local">
  <xs:import namespace="http://example.com/imported"/>
  <xs:complexType name="holder">
    <xs:sequence>
      <xs:element name="remote" type="imp:remote_type"/>
    </xs:sequence>
  </xs:complexType>
  <xs:element name="root" type="holder"/>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"namespaceImports=http://example.com/imported=example.com/remote/pkg",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `imp "example.com/remote/pkg"`)
	assertContains(t, out, `Remote imp.RemoteType `+"`"+`xml:"remote"`+"`")
}

func TestRootLimitedNamespaceGeneration(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "root-ns.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/root"
  attributeFormDefault="qualified">
  <xs:element name="root">
    <xs:complexType>
      <xs:attribute name="id" type="xs:string" use="required"/>
    </xs:complexType>
  </xs:element>
  <xs:element name="other">
    <xs:complexType>
      <xs:attribute name="id" type="xs:string" use="required"/>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"namespaced=yes",
		"root=root",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `type Root struct {`)
	assertContains(t, out, `XMLName xml.Name `+"`"+`xml:"http://example.com/root root"`+"`")
	assertContains(t, out, `type Other struct {`)
	assertContains(t, out, `XMLName xml.Name `+"`"+`xml:"other"`+"`")
	assertContains(t, out, `Id string `+"`"+`xml:"http://example.com/root id,attr"`+"`")
	assertNotContains(t, out, `XMLName xml.Name `+"`"+`xml:"http://example.com/root other"`+"`")
}

func TestMultipleRootLimitedNamespaceGeneration(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "multi-root-ns.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/root">
  <xs:element name="rootA">
    <xs:complexType/>
  </xs:element>
  <xs:element name="rootB">
    <xs:complexType/>
  </xs:element>
  <xs:element name="other">
    <xs:complexType/>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"namespaced=yes",
		"root=rootA, rootB",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `XMLName xml.Name `+"`"+`xml:"http://example.com/root rootA"`+"`")
	assertContains(t, out, `XMLName xml.Name `+"`"+`xml:"http://example.com/root rootB"`+"`")
	assertContains(t, out, `XMLName xml.Name `+"`"+`xml:"other"`+"`")
	assertNotContains(t, out, `XMLName xml.Name `+"`"+`xml:"http://example.com/root other"`+"`")
}

func TestRootTagConstFromRootParam(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "root-tag-param.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
  <xs:element name="rootA">
    <xs:complexType/>
  </xs:element>
  <xs:element name="rootB">
    <xs:complexType/>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"root=rootB",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)
	assertContains(t, out, `ROOT_TAG = "rootB"`)
}

func TestRootTagConstDefaultsToFirstGlobalElement(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "root-tag-default.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
  <xs:element name="firstRoot">
    <xs:complexType/>
  </xs:element>
  <xs:element name="secondRoot">
    <xs:complexType/>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)
	assertNotContains(t, out, `ROOT_TAG = "`)
}

func TestSchemaLocationMarshalXMLForSelectedRoots(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "schema-location.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/schema-location">
  <xs:element name="root">
    <xs:complexType/>
  </xs:element>
  <xs:element name="other">
    <xs:complexType/>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"root=root",
		"schemaLocation=schema.xsd",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `XMLNS_XSI = "http://www.w3.org/2001/XMLSchema-instance"`)
	assertContains(t, out, `XSI_SCHEMA_LOCATION = "http://example.com/schema-location schema.xsd"`)
	assertNotContains(t, out, `TARGET_NAMESPACE = "`)
	assertContains(t, out, `func (x *Root) MarshalXML(e *xml.Encoder, start xml.StartElement) error {`)
	assertContains(t, out, `xml.Attr{Name: xml.Name{Local: "xmlns:xsi"}, Value: XMLNS_XSI},`)
	assertContains(t, out, `xml.Attr{Name: xml.Name{Local: "xsi:schemaLocation"}, Value: XSI_SCHEMA_LOCATION},`)
	assertContains(t, out, `return e.EncodeElement((*Alias)(x), start)`)
	assertNotContains(t, out, `func (x *Other) MarshalXML(e *xml.Encoder, start xml.StartElement) error {`)
}

func TestSchemaLocationRawPairsPreserved(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "schema-location-raw.xsd")
	outputFile := filepath.Join(dir, "schema.go")
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/raw">
  <xs:element name="root">
    <xs:complexType/>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"root=root",
		"schemaLocation=http://a/ns a.xsd http://b/ns b.xsd",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	out := string(generated)

	assertContains(t, out, `XSI_SCHEMA_LOCATION = "http://a/ns a.xsd http://b/ns b.xsd"`)
}

func TestSchemaLocationMarshalXMLRuntime(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "runtime-schema.xsd")
	genDir := filepath.Join(dir, "gen")
	outputFile := filepath.Join(genDir, "schema.go")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		t.Fatal(err)
	}
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/runtime">
  <xs:element name="root">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="value" type="xs:string" minOccurs="0"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"package=gen",
		"namespaced=yes",
		"root=root",
		"schemaLocation=runtime.xsd",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}
	generated, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	genOut := string(generated)
	assertContains(t, genOut, `TARGET_NAMESPACE = "http://example.com/runtime"`)
	assertContains(t, genOut, `start.Name.Space = TARGET_NAMESPACE`)

	goMod := `module example.com/runtime-test

go 1.22
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	mainGo := `package main

import (
	"encoding/xml"
	"fmt"

	"example.com/runtime-test/gen"
)

func main() {
	doc := gen.Root{Value: "hello"}
	out, err := xml.Marshal(&doc)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(out))
	}
	xmlOut := string(out)

	assertContains(t, xmlOut, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
	assertContains(t, xmlOut, `xsi:schemaLocation="http://example.com/runtime runtime.xsd"`)
	assertContains(t, xmlOut, `<root`)
	assertContains(t, xmlOut, `xmlns="http://example.com/runtime"`)
}

func TestSchemaLocationMarshalXMLRuntimeJSONRoundtrip(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "runtime-schema-json.xsd")
	genDir := filepath.Join(dir, "gen")
	outputFile := filepath.Join(genDir, "schema.go")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		t.Fatal(err)
	}
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/runtime-json">
  <xs:element name="root">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="value" type="xs:string" minOccurs="0"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"package=gen",
		"namespaced=yes",
		"root=root",
		"schemaLocation=runtime-json.xsd",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	goMod := `module example.com/runtime-json-test

go 1.22
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	mainGo := `package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"example.com/runtime-json-test/gen"
)

func main() {
	var doc gen.Root
	if err := json.Unmarshal([]byte("{\"value\":\"hello\"}"), &doc); err != nil {
		panic(err)
	}
	out, err := xml.Marshal(&doc)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(out))
	}
	xmlOut := string(out)

	assertContains(t, xmlOut, `<root`)
	assertContains(t, xmlOut, `xmlns="http://example.com/runtime-json"`)
	assertContains(t, xmlOut, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
	assertContains(t, xmlOut, `xsi:schemaLocation="http://example.com/runtime-json runtime-json.xsd"`)
}

func TestSchemaLocationMarshalXMLRuntimeNoNamespace(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	schemaFile := filepath.Join(dir, "runtime-schema-no-ns.xsd")
	genDir := filepath.Join(dir, "gen")
	outputFile := filepath.Join(genDir, "schema.go")
	if err := os.MkdirAll(genDir, 0o755); err != nil {
		t.Fatal(err)
	}
	schema := `<?xml version="1.0" encoding="UTF-8"?>
<xs:schema
  xmlns:xs="http://www.w3.org/2001/XMLSchema"
  targetNamespace="http://example.com/runtime-no-ns">
  <xs:element name="root">
    <xs:complexType>
      <xs:sequence>
        <xs:element name="value" type="xs:string" minOccurs="0"/>
      </xs:sequence>
    </xs:complexType>
  </xs:element>
</xs:schema>
`
	if err := os.WriteFile(schemaFile, []byte(schema), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := run([]string{
		schemaFile,
		outputFile,
		"package=gen",
		"namespaced=no",
		"root=root",
		"schemaLocation=runtime-no-ns.xsd",
	}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("run failed with exit code %d: %s", exitCode, stderr.String())
	}

	goMod := `module example.com/runtime-no-ns-test

go 1.22
`
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}

	mainGo := `package main

import (
	"encoding/xml"
	"fmt"

	"example.com/runtime-no-ns-test/gen"
)

func main() {
	doc := gen.Root{Value: "hello"}
	out, err := xml.Marshal(&doc)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(out))
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(out))
	}
	xmlOut := string(out)

	assertContains(t, xmlOut, `<root`)
	assertNotContains(t, xmlOut, `xmlns="http://example.com/runtime-no-ns"`)
	assertContains(t, xmlOut, `xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`)
	assertContains(t, xmlOut, `xsi:schemaLocation="http://example.com/runtime-no-ns runtime-no-ns.xsd"`)
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("generated output missing substring:\n%s\n\nfull output:\n%s", want, got)
	}
}

func assertNotContains(t *testing.T, got, want string) {
	t.Helper()
	if strings.Contains(got, want) {
		t.Fatalf("generated output unexpectedly contained substring:\n%s\n\nfull output:\n%s", want, got)
	}
}
