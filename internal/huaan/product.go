package huaan

import (
	"context"
	"encoding/json"
	"fmt"
)

const productInfoByChannelPath = "/getProductInfoByChannel"

// ChannelProduct getProductInfoByChannel 响应中的可售产品。
type ChannelProduct struct {
	ProductCode string
	ProductName string
}

// ParseChannelProducts 从 getProductInfoByChannel 华安响应 JSON 解析产品列表。
func ParseChannelProducts(respBytes []byte) ([]ChannelProduct, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if ParseCode(respBytes) != 200 {
		return nil, fmt.Errorf("getProductInfoByChannel code=%d", ParseCode(respBytes))
	}
	items, ok := resp["data"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("getProductInfoByChannel data empty or invalid")
	}
	out := make([]ChannelProduct, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		p := ChannelProduct{
			ProductCode: strVal(m, "productCode", "product_code"),
			ProductName: strVal(m, "productName", "product_name"),
		}
		if p.ProductCode == "" {
			continue
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no productCode in getProductInfoByChannel data")
	}
	return out, nil
}

// SelectChannelProduct 按 productCode 选取产品；productCode 为空时取列表第一项。
func SelectChannelProduct(products []ChannelProduct, productCode string) (ChannelProduct, error) {
	if len(products) == 0 {
		return ChannelProduct{}, fmt.Errorf("empty product list")
	}
	if productCode == "" {
		return products[0], nil
	}
	for _, p := range products {
		if p.ProductCode == productCode {
			return p, nil
		}
	}
	return ChannelProduct{}, fmt.Errorf("productCode %q not found in channel products (%d items)", productCode, len(products))
}

// FetchChannelProducts 调用 getProductInfoByChannel 并解析可售产品列表。
func FetchChannelProducts(ctx context.Context, client *Client, channelCode string) ([]ChannelProduct, *CallResult, error) {
	body := BuildRequestBody(channelCode, nil)
	result, err := client.Call(ctx, productInfoByChannelPath, body)
	if err != nil {
		return nil, nil, err
	}
	products, err := ParseChannelProducts(result.Body)
	if err != nil {
		return nil, result, err
	}
	return products, result, nil
}
