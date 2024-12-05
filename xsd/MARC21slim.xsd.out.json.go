// Code generated by xsd2go.xsl; DO NOT EDIT.
//go:build checkxsd

package slim

import (
  "encoding/xml"
)

// 			MARCXML: The MARC 21 XML Schema
// 			Prepared by Corey Keith
// 			
// 				May 21, 2002 - Version 1.0  - Initial Release
// **********************************************
// Changes.
// August 4, 2003 - Version 1.1 - 
// Removed import of xml namespace and the use of xml:space="preserve" attributes on the leader and controlfields. 
//                     Whitespace preservation in these subfields is accomplished by the use of xsd:whiteSpace value="preserve"
// May 21, 2009  - Version 1.2 - 
// in subfieldcodeDataType  the pattern 
//                           "[\da-z!"#$%&'()*+,-./:;<=>?{}_^`~\[\]\\]{1}"
// 	changed to:	
//                          "[\dA-Za-z!"#$%&'()*+,-./:;<=>?{}_^`~\[\]\\]{1}"
//     i.e "A-Z" added after "[\d" before "a-z"  to allow upper case.  This change is for consistency with the documentation.
// 	
// ************************************************************
// 			This schema supports XML markup of MARC21 records as specified in the MARC documentation (see www.loc.gov).  It allows tags with
// 			alphabetics and subfield codes that are symbols, neither of which are as yet used in  the MARC 21 communications formats, but are 
// 			allowed by MARC 21 for local data.  The schema accommodates all types of MARC 21 records: bibliographic, holdings, bibliographic 
// 			with embedded holdings, authority, classification, and community information.
// 		
type Record struct {
  XMLName xml.Name `xml:"record" json:"-"`
  RecordType
}

type Collection struct {
  XMLName xml.Name `xml:"collection" json:"-"`
  CollectionType
}

type CollectionType struct {
  Record []RecordType `xml:"record,omitempty" json:"record,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
}

type RecordType struct {
  Leader *LeaderFieldType `xml:"leader,omitempty" json:"leader,omitempty"`
  Controlfield []ControlFieldType `xml:"controlfield,omitempty" json:"controlfield,omitempty"`
  Datafield []DataFieldType `xml:"datafield,omitempty" json:"datafield,omitempty"`
  Type string `xml:"type,attr,omitempty" json:"@type,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
}

type RecordTypeType string

const RecordTypeTypeBibliographic RecordTypeType = "Bibliographic"
const RecordTypeTypeAuthority RecordTypeType = "Authority"
const RecordTypeTypeHoldings RecordTypeType = "Holdings"
const RecordTypeTypeClassification RecordTypeType = "Classification"
const RecordTypeTypeCommunity RecordTypeType = "Community"

type LeaderFieldType struct {
  // MARC21 Leader, 24 bytes
  Text LeaderDataType `xml:",chardata" json:"#text,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
}

type LeaderDataType string


type ControlFieldType struct {
  // MARC21 Fields 001-009
  Text ControlDataType `xml:",chardata" json:"#text,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
  Tag string `xml:"tag,attr" json:"@tag"`
}

type ControlDataType string


type ControltagDataType string


type DataFieldType struct {
  // MARC21 Variable Data Fields 010-999
  Subfield []SubfieldatafieldType `xml:"subfield,omitempty" json:"subfield,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
  Tag string `xml:"tag,attr" json:"@tag"`
  Ind1 string `xml:"ind1,attr" json:"@ind1"`
  Ind2 string `xml:"ind2,attr" json:"@ind2"`
}

type TagDataType string


type IndicatorDataType string


type SubfieldatafieldType struct {
  Text SubfieldDataType `xml:",chardata" json:"#text,omitempty"`
  Id string `xml:"id,attr,omitempty" json:"@id,omitempty"`
  Code string `xml:"code,attr" json:"@code"`
}

type SubfieldDataType string


type SubfieldcodeDataType string

//  "A-Z" added after "\d" May 21, 2009 

type IdDataType string


