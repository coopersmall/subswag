import { useMemo } from 'react'
import { HttpClient } from '../../clients/HttpClient'
import { IUsersGateway, UsersGateway } from '../../gateways/UsersGateway'

export interface IGateways {
  useUsersGateway: () => IUsersGateway
}

export function useGateways(): IGateways {
  const usersGateway = useMemo(() => {
    const httpClient = new HttpClient({
      baseURL: 'http://localhost:9000',
      beforeRequest: (config) => {
        const authHeader = `Bearer ${localStorage.getItem('token')}` // update this to use a real token
        config.headers.Authorization = authHeader
        config.headers['Content-Type'] = 'application/json'
      },
    })
    return new UsersGateway(httpClient)
  }, [])

  return {
    useUsersGateway: () => usersGateway,
  }
}
