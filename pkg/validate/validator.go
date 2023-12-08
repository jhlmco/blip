package validate

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/Masterminds/semver"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed bom-1.5.schema.json
var cdx_1_5_schema []byte

func Validate(bom string) error {
	schemaLoader := gojsonschema.NewBytesLoader(cdx_1_5_schema)

	documentLoader := gojsonschema.NewReferenceLoader("file://" + bom)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid cyclondx json\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	bomBytes, err := os.ReadFile(bom)
	if err != nil {
		panic(err)
	}
	// Decode the BOM
	sbom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(bytes.NewReader(bomBytes), cdx.BOMFileFormatJSON)
	if err = decoder.Decode(sbom); err != nil {
		panic(err)
	}

	c, _ := semver.NewConstraint(">= 1.5")

	v, _ := semver.NewVersion(sbom.SpecVersion.String())

	if !c.Check(v) {
		fmt.Printf("SBOM Specification does not meet minimum requirements.  Version is %s, should be >= 1.5\n.", v.Original())
	}

	return nil
}
