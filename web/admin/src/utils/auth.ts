export const TOKEN_KEY = 'admin_token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

export function maskSecret(value: string | undefined | null, visible = 4): string {
  if (!value) return '***'
  if (value.length <= visible * 2) return '***'
  return `${value.slice(0, visible)}***${value.slice(-visible)}`
}
