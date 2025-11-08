package assets

import "embed"

//go:embed init.slpx
var embeddedAssets embed.FS

func LoadDefaultInitFile() string {

	content, err := embeddedAssets.ReadFile("init.slpx")
	if err != nil {
		panic(err)
	}

	return string(content)
}
