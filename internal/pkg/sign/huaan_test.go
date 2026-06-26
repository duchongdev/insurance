package sign

import (
	"strings"
	"testing"
)

func TestBuildHuaAn_skipsEmptyAndKeyInBody(t *testing.T) {
	got := BuildHuaAn(map[string]interface{}{
		"channelCode": "IF8XE0",
		"timestamp":   "1780701234123",
		"extra":       "",
		"key":         "should-not-appear-in-middle",
		"sign":        "ignored",
	}, "97dce65abf194b87b5d5b661f22c7366")
	if len(got) != 32 {
		t.Fatalf("len=%d", len(got))
	}
	if got != strings.ToUpper(got) {
		t.Fatalf("sign not uppercase: %s", got)
	}
	if !VerifyHuaAnHeader(map[string]interface{}{
		"channelCode": "IF8XE0",
		"timestamp":   "1780701234123",
	}, "97dce65abf194b87b5d5b661f22c7366", got) {
		t.Fatal("verify failed")
	}
}

func TestBuildHuaAn_nestedObjectSorted(t *testing.T) {
	a := BuildHuaAn(map[string]interface{}{
		"channelCode": "CH1",
		"timestamp":   "1",
		"productPriceList": []interface{}{
			map[string]interface{}{
				"price":         0.0,
				"kindCode":      "102",
				"liabilityName": "罕见保险金",
				"uwCount":       10,
			},
		},
	}, "secret")
	b := BuildHuaAn(map[string]interface{}{
		"channelCode": "CH1",
		"timestamp":   "1",
		"productPriceList": []interface{}{
			map[string]interface{}{
				"liabilityName": "罕见保险金",
				"kindCode":      "102",
				"uwCount":       10,
				"price":         0.0,
			},
		},
	}, "secret")
	if a != b {
		t.Fatalf("nested key order changed sign: %s vs %s", a, b)
	}
}

func TestBuildHuaAn_plaintextOrder(t *testing.T) {
	plain := buildHuaAnPlaintext(map[string]interface{}{
		"channelCode": "IF8XE0",
		"timestamp":   "1780701234123",
	}, "97dce65abf194b87b5d5b661f22c7366")
	want := "channelCode=IF8XE0&timestamp=1780701234123&key=97dce65abf194b87b5d5b661f22c7366"
	if plain != want {
		t.Fatalf("plain=%q want=%q", plain, want)
	}
}
