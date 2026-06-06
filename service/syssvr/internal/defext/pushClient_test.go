package defext

import "testing"

func TestNormalizePushClientPlatform(t *testing.T) {
	cases := map[string]string{
		"harmonyos": "harmony",
		"HarmonyOS": "harmony",
		"ohos":      "harmony",
		"harmony":   "harmony",
		"android":   "android",
		"ios":       "ios",
		"unbind":    "unbind",
	}
	for in, want := range cases {
		if got := NormalizePushClientPlatform(in); got != want {
			t.Fatalf("%q => %q, want %q", in, got, want)
		}
	}
}
