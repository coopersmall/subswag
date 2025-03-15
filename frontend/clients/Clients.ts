export interface IHttpClient {
  post(
    path: string,
    body: unknown,
    opts?: Record<string, string>
  ): Promise<unknown>
  patch(
    path: string,
    body: unknown,
    opts?: Record<string, string>
  ): Promise<unknown>
  put(
    path: string,
    body: unknown,
    opts?: Record<string, string>
  ): Promise<unknown>
  delete(path: string, opts?: Record<string, string>): Promise<unknown>
  get(path: string, opts?: Record<string, string>): Promise<unknown>
}
