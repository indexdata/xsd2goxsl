// Code generated by xsd2go.xsl; DO NOT EDIT.
//go:build checkxsd

//     Search Web Services searchRetrieve Version 1.0
//     OASIS Standard
//     30 January 2013
//     Source: http://docs.oasis-open.org/search-ws/searchRetrieve/v1.0/os/schemas/
//     Copyright (c) OASIS Open 2013.  All Rights Reserved.
// Editor: Ray Denenberg, Library of Congress.  rden@loc.gov
// *****
package searchResultAnalysis

import (
  "encoding/xml"
)

type SearchResultAnalysis struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis searchResultAnalysis"`
  SearchResultAnalysisDefinition
}

type SearchResultAnalysisDefinition struct {
  Datasource []DatasourceDefinition `xml:"datasource,omitempty"`
  SubqueryResult []SubqueryResultDefinition `xml:"subqueryResult,omitempty"`
}

type DatasourceDefinition struct {
  DatasourceDisplayLabel string `xml:"datasourceDisplayLabel,omitempty"`
  DatasourceDescription string `xml:"datasourceDescription,omitempty"`
  BaseURL string `xml:"baseURL,omitempty"`
  SubqueryResults SubqueryResultsDefinition `xml:"subqueryResults"`
  Full string `xml:"full,attr,omitempty"`
}

type SubqueryResultsDefinition struct {
  SubqueryResult []SubqueryResultDefinition `xml:"subqueryResult,omitempty"`
}

type SubqueryResultDefinition struct {
  SubqueryDisplayLabel string `xml:"subqueryDisplayLabel,omitempty"`
  Subquery string `xml:"subquery"`
  Count int64 `xml:"count"`
  RequestUrl string `xml:"requestUrl,omitempty"`
  Full string `xml:"full,attr,omitempty"`
}

type FullDefinition string

const FullDefinitionTrue FullDefinition = "true"

type BaseURL struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis baseURL"`
  Text string `xml:",chardata"`
}

type Count struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis count"`
  Text int64 `xml:",chardata"`
}

type Datasource struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis datasource"`
  DatasourceDefinition
}

type DatasourceDescription struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis datasourceDescription"`
  Text string `xml:",chardata"`
}

type DatasourceDisplayLabel struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis datasourceDisplayLabel"`
  Text string `xml:",chardata"`
}

type RequestUrl struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis requestUrl"`
  Text string `xml:",chardata"`
}

type Subquery struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis subquery"`
  Text string `xml:",chardata"`
}

type SubqueryDisplayLabel struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis subqueryDisplayLabel"`
  Text string `xml:",chardata"`
}

type SubqueryResult struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis subqueryResult"`
  SubqueryResultDefinition
}

type SubqueryResults struct {
  XMLName xml.Name `xml:"http://docs.oasis-open.org/ns/search-ws/searchResultAnalysis subqueryResults"`
  SubqueryResultsDefinition
}

