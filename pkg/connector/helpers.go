package connector

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ResourcesPageSize = 50

var titleCaser = cases.Title(language.English)
