package huaan

import "strings"

// resolveUpstreamPath 解析华安 URL 路径；部分环境 legacy /upChannelApi 须走 /proxy 前缀。
func (c *Client) resolveUpstreamPath(channelPath string) string {
	p := ResolveUpstreamPath(c.cfg.HuaAn.APIPath, channelPath)
	if !c.useProxyUpstreamPrefix() {
		return p
	}
	if strings.HasPrefix(p, "/common/channel/api/") || strings.HasPrefix(p, "/proxy/") {
		return p
	}
	if strings.HasPrefix(p, c.cfg.HuaAn.APIPath+"/") {
		suffix := strings.TrimPrefix(p, c.cfg.HuaAn.APIPath)
		return "/proxy" + c.cfg.HuaAn.APIPath + suffix
	}
	return p
}

func (c *Client) useProxyUpstreamPrefix() bool {
	host := strings.ToLower(c.cfg.HuaAn.BaseURL)
	return strings.Contains(host, "hahealth.ink") || strings.Contains(host, "47.97.156.18")
}
