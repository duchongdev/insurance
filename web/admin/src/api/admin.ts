import request from './request'

export interface PageResult<T> {
  list: T[]
  total: number
}

export interface Channel {
  id: number
  channelCode: string
  channelName?: string
  channelKey: string
  huaAnKey: string
  huaanSettingId?: number
  envType?: HuaAnEnvType
  callbackUrl: string
  status: number
  remark?: string
  createdAt?: string
  updatedAt?: string
}

export type HuaAnEnvType = 'prod' | 'test'

export interface HuaAnConfig {
  id?: number
  name?: string
  envType: HuaAnEnvType
  baseUrl: string
  channelCode: string
  channelSecret: string
  isActive?: boolean
  createdAt?: string
  updatedAt?: string
}

export interface HuaAnConnectivityResult {
  ok: boolean
  durationMs: number
  upstreamUrl: string
  httpStatus?: number
  message: string
}

export interface BankRecord {
  id: number
  bankCode: string
  bankName: string
  debitCard: number
  creditCard: number
  status: number
  createTime?: string
  updateTime?: string
}

export function login(username: string, password: string) {
  return request.post<{ token: string }>('/login', { username, password })
}

export function fetchChannels(page: number, size: number) {
  return request.get<PageResult<Channel>>('/channels', { params: { page, size } })
}

export function createChannel(data: {
  huaanSettingId: number
  channelName?: string
  channelKey?: string
  callbackUrl: string
  status?: number
}) {
  return request.post<Channel>('/channels', data)
}

export function updateChannel(id: number, data: Partial<Channel>) {
  return request.put<Channel>(`/channels/${id}`, data)
}

export function deleteChannel(id: number) {
  return request.delete(`/channels/${id}`)
}

export function fetchBanks(page: number, size: number, status?: number) {
  return request.get<PageResult<BankRecord>>('/banks', {
    params: { page, size, status: status || undefined },
  })
}

export function getBankListCache(channelCode: string) {
  return request.get<string>('/bank-list', {
    params: { channelCode },
    transformResponse: [(data) => data],
    responseType: 'text',
  })
}

export function refreshBankList(huaanSettingId: number) {
  return request.post<string>('/bank-list/refresh', { huaanSettingId }, {
    transformResponse: [(data) => data],
    responseType: 'text',
  })
}

export function fetchHuaAnConfigs() {
  return request.get<{ list: HuaAnConfig[] }>('/huaan-configs')
}

export function createHuaAnConfig(data: Partial<HuaAnConfig>) {
  return request.post<HuaAnConfig>('/huaan-configs', data)
}

export function updateHuaAnConfig(id: number, data: Partial<HuaAnConfig>) {
  return request.put<HuaAnConfig>(`/huaan-configs/${id}`, data)
}

export function deleteHuaAnConfig(id: number) {
  return request.delete(`/huaan-configs/${id}`)
}

export function testHuaAnConnectivity(id: number) {
  return request.post<HuaAnConnectivityResult>(`/huaan-configs/${id}/test`)
}
