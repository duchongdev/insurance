// Package callback 定义华安调用本服务的回调路径（非渠道 upChannelApi 前缀）。
package callback

// InsureNotifyPath 投保结果回调：华安 POST 本服务，本服务再 POST 渠道 callbackUrl。
const InsureNotifyPath = "/huaan/callback/insureNotify"
