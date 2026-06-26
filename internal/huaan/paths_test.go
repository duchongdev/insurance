package huaan

import "testing"

func TestAPIPathsContainsCommonChannelAPIs(t *testing.T) {
	want := []string{
		ProductInfoPath,
		SmsSendPath,
		SmsValidPath,
		SmsNoValidPath,
		PriceByUserPath,
		PolicyByPhonePath,
		UserInfoByPhoneNoPath,
		LiabilitiesByProductIDPath,
		GetPhoneByTokenPath,
	}
	for _, path := range want {
		found := false
		for _, p := range APIPaths {
			if p == path {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("APIPaths missing %s", path)
		}
	}
}
