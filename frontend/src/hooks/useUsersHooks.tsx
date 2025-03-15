import { useCallback, useState } from 'react'
import { IUsersGateway } from '../../gateways/UsersGateway'
import { User, UserData, UserID } from '../types/user.generated'
import { useGateways } from './useGateways'

export interface IUsersHooks {
  useGetAllUsers: () => [
    isLoading: boolean,
    loadUsers: () => Promise<Array<User>>,
  ]
  useGetUserById: () => [
    isLoading: boolean,
    loadUser: (id: UserID) => Promise<User>,
  ]
  useCreateUser: () => [
    isLoading: boolean,
    createUser: (data: UserData) => Promise<User>,
  ]
  useUpdateUser: () => [
    isLoading: boolean,
    updateUser: (id: UserID, data: UserData) => Promise<User>,
  ]
  useDeleteUser: () => [
    isLoading: boolean,
    deleteUser: (id: UserID) => Promise<void>,
  ]
}

export function useUsersHooks(): IUsersHooks {
  const { useUsersGateway } = useGateways()
  const usersGateway = useUsersGateway()
  return {
    useGetAllUsers: () => useGetAllUsers(usersGateway),
    useGetUserById: () => useGetUserById(usersGateway),
    useCreateUser: () => useCreateUser(usersGateway),
    useUpdateUser: () => useUpdateUser(usersGateway),
    useDeleteUser: () => useDeleteUser(usersGateway),
  }
}

export function useGetAllUsers(
  usersGateway: IUsersGateway
): [isLoading: boolean, loadUsers: () => Promise<Array<User>>] {
  const [isLoading, setIsLoading] = useState(false)

  const loadUsers = useCallback(async () => {
    setIsLoading(true)
    const users = await usersGateway.listUsers()
    setIsLoading(false)
    return users
  }, [usersGateway])

  return [isLoading, loadUsers]
}

export function useGetUserById(
  usersGateway: IUsersGateway
): [isLoading: boolean, loadUser: (id: UserID) => Promise<User>] {
  const [isLoading, setIsLoading] = useState(false)

  const loadUser = useCallback(
    async (id: number) => {
      setIsLoading(true)
      const user = await usersGateway.getUser(id)
      setIsLoading(false)
      return user
    },
    [usersGateway]
  )

  return [isLoading, loadUser]
}

export function useCreateUser(
  usersGateway: IUsersGateway
): [isLoading: boolean, createUser: (data: UserData) => Promise<User>] {
  const [isLoading, setIsLoading] = useState(false)

  const createUser = useCallback(
    async (data: UserData) => {
      setIsLoading(true)
      const user = await usersGateway.createUser(data)
      setIsLoading(false)
      return user
    },
    [usersGateway]
  )

  return [isLoading, createUser]
}

export function useUpdateUser(
  usersGateway: IUsersGateway
): [
  isLoading: boolean,
  updateUser: (id: UserID, data: UserData) => Promise<User>,
] {
  const [isLoading, setIsLoading] = useState(false)

  const updateUser = useCallback(
    async (id: UserID, data: UserData) => {
      setIsLoading(true)
      const user = await usersGateway.updateUser(id, data)
      setIsLoading(false)
      return user
    },
    [usersGateway]
  )

  return [isLoading, updateUser]
}

export function useDeleteUser(
  usersGateway: IUsersGateway
): [isLoading: boolean, deleteUser: (id: UserID) => Promise<void>] {
  const [isLoading, setIsLoading] = useState(false)

  const deleteUser = useCallback(
    async (id: UserID) => {
      setIsLoading(true)
      await usersGateway.deleteUser(id)
      setIsLoading(false)
    },
    [usersGateway]
  )

  return [isLoading, deleteUser]
}
