import { Flex } from '@radix-ui/themes'
import { PersonIcon, BookmarkIcon, BellIcon } from '@radix-ui/react-icons'
import styled from 'styled-components'
import { SidePanel } from './SidePanel'
import { useCallback, useContext } from 'react'
import { UIContext } from '../../contexts/UIStateContext'
import { StyledIconButton } from './IconButton'

export const LeftButtonIds = ['account', 'alerts', 'saved'] as const
export type LeftButtonId = (typeof LeftButtonIds)[number]

export default function LeftMenu() {
  const state = useContext(UIContext)
  return (
    <>
      <Flex
        style={{
          display: 'flex',
          flexDirection: 'column',
          height: state.mainHeight,
        }}
      >
        <Flex
          style={{
            padding: '10px',
            height: '70%',
            border: '1px solid #ff0000',
          }}
        >
          <LeftButtons />
        </Flex>
        <Flex
          style={{
            height: '30%',
            width: '100%',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            border: '1px solid #ff0000',
          }}
        ></Flex>
      </Flex>
      <LeftPanel />
      <Flex
        style={{
          display: 'flex',
          width: '2%',
        }}
      ></Flex>
    </>
  )
}

export function LeftButtons() {
  const state = useContext(UIContext)

  const onAccountClick = () => {
    if (state.lastLeftMenuButtonClicked === 'account') {
      state.updateLastLeftMenuButtonClicked('')
      state.updateLeftPanelExpanded(false)
      return
    }
    state.updateLastLeftMenuButtonClicked('account')
    state.updateLeftPanelExpanded(true)
  }

  const onAlertsClick = () => {
    if (state.lastLeftMenuButtonClicked === 'alerts') {
      state.updateLastLeftMenuButtonClicked('')
      state.updateLeftPanelExpanded(false)
      return
    }
    state.updateLastLeftMenuButtonClicked('alerts')
    state.updateLeftPanelExpanded(true)
  }

  const onSaveClick = () => {
    if (state.lastLeftMenuButtonClicked === 'saved') {
      state.updateLastLeftMenuButtonClicked('')
      state.updateLeftPanelExpanded(false)
      return
    }
    state.updateLastLeftMenuButtonClicked('saved')
    state.updateLeftPanelExpanded(true)
  }

  return (
    <ButtonContainer>
      <AccountButton
        isSelected={state.lastLeftMenuButtonClicked === 'account'}
        onClick={onAccountClick}
      />
      <AlertsButton
        isSelected={state.lastLeftMenuButtonClicked === 'alerts'}
        onClick={onAlertsClick}
      />
      <SaveButton
        isSelected={state.lastLeftMenuButtonClicked === 'saved'}
        onClick={onSaveClick}
      />
    </ButtonContainer>
  )
}

function AccountButton({
  isSelected,
  onClick,
}: {
  isSelected: boolean
  onClick: () => void
}) {
  const backgroundColor = '#000000'
  const hoverColor = '#ffffff'
  const solidSelectedColor = '#ffffff'
  const invertIcon = true

  const resolve = () => {
    onClick()
  }

  return (
    <StyledButton
      onClick={resolve}
      backgroundColor={backgroundColor}
      hoverColor={hoverColor}
      solidSelected={isSelected}
      solidColor={solidSelectedColor}
      invertIcon={invertIcon}
    >
      <PersonIcon data-icon />
    </StyledButton>
  )
}

interface AlertsButtonProps {
  isSelected: boolean
  onClick: () => void
}

function AlertsButton({ isSelected, onClick }: AlertsButtonProps) {
  const backgroundColor = '#000000'
  const hoverColor = '#ffffff'
  const solidSelectedColor = '#ffffff'
  const invertIcon = true

  const resolve = () => {
    onClick()
  }

  return (
    <StyledButton
      onClick={resolve}
      backgroundColor={backgroundColor}
      hoverColor={hoverColor}
      solidSelected={isSelected}
      solidColor={solidSelectedColor}
      invertIcon={invertIcon}
    >
      <BellIcon data-icon />
    </StyledButton>
  )
}

interface SaveButtonProps {
  isSelected: boolean
  onClick: () => void
}

function SaveButton({ onClick }: SaveButtonProps) {
  const backgroundColor = '#000000'
  const hoverColor = '#ffffff'
  const invertIcon = true

  const resolve = () => {
    onClick()
  }

  // const icon = state.isSessionSaved ? (
  //   <BookmarkFilledIcon data-icon="bookmark-filled" />
  // ) : (
  //   <BookmarkIcon data-icon="bookmark" />
  // )

  const icon = <BookmarkIcon data-icon="bookmark" />
  return (
    <StyledButton
      onClick={resolve}
      backgroundColor={backgroundColor}
      hoverColor={hoverColor}
      invertIcon={invertIcon}
    >
      {icon}
    </StyledButton>
  )
}

export function LeftPanel() {
  const maxWidth = '400px'
  const state = useContext(UIContext)

  const wasLeftPanelDoubleClicked = useCallback(
    (direction: string): boolean => {
      if (direction !== state.lastLeftBorderClicked) return false
      if (state.leftborderClickTime === 0) return false
      if (Date.now() - state.leftborderClickTime > 300) return false
      return true
    },
    [state]
  )

  const onResizeStart = useCallback(
    (e, direction, ref, delta, position) => {
      state.updateLastLeftBorderClicked(direction)
      if (wasLeftPanelDoubleClicked(direction)) {
        if (state.leftPanelWidth === maxWidth) {
          ref.style.width = state.leftPanelStartWidth
        }
        if (state.leftPanelWidth === state.leftPanelStartWidth) {
          ref.style.width = state.leftPanelWidth
        }
      }
    },
    [state, wasLeftPanelDoubleClicked]
  )

  const onResizeStop = useCallback(
    (e, direction, ref, delta, position) => {
      state.updateLeftPanelWidth(ref.style.width)
    },
    [state]
  )

  return (
    <SidePanel
      x={56}
      y={0}
      height={state.mainHeight}
      width={state.leftPanelStartWidth}
      side="left"
      maxWidth={maxWidth}
      isExpanded={state.isLeftPanelExpanded}
      onResizeStart={onResizeStart}
      onResizeStop={onResizeStop}
    >
      <Flex
        style={{
          display: 'flex',
          flexDirection: 'row',
          height: '5%',
          minHeight: '20px',
          border: '1px solid #ff0000',
        }}
      ></Flex>
      <Flex
        style={{
          backgroundColor: '#ff0000',
          height: '95%',
        }}
      ></Flex>
    </SidePanel>
  )
}

const ButtonContainer = styled(Flex)`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 10px;
`

const StyledButton = styled(StyledIconButton)`
  margin-top: 10px;
  margin-bottom: 10px;
`
