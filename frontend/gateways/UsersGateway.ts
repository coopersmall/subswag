import { IHttpClient } from '../clients/Clients'
import { User, UserData, UserID } from '../src/types/user.generated'
import { userSchema } from '../src/types/user.generated.zod'

export interface IUsersGateway {
  createUser(data: UserData): Promise<User>
  getUser(id: UserID): Promise<User>
  updateUser(id: UserID, data: UserData): Promise<User>
  deleteUser(id: UserID): Promise<void>
  listUsers(): Promise<User[]>
}

export class UsersGateway implements IUsersGateway {
  constructor(private readonly httpClient: IHttpClient) {}
  async createUser(data: UserData): Promise<User> {
    const response = await this.httpClient.post('/users', data)
    return userSchema.parse(response)
  }

  async listUsers(): Promise<User[]> {
    const response = (await this.httpClient.get('/api/users')) as Promise<
      User[]
    >
    return userSchema.array().parse(response)
  }

  async getUser(id: UserID): Promise<User> {
    const response = (await this.httpClient.get(
      `/users/${id}`
    )) as Promise<User>
    return userSchema.parse(response)
  }

  async updateUser(id: UserID, data: UserData): Promise<User> {
    const response = (await this.httpClient.patch(
      `/users/${id}`,
      data
    )) as Promise<User>
    return userSchema.parse(response)
  }

  async deleteUser(id: UserID): Promise<void> {
    return (await this.httpClient.delete(`/users/${id}`)) as Promise<void>
  }
}
