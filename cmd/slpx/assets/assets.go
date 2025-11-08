package assets

import "embed"

//go:embed default/* advanced/*
var embeddedAssets embed.FS

func LoadDefaultInitFile() string {
	content, err := embeddedAssets.ReadFile("default/init.slpx")
	if err != nil {
		panic(err)
	}
	return string(content)
}

func LoadDefaultThemesFile() string {
	content, err := embeddedAssets.ReadFile("default/themes.slpx")
	if err != nil {
		panic(err)
	}
	return string(content)
}

func LoadDefaultVariant() map[string]string {
	return map[string]string{
		"init.slpx":   mustRead("default/init.slpx"),
		"themes.slpx": mustRead("default/themes.slpx"),
	}
}

func LoadAdvancedVariant() map[string]string {
	return map[string]string{
		"init.slpx":     mustRead("advanced/init.slpx"),
		"themes.slpx":   mustRead("advanced/themes.slpx"),
		"commands.slpx": mustRead("advanced/commands.slpx"),
		"preload.slpx":  mustRead("advanced/preload.slpx"),
		"handler.slpx":  mustRead("advanced/handler.slpx"),
	}
}

func mustRead(path string) string {
	content, err := embeddedAssets.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

