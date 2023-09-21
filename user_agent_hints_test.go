package ch

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBrand_Success(t *testing.T) {
	type testData struct {
		name    string
		version string
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUA, `"Opera";v="81", "Chromium";v="96", " Not A;Brand";v="99"`)

	h2 := http.Header{}
	h2.Add(SecCHUA, `Opera;v=82, Chromium;v=96`)

	h3 := http.Header{}
	h3.Add(SecCHUA, `"Microsoft Edge";v="33",Chromium;v=96`)

	h4 := http.Header{}
	h4.Add(SecCHUA, `"Google Chrome";v="104",Chromium;v=96`)

	h5 := http.Header{}
	h5.Add(SecCHUA, `"(Not(A:Brand";v="8","Chromium";v="96"`)

	tests := map[string]testData{
		"Opera":          {"Opera", "81", h1},
		"Opera 2":        {"Opera", "82", h2},
		"Microsoft Edge": {"Microsoft Edge", "33", h3},
		"Google Chrome":  {"Google Chrome", "104", h4},
		"Chromium":       {"Chromium", "96", h5},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetBrand(val.headers)

			assert.NotNil(t, got)
			assert.Equal(t, val.name, got.Name)
			assert.Equal(t, val.version, got.Version)
		})
	}
}

func TestGetBrand_NotFound(t *testing.T) {
	h1 := http.Header{}
	h1.Add(SecCHUA, `"(Not(A:Brand";v="8"`)
	h2 := http.Header{}
	h2.Add(SecCHUA, `"(Not(A:Brand";v="8"`)

	tests := map[string]http.Header{
		"No Sec-CH-UA header":          {},
		"Brand in Sec-CH-UA header":    h1,
		"No brand in Sec-CH-UA header": h2,
	}

	for name, headers := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Nil(t, GetBrand(headers))
		})
	}
}

func TestGetArch(t *testing.T) {
	type testData struct {
		arch    string
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUAArch, "x86")

	h2 := http.Header{}
	h2.Add(SecCHUAArch, "ARM")

	tests := map[string]testData{
		"x86":               {"x86", h1},
		"ARM":               {"ARM", h2},
		"No Arch in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.arch, GetArch(val.headers))
		})
	}
}

func TestGetBitness(t *testing.T) {
	type testData struct {
		bitness string
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUABitness, "64")

	h2 := http.Header{}
	h2.Add(SecCHUABitness, "32")

	tests := map[string]testData{
		"32":                   {"32", h2},
		"64":                   {"64", h1},
		"No Bitness in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.bitness, GetBitness(val.headers))
		})
	}
}

func TestIsMobile(t *testing.T) {
	type testData struct {
		val     bool
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUAMobile, "?1")

	h2 := http.Header{}
	h2.Add(SecCHUAMobile, "?0")

	tests := map[string]testData{
		"mobile":               {true, h1},
		"not mobile":           {false, h2},
		"No Bitness in header": {false, http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.val, IsMobile(val.headers))
		})
	}
}

func TestGetColorScheme(t *testing.T) {
	type testData struct {
		colorScheme string
		headers     http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHPrefersColorScheme, "light")

	h2 := http.Header{}
	h2.Add(SecCHPrefersColorScheme, "dark")

	tests := map[string]testData{
		"light":                     {"light", h1},
		"dark":                      {"dark", h2},
		"No Color Scheme in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.colorScheme, GetColorScheme(val.headers))
		})
	}
}

func TestGetPlatformVersion(t *testing.T) {
	type testData struct {
		version string
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUAPlatformVersion, "11.0.0")

	tests := map[string]testData{
		"11.0.0":               {"11.0.0", h1},
		"No Version in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.version, GetPlatformVersion(val.headers))
		})
	}
}

func TestGetPlatform(t *testing.T) {
	type testData struct {
		platform string
		headers  http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUAPlatform, "Android")

	h2 := http.Header{}
	h2.Add(SecCHUAPlatform, "Windows")

	h3 := http.Header{}
	h3.Add(SecCHUAPlatform, "iOS")

	tests := map[string]testData{
		"Android":               {"Android", h1},
		"Windows":               {"Windows", h2},
		"iOS":                   {"iOS", h3},
		"No Platform in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.platform, GetPlatform(val.headers))
		})
	}
}

func TestGetModel(t *testing.T) {
	type testData struct {
		model   string
		headers http.Header
	}

	h1 := http.Header{}
	h1.Add(SecCHUAModel, "Pixel 3 XL")

	tests := map[string]testData{
		"Pixel 3 XL":            {"Pixel 3 XL", h1},
		"No Platform in header": {"", http.Header{}},
	}

	for name, val := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, val.model, GetModel(val.headers))
		})
	}
}
