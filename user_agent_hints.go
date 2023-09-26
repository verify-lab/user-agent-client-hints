package ch

import (
	"net/http"
	"strings"

	"github.com/verify-lab/strutil"
)

const (
	// Response headers
	CriticalCH = "Critical-CH"
	AcceptCH   = "Accept-CH"
	Vary       = "Vary"

	// Request headers
	SecCHUA                 = "Sec-CH-UA"
	SecCHUAArch             = "Sec-CH-UA-Arch"
	SecCHUABitness          = "Sec-CH-UA-Bitness"
	SecCHUAMobile           = "Sec-CH-UA-Mobile"
	SecCHUAModel            = "Sec-CH-UA-Model"
	SecCHUAPlatform         = "Sec-CH-UA-Platform"
	SecCHUAPlatformVersion  = "Sec-CH-UA-Platform-Version"
	SecCHPrefersColorScheme = "Sec-CH-Prefers-Color-Scheme"
)

const headerBoolTrue = "?1"

type Brand struct {
	Name    string
	Version string
}

// GetBrand provides the user-agent's branding and significant version information.
func GetBrand(headers http.Header) *Brand {
	ua := headers.Get(SecCHUA)
	data := strings.Split(ua, ",")

	for _, val := range data {
		val = strings.TrimSpace(val)

		if val == "" {
			continue
		}

		if strings.Contains(val, "Not") || strings.Contains(val, "not") {
			continue
		}

		if strings.Contains(val, "Brand") || strings.Contains(val, "brand") {
			continue
		}

		if strings.Contains(val, "Chromium") || strings.Contains(val, "chromium") {
			// Sec-CH-UA: "(Not(A:Brand";v="8", "Chromium";v="98"
			if len(data) > 2 {
				continue
			}
		}

		return headerValueToBrand(val)
	}

	return nil
}

func headerValueToBrand(val string) *Brand {
	data := strings.Split(val, ";v=")
	data[0] = strings.Trim(data[0], `"`)
	data[0] = strutil.StripNonPrintable(data[0])
	data[0] = strings.TrimSpace(data[0])

	data[1] = strings.Trim(data[1], `"`)
	data[1] = strutil.StripNonPrintable(data[1])
	data[1] = strings.TrimSpace(data[1])

	return &Brand{
		Name:    data[0],
		Version: data[1],
	}
}

// GetArch provides the user-agent's underlying CPU architecture, such as ARM or x86.
func GetArch(headers http.Header) string {
	return strutil.StripNonPrintable(headers.Get(SecCHUAArch))
}

// GetBitness provides the "bitness" of the user-agent's underlying CPU architecture.
// This is the size in bits of an integer or memory address—typically 64 or 32 bits.
func GetBitness(headers http.Header) string {
	return strutil.StripNonPrintable(headers.Get(SecCHUABitness))
}

// IsMobile indicates whether the browser is on a mobile device.
func IsMobile(headers http.Header) bool {
	return headerBoolTrue == headers.Get(SecCHUAMobile)
}

// GetModel provides the device model on which the browser is running.
// For example "Pixel 3".
func GetModel(headers http.Header) string {
	return strings.Trim(strutil.StripNonPrintable(headers.Get(SecCHUAModel)), `"`)
}

// GetPlatform provides the platform or operating system on which the user agent is running.
// One of the following strings: "Android", "Chrome OS", "Chromium OS", "iOS", "Linux", "macOS", "Windows", or "Unknown"
func GetPlatform(headers http.Header) string {
	return strutil.StripNonPrintable(headers.Get(SecCHUAPlatform))
}

// GetPlatformVersion provides the version of the operating system on which the user agent is running.
// The version string typically contains the operating system version in a string, consisting of dot-separated major, minor and patch version numbers.
// For example, "11.0.0"
// The version string on Linux is always empty.
func GetPlatformVersion(headers http.Header) string {
	return strutil.StripNonPrintable(headers.Get(SecCHUAPlatformVersion))
}

// GetColorScheme provides user color scheme preference at request time
// A string indicating the user agent's preference for dark or light content: "light" or "dark".
func GetColorScheme(headers http.Header) string {
	return strutil.StripNonPrintable(headers.Get(SecCHPrefersColorScheme))
}
