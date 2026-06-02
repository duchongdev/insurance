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
  status: number
  remark?: string
  createdAt?: string
  updatedAt?: string
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

export function createChannel(data: Partial<Channel>) {
  return request.post<Channel>('/channels', data)
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

export function refreshBankList(channelCode: string, huaAnKey: string) {
  return request.post<string>('/bank-list/refresh', { channelCode, huaAnKey }, {
    transformResponse: [(data) => data],
    responseType: 'text',
  })
}
