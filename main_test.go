package main

import (
	"bytes"
	"os"
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
