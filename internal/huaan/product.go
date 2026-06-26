package huaan

import (
	"context"
	"encoding/json"
	"fmt"
)

// ChannelProduct 产品列表项（product/info 等接口响应中的 productCode/productName）。
type ChannelProduct struct {
	ProductCode string
	ProductName string
}

// ParseChannelProducts 从华安产品列表响应 JSON 解析（data 为数组，含 productCode/productName）。
// 用于 product/info；getProductInfoByChannel 已废弃，响应结构相同但不应再调用。
func ParseChannelProducts(respBytes []byte) ([]ChannelProduct, error) {
	var resp map[string]interface{}
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if ParseCode(respBytes) != 200 {
		return nil, fmt.Errorf("product list code=%d", ParseCode(respBytes))
	}
	items, ok := resp["data"].([]interface{})
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("product list data empty or invalid")
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
		return nil, fmt.Errorf("no productCode in product list data")
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
	return ChannelProduct{}, fmt.Errorf("productCode %q not found in product list (%d items)", productCode, len(products))
}

// FetchProductInfoProducts 调用 product/info 并解析可售产品列表（须传三要素等业务字段）。
func FetchProductInfoProducts(ctx context.Context, client *Client, channelCode string, fields map[string]interface{}) ([]ChannelProduct, *CallResult, error) {
	body := BuildRequestBody(channelCode, fields)
	result, err := client.Call(ctx, ProductInfoPath, body)
	if err != nil {
		return nil, nil, err
	}
	products, err := ParseChannelProducts(result.Body)
	if err != nil {
		return nil, result, err
	}
	return products, result, nil
}

// FetchChannelProducts 已废弃：华安 getProductInfoByChannel 不再使用，请改用 FetchProductInfoProducts。
//
// Deprecated: Use FetchProductInfoProducts with product/info instead.
func FetchChannelProducts(ctx context.Context, client *Client, channelCode string) ([]ChannelProduct, *CallResult, error) {
	body := BuildRequestBody(channelCode, nil)
	result, err := client.Call(ctx, ProductInfoByChannelPath, body)
	if err != nil {
		return nil, nil, err
	}
	products, err := ParseChannelProducts(result.Body)
	if err != nil {
		return nil, result, err
	}
	return products, result, nil
}
