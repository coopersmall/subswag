import React, { forwardRef, useCallback, useContext, useEffect, useRef } from 'react'
import { Rnd, RndResizeCallback } from 'react-rnd'
import styled from 'styled-components'
import { Flex } from '@radix-ui/themes'
import { DotFilledIcon } from '@radix-ui/react-icons'
import { UIContext } from '../../contexts/UIStateContext'
import RotatingNotification from './RotatingNotification'

export function ChatMenu() {
  const state = useContext(UIContext)

  const panelRef = useRef<Rnd>(null)
  useEffect(() => {
    state.updateChatPanel(panelRef)
  }, [state, panelRef])

  const wasBorderDoubleClicked = useCallback(
    (direction: string): boolean => {
      if (direction !== state.lastChatBorderClicked) return false
      if (state.lastChatBorderClickedTime === 0) return false
      if (Date.now() - state.lastChatBorderClickedTime > 300) return false
      return true
    },
    [state.lastChatBorderClicked, state.lastChatBorderClickedTime]
  )

  const onResizeStart = useCallback(
    (e, direction, ref, delta, position) => {
      state.updateLastChatBorderClicked(direction)
      state.updateLastChatBorderClickedTime(Date.now())
      if (wasBorderDoubleClicked(direction)) {
      if (state.chatHeight === state.chatStartHeight) {
        state.updateChatHeight(state.chatMaxHeight)
      } else {
         state.updateChatHeight(state.chatStartHeight)
      }
      }
    },
    [
        state,
        wasBorderDoubleClicked,
    ]
  )

  const onResizeStop = useCallback(
    (e, direction, ref, delta, position) => {
        state.updateChatHeight(ref.style.height)
        state.updateChatWidth(ref.style.width)
    },
    [state]
  )

  return (
    <>
      <Flex
        style={{
          display: 'flex',
          flexDirection: 'row',
          gap: '10px',
        }}
        gap="10px"
        direction="row"
        width="100%"
      >
    <ChatPanelRef
    ref={panelRef}
      height={state.chatHeight}
      startWidth={state.chatStartWidth}
      startHeight={state.chatStartHeight}
      maxHeigth={state.chatMaxHeight}
      onResizeStart={onResizeStart}
      onResizeStop={onResizeStop}
      >
    <StyledChatMenuBar
        height={state.chatStartHeight}
        width={state.chatStartWidth}
        isExpanded={state.chatHeight !== state.chatStartHeight}
    >
    <Flex
        justify="start"
        align="center"
        direction="row"
        gap="10px"
    >
    <RotatingNotification
     messages={[
         'Welcome to the chat!',
         'Please be respectful to others',
         'No spamming or trolling',
         'Have fun!',
      ]}
    />
    <StyledOnlineIcon
       data-icon
       color="green"
       width="30px"
       height="30px"
       />
    </Flex>
    </StyledChatMenuBar>
    </ChatPanelRef>
    </Flex>
    </>
  )
}

const ChatPanelRef = forwardRef<Rnd, {
    height: string
    startWidth: string
    startHeight: string
    maxHeigth: string
    children?: React.ReactNode
    onResizeStart: RndResizeCallback
    onResizeStop: RndResizeCallback
}>(
  (
    {
      height,
      startWidth,
      startHeight,
      maxHeigth,
      children,
      onResizeStart,
      onResizeStop,
    },
    ref
  ) => (
    <StyledRnd
      isExpanded={height !== startHeight}
      ref={ref}
      style={{
        position: 'relative',
        display: 'flex',
      }}
      enableResizing={{
        bottom: true,
      }}
      height={height}
      minHeight={startHeight}
      maxHeight={maxHeigth}
      disableDragging={true}
      minWidth={startWidth}
      width={startWidth}
      onResizeStart={onResizeStart}
      onResizeStop={onResizeStop}
    >
      {children}
    </StyledRnd>
  )
)

ChatPanelRef.displayName = 'ChatPanelRef'

const StyledRnd = styled(Rnd)<{
    isExpanded: boolean
}>`
  background-color: #1c1c1c;
  opacity: ${({ isExpanded }) => (isExpanded ? 0.8 : 0.2)};
  transition:
    opacity 0.4s ease-in-out,
    box-shadow 0.4s ease-in-out;
  border-radius: 5px;
  &:hover {
    opacity: 1;
    box-shadow: 0 0 0 2px #ff0000;
  }
`

const StyledChatMenuBar = styled(Flex)<{
    height: string;
    width: string;
    isExpanded: boolean;
}>`
height: ${({ height }) => height};
width: ${({ width }) => width};
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-end;
  padding: 10px;
  background-color: #000000;
  transition: opacity 0.4s ease-in-out;
`

const StyledOnlineIcon = styled(DotFilledIcon)`
  width: 32px;
  height: 32px;
  color: green;
  cursor: pointer;
  
  @keyframes pulse {
    0% {
      opacity: 1;
      transform: scale(1);
    }
    50% {
      opacity: 0.7;
      transform: scale(1.1);
    }
    100% {
      opacity: 1;
      transform: scale(1);
    }
  }
  animation: pulse 2s infinite;
`
