import request from './request'

export interface PageResult<T> {
  list: T[]
  total: number
}

export interface Stats {
  policyCount: number
  userCount: number
  signCount: number
}

export interface Channel {
  id: number
  channelCode: string
  channelName?: string
  channelKey: string
  huaAnKey: string
  status: number
  remark?: string
  createdAt?: string
  updatedAt?: string
}

export interface PolicyRecord {
  id: number
  channelCode: string
  policyId: string
  policyNo?: string
  userId?: string
  productCode?: string
  policyStatus?: string
  createdAt?: string
}

export interface UserRecord {
  id: number
  channelCode: string
  userId: string
  createdAt?: string
}

export interface SignRecord {
  id: number
  channelCode: string
  signId?: string
  policyId?: string
  bankCode?: string
  createdAt?: string
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

export function fetchStats(channelCode?: string) {
  return request.get<Stats>('/stats', { params: { channelCode: channelCode || undefined } })
}

export function fetchChannels(page: number, size: number) {
  return request.get<PageResult<Channel>>('/channels', { params: { page, size } })
}

export function createChannel(data: Partial<Channel>) {
  return request.post<Channel>('/channels', data)
}

export function deleteChannel(id: number) {
  return request.delete(`/channels/${id}`)
}

export function fetchPolicies(page: number, size: number, channelCode?: string) {
  return request.get<PageResult<PolicyRecord>>('/policies', {
    params: { page, size, channelCode: channelCode || undefined },
  })
}

export function fetchUsers(page: number, size: number, channelCode?: string) {
  return request.get<PageResult<UserRecord>>('/users', {
    params: { page, size, channelCode: channelCode || undefined },
  })
}

export function fetchSigns(page: number, size: number, channelCode?: string) {
  return request.get<PageResult<SignRecord>>('/signs', {
    params: { page, size, channelCode: channelCode || undefined },
  })
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

export function refreshBankList(channelCode: string, huaAnKey: string) {
  return request.post<string>('/bank-list/refresh', { channelCode, huaAnKey }, {
    transformResponse: [(data) => data],
    responseType: 'text',
  })
}
