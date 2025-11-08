package assets

import "embed"

//go:embed init.slpx themes.slpx
var embeddedAssets embed.FS

func LoadDefaultInitFile() string {

	content, err := embeddedAssets.ReadFile("init.slpx")
	if err != nil {
		panic(err)
	}

	return string(content)
}

func LoadDefaultThemesFile() string {

	content, err := embeddedAssets.ReadFile("themes.slpx")
	if err != nil {
		panic(err)
	}

	return string(content)
}
