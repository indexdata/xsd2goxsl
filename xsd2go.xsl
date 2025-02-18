<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet
    version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xs="http://www.w3.org/2001/XMLSchema"
    xmlns:str="http://exslt.org/strings">

  <xsl:output method="text" encoding="UTF-8"/>
  <xsl:strip-space elements="*" />
  <xsl:param name="indent" select="'  '"/>
  <xsl:param name="break" select="'&#10;'"/>
  <xsl:param name="attributeForm" select="/xs:schema/@attributeFormDefault"/>
  <xsl:param name="targetNamespace" select="/xs:schema/@targetNamespace"/>
  <xsl:param name="package" select="str:tokenize(str:tokenize($targetNamespace, '/')[last()],'.')[1]"/>
  <xsl:param name="qAttrType" select="'string'"/>
  <xsl:param name="qAttrImport" select="''"/>
  <xsl:param name="dateTimeType" select="''"/>
  <xsl:param name="dateType" select="''"/>
  <xsl:param name="timeType" select="''"/>
  <xsl:param name="decimalType" select="''"/>
  <xsl:param name="typeImports" select="''"/>
  <xsl:param name="omitempty" select="'yes'"/>
  <xsl:param name="xmlns_xsd" select="'http://www.w3.org/2001/XMLSchema'"/>
  <xsl:param name="debug" select="'no'"/>
  <xsl:param name="header" select="'// Code generated by xsd2go.xsl; DO NOT EDIT.'" />
  <xsl:param name="buildtag" select="''" />
  <xsl:param name="json" select="'no'"/>

  <xsl:template match="/">
    <xsl:if test="$header">
      <xsl:value-of select="$header"/>
      <xsl:value-of select="$break"/>
    </xsl:if>
    <xsl:if test="$buildtag">
      <xsl:value-of select="'//go:build '"/>
      <xsl:value-of select="$buildtag"/>
      <xsl:value-of select="$break"/>
      <xsl:value-of select="$break"/>
    </xsl:if>
    <xsl:if test="$buildtag = ''">
      <xsl:value-of select="$break"/>
    </xsl:if>
    <xsl:apply-templates/>
  </xsl:template>

  <xsl:template match="xs:schema">
    <xsl:call-template name="debug">
      <xsl:with-param name="text" select="'Schema'"/>
    </xsl:call-template>
    <xsl:text>package </xsl:text>
    <xsl:value-of select="$package"/>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$break"/>
    <xsl:text>import (</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$indent"/>
    <xsl:text>"encoding/xml"</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:if test="$qAttrImport and $attributeForm = 'qualified'">
      <xsl:value-of select="$indent"/>
      <xsl:value-of select="$qAttrImport"/>
      <xsl:value-of select="$break"/>
    </xsl:if>
    <xsl:call-template name="type-imports"/>
    <xsl:text>)</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$break"/>
    <xsl:apply-templates mode="global"/>
    <xsl:for-each select="//xs:element[(@minOccurs!='1' or @maxOccurs!='1') and (xs:complexType or xs:simpleType)
                                        or (not(parent::xs:schema) and xs:simpleType/xs:restriction/xs:enumeration)]">
      <xsl:call-template name="debug">
        <xsl:with-param name="text" select="'Hoisted element'"/>
      </xsl:call-template>
      <xsl:value-of select="$break"/>
      <xsl:apply-templates select=".">
        <xsl:with-param name="hoisted" select="true()"/>
      </xsl:apply-templates>
    </xsl:for-each>
  </xsl:template>

  <xsl:template name="type-imports">
    <xsl:for-each select="str:tokenize($typeImports, ',')">
      <xsl:value-of select="$indent"/>
      <xsl:value-of select="."/>
      <xsl:value-of select="$break"/>
    </xsl:for-each>
  </xsl:template>

  <!-- global elem with nested types, apply recursively -->
  <xsl:template match="xs:element[xs:complexType or xs:simpleType]" mode="global">
    <xsl:value-of select="$break"/>
    <xsl:apply-templates select="."/>
  </xsl:template>

  <!-- global elem w/o nested types, embed type and stop -->
  <xsl:template match="xs:element[not(xs:complexType or xs:simpleType)]" mode="global">
    <xsl:param name="name" select="@name" />
    <xsl:param name="type" select="@type" />
    <xsl:call-template name="debug">
      <xsl:with-param name="text" select="'Element'"/>
    </xsl:call-template>
    <xsl:text>type </xsl:text>
    <xsl:call-template name="convert-name">
      <xsl:with-param name="name" select="$name"/>
    </xsl:call-template>
    <xsl:text> struct {</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$indent"/>
    <xsl:text>XMLName xml.Name `xml:"</xsl:text>
    <xsl:value-of select="$name"/>
    <xsl:if test="$json = 'yes'">
      <xsl:text>" json:"-</xsl:text>
    </xsl:if>
    <xsl:text>"`</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$indent"/>
    <xsl:call-template name="nest-type">
      <xsl:with-param name="type" select="$type"/>
    </xsl:call-template>
    <xsl:value-of select="$break"/>
    <xsl:text>}</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$break"/>
  </xsl:template>

  <xsl:template name="nest-type">
    <xsl:param name="type"/>
    <xsl:variable name="ns">
      <xsl:call-template name="get-ns">
        <xsl:with-param name="name" select="$type"/>
      </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="convType">
      <xsl:call-template name="convert-type">
        <xsl:with-param name="type" select="$type"/>
      </xsl:call-template>
    </xsl:variable>
    <xsl:choose>
      <!-- embed simple type, assumes all XSD types are simple -->
      <xsl:when test="$ns = $xmlns_xsd">
        <xsl:text>Text </xsl:text>
        <xsl:value-of select="$convType"/>
        <xsl:text> `xml:",chardata</xsl:text>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"#text</xsl:text>
          <xsl:if test="$omitempty = 'yes'">
            <xsl:text>,omitempty</xsl:text>
          </xsl:if>
        </xsl:if>
        <xsl:text>"`</xsl:text>
      </xsl:when>
      <!-- embed complex type, assume all local schema types are complex-->
      <xsl:otherwise>
        <xsl:value-of select="$convType"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

  <!-- local or hoisted elem with nested type, apply recursively-->
  <xsl:template match="xs:element[xs:complexType or xs:simpleType]">
    <xsl:param name="level"/>
    <xsl:param name="hoisted"/>
    <xsl:variable name="ptr">
      <xsl:call-template name="get-cardinality">
        <xsl:with-param name="minOccurs" select="@minOccurs"/>
        <xsl:with-param name="maxOccurs" select="@maxOccurs"/>
        <xsl:with-param name="container" select="'element'"/>
        </xsl:call-template>
    </xsl:variable>
    <xsl:apply-templates>
      <xsl:with-param name="name" select="@name"/>
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="parentPtr" select="$ptr"/>
    </xsl:apply-templates>
  </xsl:template>

  <!-- local elem w/o nested type -->
  <xsl:template match="xs:element[not(xs:complexType or xs:simpleType)]">
    <xsl:param name="level" select="''"/>
    <xsl:param name="parentPtr"/>
    <xsl:variable name="ptr">
      <xsl:call-template name="get-cardinality">
        <xsl:with-param name="minOccurs" select="@minOccurs"/>
        <xsl:with-param name="maxOccurs" select="@maxOccurs"/>
        <xsl:with-param name="parentPtr" select="$parentPtr"/>
        <xsl:with-param name="container" select="'element'"/>
        </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="ref" select="@ref" />
    <xsl:variable name="name" select="@name" />
    <xsl:variable name="type" select="@type" />
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Element'"/>
    </xsl:call-template>
    <xsl:choose>
      <xsl:when test="$ref">
        <xsl:value-of select="$level"/>
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$ref"/>
        </xsl:call-template>
        <xsl:variable name="locRef">
          <xsl:call-template name="strip-prefix">
            <xsl:with-param name="name" select="$ref"/>
          </xsl:call-template>
        </xsl:variable>
        <xsl:text> </xsl:text>
        <xsl:variable name="refType" select="/xs:schema/xs:element[@name=$locRef]/@type"/>
        <xsl:variable name="inlineType">
          <xsl:choose>
            <xsl:when test="$refType">
              <xsl:value-of select="$refType"/>
            </xsl:when>
            <xsl:otherwise>
              <xsl:value-of select="$locRef"/>
            </xsl:otherwise>
          </xsl:choose>
        </xsl:variable>
        <xsl:call-template name="convert-type">
          <xsl:with-param name="ptr" select="$ptr"/>
          <xsl:with-param name="type" select="$inlineType"/>
        </xsl:call-template>
        <xsl:text> `xml:"</xsl:text>
        <xsl:value-of select="$locRef"/>
        <xsl:if test="$omitempty = 'yes' and $ptr != ''">
          <xsl:text>,omitempty</xsl:text>
        </xsl:if>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"</xsl:text>
          <xsl:value-of select="$locRef"/>
          <xsl:if test="$omitempty = 'yes' and $ptr != ''">
            <xsl:text>,omitempty</xsl:text>
          </xsl:if>
        </xsl:if>
        <xsl:text>"`</xsl:text>
        <xsl:value-of select="$break"/>
      </xsl:when>
      <xsl:otherwise>
        <xsl:value-of select="$level"/>
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$name"/>
        </xsl:call-template>
        <xsl:text> </xsl:text>
        <xsl:call-template name="convert-type">
          <xsl:with-param name="ptr" select="$ptr"/>
          <xsl:with-param name="type" select="$type"/>
        </xsl:call-template>
        <xsl:text> `xml:"</xsl:text>
        <xsl:value-of select="$name"/>
        <xsl:if test="$omitempty = 'yes' and $ptr != ''">
          <xsl:text>,omitempty</xsl:text>
        </xsl:if>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"</xsl:text>
          <xsl:value-of select="$name"/>
          <xsl:if test="$omitempty = 'yes' and $ptr != ''">
            <xsl:text>,omitempty</xsl:text>
          </xsl:if>
        </xsl:if>
        <xsl:text>"`</xsl:text>
        <xsl:value-of select="$break"/>
      </xsl:otherwise>
    </xsl:choose>
    <xsl:apply-templates>
      <xsl:with-param name="name" select="$name"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:attribute">
    <xsl:param name="level" select="''"/>
    <xsl:variable name="ptr">
      <xsl:choose>
        <xsl:when test="not(@use) or @use = 'optional'">
          <xsl:text>*</xsl:text>
        </xsl:when>
       </xsl:choose>
    </xsl:variable>
    <xsl:variable name="name" select="@name" />
    <xsl:variable name="type" select="@type" />
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Attribute'"/>
    </xsl:call-template>
    <xsl:value-of select="$level"/>
    <xsl:call-template name="convert-name">
      <xsl:with-param name="name" select="$name"/>
    </xsl:call-template>
    <xsl:text> </xsl:text>
    <xsl:call-template name="convert-type">
      <xsl:with-param name="ptr" select="$ptr"/>
      <xsl:with-param name="type" select="'attribute'"/>
    </xsl:call-template>
    <xsl:text> `xml:"</xsl:text>
    <xsl:value-of select="$name"/>
    <xsl:text>,attr</xsl:text>
    <xsl:if test="$omitempty = 'yes' and $ptr = '*'">
      <xsl:text>,omitempty</xsl:text>
    </xsl:if>
    <xsl:if test="$json = 'yes'">
      <xsl:text>" json:"@</xsl:text>
      <xsl:value-of select="$name"/>
      <xsl:if test="$omitempty = 'yes' and $ptr = '*'">
        <xsl:text>,omitempty</xsl:text>
      </xsl:if>
    </xsl:if>
    <xsl:text>"`</xsl:text>
    <xsl:value-of select="$break"/>
  </xsl:template>

  <xsl:template match="xs:complexType" mode="global">
    <xsl:variable name="name" select="@name" />
    <xsl:call-template name="debug">
      <xsl:with-param name="text" select="'ComplexType'"/>
    </xsl:call-template>
    <xsl:text>type </xsl:text>
    <xsl:call-template name="convert-name">
      <xsl:with-param name="name" select="$name"/>
    </xsl:call-template>
    <xsl:text> struct {</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:apply-templates>
      <xsl:with-param name="name" select="$name"/>
      <xsl:with-param name="level" select="$indent" />
    </xsl:apply-templates>
    <xsl:text>}</xsl:text>
    <xsl:value-of select="$break"/>
    <xsl:value-of select="$break"/>
  </xsl:template>

  <xsl:template match="xs:complexType">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name"/>
    <xsl:param name="parentPtr"/>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'ComplexType'"/>
    </xsl:call-template>
    <xsl:value-of select="$level"/>
    <xsl:choose>
      <xsl:when test="$level != '' and $parentPtr != ''">
        <!--hoisted type-->
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$name"/>
        </xsl:call-template>
        <xsl:text> </xsl:text>
        <xsl:value-of select="$parentPtr"/>
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$name"/>
        </xsl:call-template>
        <xsl:text> `xml:"</xsl:text>
        <xsl:value-of select="$name"/>
        <xsl:if test="$omitempty = 'yes' and $parentPtr != ''">
          <xsl:text>,omitempty</xsl:text>
        </xsl:if>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"</xsl:text>
          <xsl:value-of select="$name"/>
          <xsl:if test="$omitempty = 'yes' and $parentPtr != ''">
            <xsl:text>,omitempty</xsl:text>
          </xsl:if>
        </xsl:if>
        <xsl:text>"`</xsl:text>
        <xsl:value-of select="$break"/>
      </xsl:when>
      <xsl:otherwise>
        <xsl:if test="$level = ''">
          <xsl:text>type </xsl:text>
        </xsl:if>
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$name"/>
        </xsl:call-template>
        <xsl:text> struct {</xsl:text>
        <xsl:value-of select="$break"/>
        <xsl:value-of select="$level"/>
        <xsl:value-of select="$indent"/>
        <xsl:text>XMLName xml.Name `xml:"</xsl:text>
        <xsl:value-of select="$name"/>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"-</xsl:text>
        </xsl:if>
        <xsl:text>"`</xsl:text>
        <xsl:value-of select="$break"/>
        <xsl:apply-templates>
          <xsl:with-param name="name" select="$name"/>
          <xsl:with-param name="level" select="concat($level,$indent)"/>
        </xsl:apply-templates>
        <xsl:value-of select="$level"/>
        <xsl:text>}</xsl:text>
        <xsl:if test="$level = ''">
          <xsl:value-of select="$break"/>
        </xsl:if>
        <xsl:value-of select="$break"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

  <xsl:template match="xs:simpleType" mode="global">
    <xsl:variable name="name" select="@name" />
    <xsl:call-template name="debug">
      <xsl:with-param name="text" select="'SimpleType'"/>
    </xsl:call-template>
    <xsl:apply-templates>
      <xsl:with-param name="name" select="@name"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:simpleType">
    <xsl:param name="level"/>
    <xsl:param name="parentPtr"/>
    <xsl:param name="name"/>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'SimpleType'"/>
    </xsl:call-template>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="name" select="$name"/>
      <xsl:with-param name="parentPtr" select="$parentPtr"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:simpleContent">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name"/>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'SimpleContent'"/>
    </xsl:call-template>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="name" select="$name"/>
      <xsl:with-param name="content" select="'simple'"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:complexContent">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name"/>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'ComplexContent'"/>
    </xsl:call-template>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="name" select="$name"/>
      <xsl:with-param name="content" select="'complex'"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:extension">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name" />
    <xsl:param name="content" />
    <xsl:variable name="normName">
      <xsl:call-template name="convert-name">
        <xsl:with-param name="name" select="$name" />
      </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="type" select="@base" />
    <xsl:variable name="normType">
      <xsl:call-template name="convert-type">
        <xsl:with-param name="type" select="$type" />
      </xsl:call-template>
    </xsl:variable>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Extension'"/>
    </xsl:call-template>
    <xsl:value-of select="$level"/>
    <xsl:if test="$content = 'simple'">
      <xsl:text>Text </xsl:text>
    </xsl:if>
    <!-- embed type for complex content-->
    <xsl:value-of select="$normType"/>
    <xsl:if test="$content = 'simple'">
      <xsl:text> `xml:",chardata</xsl:text>
          <xsl:if test="$json = 'yes'">
      <xsl:text>" json:"#text</xsl:text>
      <xsl:if test="$omitempty = 'yes'">
        <xsl:text>,omitempty</xsl:text>
      </xsl:if>
      </xsl:if>
      <xsl:text>"`</xsl:text>
    </xsl:if>
    <xsl:value-of select="$break"/>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="normName" select="$normName" />
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:sequence">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name" />
    <xsl:param name="parentPtr" />
    <xsl:variable name="ptr">
      <xsl:call-template name="get-cardinality">
        <xsl:with-param name="minOccurs" select="@minOccurs"/>
        <xsl:with-param name="maxOccurs" select="@maxOccurs"/>
        <xsl:with-param name="parentPtr" select="$parentPtr"/>
        <xsl:with-param name="container" select="'sequence'"/>
        </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="sequence">
      <xsl:apply-templates>
        <xsl:with-param name="name" select="$name"/>
        <xsl:with-param name="level" select="$level"/>
        <xsl:with-param name="parentPtr" select="$ptr"/>
      </xsl:apply-templates>
    </xsl:variable>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Sequence'"/>
    </xsl:call-template>
    <!-- nested choice may produce duplicate fields-->
    <xsl:call-template name="distinct-lines">
      <xsl:with-param name="text" select="$sequence"/>
    </xsl:call-template>
  </xsl:template>

  <xsl:template match="xs:choice">
    <xsl:param name="level" select="''"/>
    <xsl:param name="name" />
    <xsl:param name="parentPtr" />
    <xsl:variable name="ptr">
      <xsl:call-template name="get-cardinality">
        <xsl:with-param name="minOccurs" select="@minOccurs"/>
        <xsl:with-param name="maxOccurs" select="@maxOccurs"/>
        <xsl:with-param name="parentPtr" select="$parentPtr"/>
        <xsl:with-param name="container" select="'choice'"/>
        </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="choice">
      <xsl:apply-templates>
        <xsl:with-param name="level" select="$level"/>
        <xsl:with-param name="name" select="$name"/>
        <xsl:with-param name="parentPtr" select="$ptr"/>
      </xsl:apply-templates>
    </xsl:variable>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Choice'"/>
    </xsl:call-template>
    <!-- choice may produce duplicate fields-->
    <xsl:call-template name="distinct-lines">
      <xsl:with-param name="text" select="$choice"/>
    </xsl:call-template>
  </xsl:template>

  <xsl:template match="xs:restriction">
    <xsl:param name="level"/>
    <xsl:param name="parentPtr"/>
    <xsl:param name="name" />
    <xsl:variable name="normName">
      <xsl:call-template name="convert-name">
        <xsl:with-param name="name" select="$name" />
      </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="type" select="@base" />
    <xsl:variable name="normType">
      <xsl:call-template name="convert-type">
        <xsl:with-param name="type" select="$type" />
      </xsl:call-template>
    </xsl:variable>
    <xsl:call-template name="debug">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="'Restriction'"/>
    </xsl:call-template>
    <xsl:choose>
      <xsl:when test="$level = ''">
        <xsl:text>type </xsl:text>
        <xsl:value-of select="$normName"/>
        <xsl:text> </xsl:text>
        <xsl:value-of select="$normType"/>
        <xsl:value-of select="$break"/>
        <xsl:value-of select="$break"/>
        <xsl:apply-templates>
          <xsl:with-param name="level" select="$level"/>
          <xsl:with-param name="normName" select="$normName" />
        </xsl:apply-templates>
        <xsl:value-of select="$break"/>
      </xsl:when>
      <xsl:otherwise>
        <!--hoisted if enumeration-->
        <xsl:value-of select="$level"/>
        <xsl:value-of select="$normName"/>
        <xsl:text> </xsl:text>
        <xsl:choose>
          <xsl:when test="child::xs:enumeration">
            <xsl:value-of select="$normName"/>
          </xsl:when>
          <xsl:otherwise>
            <xsl:value-of select="$normType"/>
          </xsl:otherwise>
        </xsl:choose>
        <xsl:text> `xml:"</xsl:text>
        <xsl:value-of select="$name"/>
        <xsl:if test="$omitempty = 'yes' and $parentPtr != ''">
          <xsl:text>,omitempty</xsl:text>
        </xsl:if>
        <xsl:if test="$json = 'yes'">
          <xsl:text>" json:"</xsl:text>
          <xsl:value-of select="$name"/>
          <xsl:if test="$omitempty = 'yes' and $parentPtr != ''">
            <xsl:text>,omitempty</xsl:text>
          </xsl:if>
        </xsl:if>
        <xsl:text>"`</xsl:text>
        <xsl:value-of select="$break"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

  <xsl:template match="xs:any">
    <xsl:param name="name"/>
    <xsl:value-of select="$indent"/>
    <xsl:call-template name="convert-name">
      <xsl:with-param name="name" select="$name"/>
    </xsl:call-template>
    <xsl:text> []byte</xsl:text>
    <xsl:value-of select="$indent"/>
    <xsl:text> `xml:",innerxml"`</xsl:text>
    <xsl:value-of select="$break"/>
  </xsl:template>

  <xsl:template match="xs:enumeration">
    <xsl:param name="normName" />
    <xsl:variable name="constType">
      <xsl:value-of select="$normName"/>
      <xsl:call-template name="convert-name">
        <xsl:with-param name="name" select="@value"/>
      </xsl:call-template>
    </xsl:variable>
    <xsl:call-template name="debug">
      <xsl:with-param name="text" select="'Enumeration'"/>
    </xsl:call-template>
    <xsl:text>const </xsl:text>
    <xsl:value-of select="$constType"/>
    <xsl:text> </xsl:text>
    <xsl:value-of select="$normName"/>
    <xsl:text> = "</xsl:text>
    <xsl:value-of select="@value"/>
    <xsl:text>"</xsl:text>
    <xsl:value-of select="$break"/>
  </xsl:template>

  <xsl:template match="xs:annotation">
    <xsl:param name="level" select="''"/>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:annotation" mode="global">
    <xsl:param name="level" select="''"/>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template match="xs:documentation">
    <xsl:param name="level" select="''"/>
    <xsl:apply-templates>
      <xsl:with-param name="level" select="$level"/>
    </xsl:apply-templates>
  </xsl:template>

  <xsl:template name="debug">
    <xsl:param name="level"/>
    <xsl:param name="text"/>
    <xsl:if test="$debug = 'yes'">
      <xsl:value-of select="$level"/>
      <xsl:text>//</xsl:text>
      <xsl:value-of select="$text"/>
      <xsl:value-of select="$break"/>
    </xsl:if>
  </xsl:template>

  <xsl:template match="comment()">
    <xsl:param name="level"/>
    <xsl:call-template name="print-comment">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="."/>
    </xsl:call-template>
  </xsl:template>

  <xsl:template match="text()">
    <xsl:param name="level"/>
    <xsl:call-template name="print-comment">
      <xsl:with-param name="level" select="$level"/>
      <xsl:with-param name="text" select="."/>
    </xsl:call-template>
  </xsl:template>

  <xsl:template name="print-comment">
    <xsl:param name="level"/>
    <xsl:param name="text"/>
    <xsl:for-each select="str:tokenize($text, $break)">
      <xsl:value-of select="$level"/>
      <xsl:text>// </xsl:text>
      <xsl:value-of select="."/>
      <xsl:value-of select="$break"/>
    </xsl:for-each>
  </xsl:template>

  <xsl:template name="distinct-lines">
    <xsl:param name="text"/>
    <xsl:variable name="lines" select="str:tokenize($text,$break)"/>
    <xsl:for-each select="$lines[not(text() = preceding::text())]">
      <xsl:value-of select="."/>
      <xsl:value-of select="$break"/>
    </xsl:for-each>
  </xsl:template>

  <xsl:template name="strip-prefix">
    <xsl:param name="name"/>
    <xsl:choose>
      <xsl:when test="contains($name,':')">
        <xsl:value-of select="substring-after($name,':')"/>
      </xsl:when>
      <xsl:otherwise>
        <xsl:value-of select="$name"/>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

  <xsl:template name="get-cardinality">
    <xsl:param name="minOccurs"/>
    <xsl:param name="maxOccurs"/>
    <xsl:param name="parentPtr"/>
    <xsl:param name="container"/>
    <xsl:variable name="ptr"><!-- careful, the type is result-tree-fragment not a string -->
      <xsl:choose>
        <xsl:when test="number($maxOccurs) > 1 or $maxOccurs = 'unbounded'">
          <xsl:text>[]</xsl:text>
        </xsl:when>
        <xsl:when test="$container = 'choice'">
          <xsl:choose>
            <xsl:when test="$parentPtr = '[]'">
              <xsl:text>[]</xsl:text>
            </xsl:when>
            <xsl:otherwise>
              <xsl:text>*</xsl:text>
            </xsl:otherwise>
          </xsl:choose>
        </xsl:when>
        <xsl:when test="$container = 'sequence' or $container = 'element'">
          <xsl:choose>
            <xsl:when test="$minOccurs = '0'">
              <xsl:text>*</xsl:text>
            </xsl:when>
            <xsl:otherwise>
              <xsl:value-of select="$parentPtr"/>
            </xsl:otherwise>
          </xsl:choose>
        </xsl:when>
      </xsl:choose>
    </xsl:variable>
    <xsl:value-of select="string($ptr)"/>
  </xsl:template>

  <xsl:template name="get-ns">
    <xsl:param name="name"/>
    <xsl:variable name="prefix" select="substring-before($name,':')"/>
    <xsl:value-of select="namespace::node()[local-name() = $prefix]"/>
  </xsl:template>

  <xsl:template name="convert-name">
    <xsl:param name="name"/>
    <xsl:variable name="normName">
      <xsl:choose>
        <xsl:when test="contains($name,':')">
          <xsl:value-of select="substring-after($name,':')"/>
        </xsl:when>
        <xsl:otherwise>
          <xsl:value-of select="$name"/>
        </xsl:otherwise>
      </xsl:choose>
    </xsl:variable>
    <xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz'" />
    <xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'" />
    <xsl:for-each select="str:tokenize($normName, '_- ')">
      <xsl:value-of select="translate(substring(.,1,1),$lowercase,$uppercase)"/>
      <xsl:value-of select="substring(.,2)"/>
    </xsl:for-each>
  </xsl:template>

  <xsl:template name="convert-type">
    <xsl:param name="ptr"/>
    <xsl:param name="type"/>
    <xsl:variable name="array">
      <xsl:if test="$ptr = '[]'">
        <xsl:value-of select="$ptr"/>
      </xsl:if>
    </xsl:variable>
    <xsl:variable name="ns">
      <xsl:call-template name="get-ns">
        <xsl:with-param name="name" select="$type"/>
      </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="t">
      <xsl:choose>
        <xsl:when test="contains($type,':')">
          <xsl:value-of select="substring-after($type,':')"/>
        </xsl:when>
        <xsl:otherwise>
          <xsl:value-of select="$type"/>
        </xsl:otherwise>
      </xsl:choose>
    </xsl:variable>
    <xsl:choose>
      <xsl:when test="$ns = $xmlns_xsd">
        <xsl:choose>
          <xsl:when test="$t = 'string'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'language'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'dateTime'">
            <xsl:choose>
              <xsl:when test="$dateTimeType">
                <xsl:value-of select="$ptr"/>
                <xsl:value-of select="$dateTimeType"/>
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$array"/>
                <xsl:text>string</xsl:text>
              </xsl:otherwise>
            </xsl:choose>
          </xsl:when>
          <xsl:when test="$t = 'date'">
            <xsl:choose>
              <xsl:when test="$dateType">
                <xsl:value-of select="$ptr"/>
                <xsl:value-of select="$dateType"/>
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$array"/>
                <xsl:text>string</xsl:text>
              </xsl:otherwise>
            </xsl:choose>
          </xsl:when>
          <xsl:when test="$t = 'time'">
            <xsl:choose>
              <xsl:when test="$timeType">
                <xsl:value-of select="$ptr"/>
                <xsl:value-of select="$timeType"/>
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$array"/>
                <xsl:text>string</xsl:text>
              </xsl:otherwise>
            </xsl:choose>
          </xsl:when>
          <xsl:when test="$t = 'gYear'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'gYearMonth'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'gMonthDay'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'gDay'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'gMonth'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'base64Binary'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'normalizedString'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'token'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'NCName'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'NMTOKEN'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'NMTOKENS'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'anySimpleType'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'anyType'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'anyURI'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'int'">
            <xsl:value-of select="$array"/>
            <xsl:text>int</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'integer'">
            <xsl:value-of select="$array"/>
            <xsl:text>int64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'long'">
            <xsl:value-of select="$array"/>
            <xsl:text>int64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'negativeInteger'">
            <xsl:value-of select="$array"/>
            <xsl:text>int64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'nonNegativeInteger'">
            <xsl:value-of select="$array"/>
            <xsl:text>uint64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'positiveInteger'">
            <xsl:value-of select="$array"/>
            <xsl:text>uint64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'unsignedInt'">
            <xsl:value-of select="$array"/>
            <xsl:text>uint64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'byte'">
            <xsl:value-of select="$array"/>
            <xsl:text>int8</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'unsignedByte'">
            <xsl:value-of select="$array"/>
            <xsl:text>uint8</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'short'">
            <xsl:value-of select="$array"/>
            <xsl:text>int16</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'unsignedShort'">
            <xsl:value-of select="$array"/>
            <xsl:text>uint16</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'double'">
            <xsl:value-of select="$array"/>
            <xsl:text>float64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'decimal'">
            <xsl:choose>
              <xsl:when test="$decimalType">
                <xsl:value-of select="$ptr"/>
                <xsl:value-of select="$decimalType"/>
              </xsl:when>
              <xsl:otherwise>
                <xsl:value-of select="$array"/>
                <xsl:text>float64</xsl:text>
              </xsl:otherwise>
            </xsl:choose>
          </xsl:when>
          <xsl:when test="$t = 'float'">
            <xsl:value-of select="$array"/>
            <xsl:text>float64</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'boolean'">
            <xsl:value-of select="$array"/>
            <xsl:text>bool</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'ID'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:when test="$t = 'IDREF'">
            <xsl:value-of select="$array"/>
            <xsl:text>string</xsl:text>
          </xsl:when>
          <xsl:otherwise>
            <xsl:message>
              //unknown XSD type:<xsl:value-of select="$t"/>
            </xsl:message>
            <xsl:text>string</xsl:text>
          </xsl:otherwise>
        </xsl:choose>
      </xsl:when>
      <xsl:when test="$t = 'attribute'">
        <xsl:choose>
          <xsl:when test="$attributeForm = 'qualified'">
            <xsl:value-of select="$ptr"/>
            <xsl:value-of select="$qAttrType"/>
          </xsl:when>
          <xsl:otherwise>
            <xsl:text>string</xsl:text>
          </xsl:otherwise>
        </xsl:choose>
      </xsl:when>
      <xsl:otherwise>
        <xsl:value-of select="$ptr"/>
        <xsl:call-template name="convert-name">
          <xsl:with-param name="name" select="$t"></xsl:with-param>
        </xsl:call-template>
      </xsl:otherwise>
    </xsl:choose>
  </xsl:template>

</xsl:stylesheet>
