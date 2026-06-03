package huaan

import "testing"

func TestParseBankPayPairs(t *testing.T) {
	raw := []byte(`{
		"code": 200,
		"message": "ok",
		"data": [
			{"bankCode": "BOC", "bankName": "中国银行", "debitCard": "1", "payChannelId": "aaa"},
			{"bankCode": "CCB", "bankName": "建设银行", "debitCard": "1", "creditCard": "1", "payChannelId": "bbb"}
		]
	}`)
	pairs, err := ParseBankPayPairs(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(pairs) != 2 || pairs[0].BankCode != "BOC" || pairs[0].PayChannelID != "aaa" {
		t.Fatalf("unexpected pairs: %+v", pairs)
	}

	got, err := SelectBankPayPair(pairs, "CCB")
	if err != nil || got.PayChannelID != "bbb" {
		t.Fatalf("select CCB: %+v err=%v", got, err)
	}
	if pairs[0].DefaultCardType() != "1" {
		t.Fatalf("cardType=%s", pairs[0].DefaultCardType())
	}
}

func TestSelectBankPayPair_notFound(t *testing.T) {
	pairs := []BankPayPair{{BankCode: "BOC", PayChannelID: "x"}}
	if _, err := SelectBankPayPair(pairs, "XYZ"); err == nil {
		t.Fatal("expected error")
	}
}
