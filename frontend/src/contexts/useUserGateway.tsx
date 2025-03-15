import { createContext, useMemo } from 'react'
import { IUsersGateway, UsersGateway } from '../..//gateways/UsersGateway'
import { IHttpClient } from '../../clients/Clients'

export const UsersGatewayContext = createContext<IUsersGateway | null>(null)

interface UserGatewayProviderProps {
  httpClient: IHttpClient
  children: React.ReactNode
}

export const UserGatewayProvider: React.FC<UserGatewayProviderProps> = ({
  httpClient,
  children,
}) => {
  const gateway = useMemo(() => new UsersGateway(httpClient), [httpClient])
  return (
    <UsersGatewayContext.Provider value={gateway}>
      {children}
    </UsersGatewayContext.Provider>
  )
}
