package hclparser

import (
	"errors"
	"sync"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

var (
	once   sync.Once
	parser *hclparse.Parser
)

func GetBlocks(filepath string) (hclsyntax.Blocks, error) {
	p := getParser() // get the singleton instance of the parser
	file, diags := p.ParseHCLFile(filepath)
	if diags.HasErrors() {
		return nil, errors.New("Error parsing HCL file: " + diags.Error())
	}

	rootBody := file.Body.(*hclsyntax.Body)
	return rootBody.Blocks, nil
}

func GetBlockAttributes(block *hclsyntax.Block, m map[string]cty.Value) error {
	for name, attr := range block.Body.Attributes {
		val, diag := attr.Expr.Value(nil)
		if diag.HasErrors() {
			return errors.New("Failed to evaluate " + name + ": " + diag.Error())
		}
		m[name] = val
	}

	return nil
}

func getParser() *hclparse.Parser {
	once.Do(func() {
		parser = hclparse.NewParser()
	})

	return parser
}
