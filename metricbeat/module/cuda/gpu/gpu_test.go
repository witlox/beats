package gpu

import (
	"encoding/xml"
	"testing"

	rice "github.com/GeertJohan/go.rice"
	"github.com/stretchr/testify/assert"
)

var (
	box        = rice.MustFindBox("_meta/test")
	exampleXML = box.MustBytes("example.xml")
)

func TestUnmarshal(t *testing.T) {
	var info NvidiaSmi
	err := xml.Unmarshal(exampleXML, &info)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
}

