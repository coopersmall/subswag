type SupportedMethod = 'DELETE' | 'GET' | 'PATCH' | 'POST' | 'PUT'

interface RequestConfig {
  method: SupportedMethod
  headers: Record<string, string>
  url: URL
  body?: Uint8Array
}

interface HttpClientOptions {
  baseURL: string
  beforeRequest?: (config: RequestConfig) => void
}

export class HttpClient {
  private baseURL: string
  private beforeRequest?: (config: RequestConfig) => void
  private client: typeof fetch

  constructor(options: HttpClientOptions) {
    this.baseURL = options.baseURL
    this.beforeRequest = options.beforeRequest
    this.client = fetch.bind(window)
  }

  async post(
    path: string,
    body: unknown,
    opts: Record<string, string> = {}
  ): Promise<unknown> {
    return this.fetch(path, 'POST', body, opts)
  }

  async patch(
    path: string,
    body: unknown,
    opts: Record<string, string> = {}
  ): Promise<unknown> {
    return this.fetch(path, 'PATCH', body, opts)
  }

  async put(
    path: string,
    body: unknown,
    opts: Record<string, string> = {}
  ): Promise<unknown> {
    return this.fetch(path, 'PUT', body, opts)
  }

  async delete(
    path: string,
    opts: Record<string, string> = {}
  ): Promise<unknown> {
    return this.fetch(path, 'DELETE', undefined, opts)
  }

  async get(path: string, opts: Record<string, string> = {}): Promise<unknown> {
    return this.fetch(path, 'GET', undefined, opts)
  }

  private async fetch(
    path: string,
    method: SupportedMethod,
    body?: unknown,
    opts: Record<string, string> = {}
  ): Promise<unknown> {
    const config = await this.makeConfig(path, method, body, opts)
    const response = await this.client(config.url, {
      method: config.method,
      headers: config.headers,
      body: config.body ? new Uint8Array(config.body) : undefined,
    })

    if (!response.ok) {
      const errorText = await response.text()
      if (!errorText) {
        throw new Error(`Unexpected status code: ${response.status}`)
      }
      throw new Error(
        `Unexpected status code: ${response.status}, response: ${errorText}`
      )
    }

    const text = await response.text()
    if (!text) {
      return null
    }

    return JSON.parse(text)
  }

  async fetchRawResponse(
    path: string,
    method: SupportedMethod,
    body?: unknown,
    opts: Record<string, string> = {}
  ): Promise<Response> {
    const config = await this.makeConfig(path, method, body, opts)
    return this.client(config.url, {
      method: config.method,
      headers: config.headers,
      body: config.body ? new Uint8Array(config.body) : undefined,
    })
  }

  private async makeConfig(
    path: string,
    method: SupportedMethod,
    body?: unknown,
    opts: Record<string, string> = {}
  ): Promise<RequestConfig> {
    const url = new URL(path, this.baseURL)

    let bodyBytes: Uint8Array | undefined
    if (body !== undefined) {
      bodyBytes = new TextEncoder().encode(JSON.stringify(body))
    }

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json',
      ...opts,
    }

    const config: RequestConfig = {
      method,
      headers,
      url,
      body: bodyBytes,
    }

    if (this.beforeRequest) {
      this.beforeRequest(config)
    }

    return config
  }
}
