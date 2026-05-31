package sign

import "testing"

func TestBuild_ExampleFromDoc(t *testing.T) {
	params := map[string]interface{}{
		"phoneNo":            "13968526776",
		"signSerialNumber":   "368dcfac0f534e2ea25ee3763de7b16b",
		"isSignType":         "0",
		"key":                "c7ae4fc06ca25a73b96fbe2d199e1819",
		"timestamp":          "1713236726003",
	}
	got := Build(params, "")
	want := "66aa23b057593150546a42d7e51d0b7e"
	if got != want {
		t.Fatalf("sign=%s want=%s", got, want)
	}
}

func TestVerify(t *testing.T) {
	params := map[string]interface{}{
		"phoneNo":   "13968526776",
		"key":       "secret",
		"timestamp": "1713236726003",
	}
	params["sign"] = Build(params, "")
	if !Verify(params, "secret") {
		t.Fatal("verify should pass")
	}
}
