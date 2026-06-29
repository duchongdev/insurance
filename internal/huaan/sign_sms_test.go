package huaan

import (
	"testing"

	"github.com/huaan/insurance-bridge/internal/pkg/sign"
)

func TestSmsSignExcludeKeys(t *testing.T) {
	key := "fb9ec7236b6c45b7bfd562672e0373ea"
	base := map[string]interface{}{
		"channelCode": "BLtJjF",
		"timestamp":   "1782469175918",
	}

	tests := []struct {
		path     string
		body     map[string]interface{}
		wantSign string
	}{
		{
			path: SmsNoValidPath,
			body: mergeMaps(base, map[string]interface{}{"mobile": "13811045503"}),
		},
		{
			path: SmsValidPath,
			body: mergeMaps(base, map[string]interface{}{"phoneNo": "13811045503", "code": "1234"}),
			wantSign: sign.BuildHuaAn(mergeMaps(base, map[string]interface{}{
				"phoneNo": "13811045503",
				"code":    "1234",
			}), key),
		},
		{
			path: SmsSendPath,
			body: mergeMaps(base, map[string]interface{}{
				"idCard":  "13068319940517031X",
				"name":    "杜冲",
				"phoneNo": "13811045503",
			}),
			wantSign: sign.BuildHuaAn(mergeMaps(base, map[string]interface{}{
				"phoneNo": "13811045503",
			}), key),
		},
	}
	for _, tc := range tests {
		got := buildHuaAnSign(tc.path, tc.body, key)
		want := tc.wantSign
		if want == "" {
			want = sign.BuildHuaAn(base, key)
		}
		if got != want {
			t.Fatalf("%s sign=%q want=%q", tc.path, got, want)
		}
	}
}

func mergeMaps(base, extra map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
