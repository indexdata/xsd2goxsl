// Code generated by xsd2go.xsl; DO NOT EDIT.
//go:build checkxsd

//     Search Web Services searchRetrieve Version 1.0
//     OASIS Standard
//     30 January 2013
//     Source: http://docs.oasis-open.org/search-ws/searchRetrieve/v1.0/os/schemas/
//     Copyright (c) OASIS Open 2013.  All Rights Reserved.
// Editor: Ray Denenberg, Library of Congress.  rden@loc.gov
// *****
package xcql

import (
  "encoding/xml"
)

type Xcql struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql xcql"`
  XcqlDefinition
}

type XcqlDefinition struct {
  Prefixes *PrefixesDefinition `xml:"prefixes,omitempty"`
  Triple TripleDefinition `xml:"triple"`
  SortKeys *SortKeysDefinition `xml:"sortKeys,omitempty"`
}

type PrefixesDefinition struct {
  Prefix []PrefixDefinition `xml:"prefix,omitempty"`
}

type PrefixDefinition struct {
  Name string `xml:"name"`
  Identifier string `xml:"identifier"`
}

type TripleDefinition struct {
  SearchClause *SearchClauseDefinition `xml:"searchClause,omitempty"`
  Boolean *BooleanPlusModifier `xml:"Boolean,omitempty"`
  LeftOperand *OperandDefinition `xml:"leftOperand,omitempty"`
  RightOperand *OperandDefinition `xml:"rightOperand,omitempty"`
}

type SortKeysDefinition struct {
  Key []KeyDefinition `xml:"key,omitempty"`
}

type BooleanPlusModifier struct {
  Value BooleanValue `xml:"value"`
  Modifiers *ModifiersDefinition `xml:"modifiers,omitempty"`
}

type BooleanValue string

const BooleanValueAnd BooleanValue = "and"
const BooleanValueOr BooleanValue = "or"
const BooleanValueNot BooleanValue = "not"
const BooleanValueProx BooleanValue = "prox"

type KeyDefinition struct {
  Index string `xml:"index"`
  Modifiers ModifiersDefinition `xml:"modifiers"`
}

type ModifierDefinition struct {
  Type string `xml:"type"`
  Comparison string `xml:"comparison,omitempty"`
  Value string `xml:"value,omitempty"`
}

type ModifiersDefinition struct {
  Modifier []ModifierDefinition `xml:"modifier,omitempty"`
}

type OperandDefinition struct {
  SearchClause *SearchClauseDefinition `xml:"searchClause,omitempty"`
  Triple *TripleDefinition `xml:"triple,omitempty"`
}

type SearchClauseDefinition struct {
  Term string `xml:"term,omitempty"`
  Index string `xml:"index,omitempty"`
  Relation *ValuePlusModifier `xml:"relation,omitempty"`
}

type ValuePlusModifier struct {
  Value string `xml:"value"`
  Modifiers *ModifiersDefinition `xml:"modifiers,omitempty"`
}

type Boolean struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql Boolean"`
  BooleanPlusModifier
}

type Comparison struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql comparison"`
  Text string `xml:",chardata"`
}

type Identifier struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql identifier"`
  Text string `xml:",chardata"`
}

type Index struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql index"`
  Text string `xml:",chardata"`
}

type Key struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql key"`
  KeyDefinition
}

type LeftOperand struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql leftOperand"`
  OperandDefinition
}

type Modifier struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql modifier"`
  ModifierDefinition
}

type Modifiers struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql modifiers"`
  ModifiersDefinition
}

type Name struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql name"`
  Text string `xml:",chardata"`
}

type Prefix struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql prefix"`
  PrefixDefinition
}

type Prefixes struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql prefixes"`
  PrefixesDefinition
}

type Relation struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql relation"`
  ValuePlusModifier
}

type RightOperand struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql rightOperand"`
  OperandDefinition
}

type SearchClause struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql searchClause"`
  SearchClauseDefinition
}

type SortKeys struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql sortKeys"`
  SortKeysDefinition
}

type Term struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql term"`
  Text string `xml:",chardata"`
}

type Type struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql type"`
  Text string `xml:",chardata"`
}

type Triple struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql triple"`
  TripleDefinition
}

type Value struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/xcql value"`
  Text string `xml:",chardata"`
}

