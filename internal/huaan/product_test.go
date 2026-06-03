package huaan

import "testing"

func TestParseChannelProducts(t *testing.T) {
	raw := []byte(`{
		"code": 200,
		"message": "操作成功",
		"data": [
			{"productCode": "ZFHLW1041001", "productName": "百万医疗险-体验版"},
			{"productCode": "ZFHLW1040003", "productName": "抗癌险-体验版"}
		]
	}`)
	products, err := ParseChannelProducts(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 2 || products[0].ProductCode != "ZFHLW1041001" {
		t.Fatalf("unexpected: %+v", products)
	}
	got, err := SelectChannelProduct(products, "ZFHLW1040003")
	if err != nil || got.ProductName != "抗癌险-体验版" {
		t.Fatalf("select: %+v err=%v", got, err)
	}
}
