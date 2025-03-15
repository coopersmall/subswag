import { Flex } from '@radix-ui/themes'
import styled from 'styled-components'
import { UIContext } from '../../contexts/UIStateContext'
import { useContext, useCallback } from 'react'
// import { ChatMenu } from './ChatMenu'

export default function TopMenu({
    bottomChildren,
}: {
    bottomChildren: React.ReactNode;
}) {
  const { chatHeight, chatStartHeight } = useContext(UIContext)

  const isOpen = useCallback(() => {
    return chatHeight != chatStartHeight
  }, [chatHeight, chatStartHeight])

  return (
    <StyledTopMenu isOpen={isOpen}>
      <Flex
        style={{
          width: '70%',
          alignItems: 'center',
          gap: '100px',
        }}
      ></Flex>
      <Flex
        style={{
          width: '30%',
          justifyContent: 'flex-end',
          minWidth: '200px',
        }}
      >
      {bottomChildren}
      </Flex>
    </StyledTopMenu>
  )
}

const StyledTopMenu = styled(Flex)<{
  isOpen: () => boolean
}>`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  height: 32px;
  max-height: 100px;
  background-color: #ff0000;
  opacity: ${(props) => (props.isOpen() ? 1 : 0)};
  transition:
    opacity 0.5s,
    visibility 0.5s;
  &:hover {
    opacity: 1;
  }
`
