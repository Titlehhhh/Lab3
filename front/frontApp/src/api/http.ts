const API_BASE = '/api'

type ErrorPayload = {
  error?: string
}

const resolveErrorMessage = async (response: Response) => {
  try {
    const data = (await response.json()) as ErrorPayload
    if (data?.error) {
      return data.error
    }
  } catch {
    // ignore parsing errors
  }

  if (response.status === 401) {
    return 'Требуется авторизация.'
  }
  if (response.status === 403) {
    return 'Недостаточно прав.'
  }
  if (response.status === 404) {
    return 'Ресурс не найден.'
  }

  return 'Произошла ошибка запроса.'
}

export const request = async <T>(path: string, options: RequestInit = {}): Promise<T> => {
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: 'include',
    ...options,
  })

  if (!response.ok) {
    throw new Error(await resolveErrorMessage(response))
  }

  if (response.status === 204) {
    return undefined as T
  }

  const contentType = response.headers.get('content-type')
  if (contentType?.includes('application/json')) {
    return (await response.json()) as T
  }

  return (await response.text()) as unknown as T
}
