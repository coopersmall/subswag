import { useCallback, useState } from 'react'
import { User } from '../types/user.generated'
import { useUsersHooks } from '../hooks/useUsersHooks'
import { Button, Flex, Heading, Text } from '@radix-ui/themes'
import DefaultLayout from '../components/lib/DefaultLayout'
import styled from 'styled-components'
import TopMenu from '../components/lib/TopMenu'
import LeftMenu from '../components/lib/LeftMenu'
import HomeMenu from '../components/lib/HomeMenu'
import { ChatMenu } from '../components/chatmenu/ChatMenu'

export default function HomePage() {
  return (
    <DefaultLayout>
      <HomePageContent />
    </DefaultLayout>
  )
}

export function HomePageContent() {
  const [users, setUsers] = useState<User[]>([])
  const { useGetAllUsers } = useUsersHooks()
  const [isGettingUsers, getUsers] = useGetAllUsers()

  const onGetUsers = useCallback(async () => {
    const users = await getUsers()
    setUsers(users)
  }, [setUsers, getUsers])

  const isLoading = isGettingUsers

  return (
    <HomePageLayout>
      <Flex direction="column" align="start" p="9" gap="2">
        <Heading as="h3">Home Page</Heading>
        <p>Welcome to the home page!</p>
        {users.length === 0 ? (
          <Button onClick={onGetUsers}>Load Users</Button>
        ) : (
          <Flex direction="column" align="start" gap="4">
            <ul>
              {users.map((user) => (
                <li key={user.id}>
                  {user.first_name} - {user.email}
                </li>
              ))}
            </ul>
            <Flex direction="row" align="start" gap="4">
              <Button onClick={onGetUsers}>Refresh Users</Button>
              <Button onClick={() => setUsers([])}>Clear Users</Button>
            </Flex>
          </Flex>
        )}
      </Flex>
    </HomePageLayout>
  )
}

export function HomePageLayout({ children }: { children: React.ReactNode }) {
  return (
    <StyledMain height="85vh">
      <LeftMenu />
      <Flex
        style={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
          width: '100%',
        }}
      >
        <Flex
          style={{
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'flex-end',
            height: '90%',
            border: '1px solid #ff0000',
          }}
        >
          <TopMenu 
          bottomChildren={<ChatMenu />}
          />
          <Flex
            style={{
              display: 'flex',
              flexDirection: 'row',
              height: '100%',
              width: '100%',
            }}
          >
            <HomeMenu />
            {children}
          </Flex>
        </Flex>
        <Flex
          style={{
            display: 'flex',
            flexDirection: 'row',
            height: '10%',
            minHeight: '70px',
          }}
        >
          <Flex
            style={{
              width: '15%',
              border: '1px solid #ff0000',
            }}
          >
            <Text>Middle Top</Text>
          </Flex>
          <Flex
            style={{
              width: '55%',
              border: '1px solid #ff0000',
            }}
          >
            <Text>Middle Top</Text>
          </Flex>
          <Flex
            style={{
              width: '30%',
              backgroundColor: '#000000',
              border: '1px solid #ff0000',
            }}
          >
            <Text>Middle Top</Text>
          </Flex>
        </Flex>
      </Flex>
    </StyledMain>
  )
}

const StyledMain = styled.main<{ height: string }>`
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  margin-right: 20px;
  height: ${(props) => props.height};
  width: 100%;
`
