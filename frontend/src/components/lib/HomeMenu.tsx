import styled from 'styled-components'
import React, { forwardRef, useCallback, useContext, useEffect } from 'react'
import { Position, ResizableDelta, Rnd, RndResizeCallback } from 'react-rnd'
import { BoxIcon } from '@radix-ui/react-icons'
import { Flex } from '@radix-ui/themes'
import { useRef } from 'react'
import { UIContext } from '../../contexts/UIStateContext'

export const HomeButtons = 'home'
export type HomeButtonId = (typeof HomeButtons)[number]

export default function HomeMenu() {
  const state = useContext(UIContext)

  const panelRef = useRef<Rnd>(null)
  useEffect(() => {
    state.updateHomePanel(panelRef)
  }, [state, panelRef])

  const onHomeButtonClicked = useCallback(() => {
    state.updateLastHomeButtonClicked('home')
    if (!state.homePanel.current) return
    if (
      state.homeWidth === state.homeStartWidth &&
      state.homeHeight === state.homeStartHeight
    ) {
      state.updateHomeWidth('100%')
      state.updateHomeHeight('100%')
    } else {
      state.updateHomeWidth(state.homeStartWidth)
      state.updateHomeHeight(state.homeStartHeight)
    }
  }, [state])

  const wasBorderDoubleClicked = useCallback(
    (direction: string): boolean => {
      if (direction !== state.lastHomeBorderClicked) return false
      if (state.lastHomeBorderClickedTime === 0) return false
      if (Date.now() - state.lastHomeBorderClickedTime > 300) return false
      return true
    },
    [state]
  )

  const onResizeStart = useCallback(
    (
      e: MouseEvent | TouchEvent,
      direction,
      ref: HTMLElement,
      delta: ResizableDelta,
      position: Position
    ) => {
      state.updateLastHomeBorderClicked(direction)
      if (wasBorderDoubleClicked(direction)) {
        const { height: newHeigth, width: newWidth } = onDoubleCick(
          direction,
          state.homeStartWidth,
          state.homeStartHeight,
          state.homeHeight,
          state.homeWidth
        )
        state.updateHomeHeight(newHeigth)
        state.updateHomeWidth(newWidth)
      }
    },
    [state, wasBorderDoubleClicked]
  )

  const onResizeStop = useCallback(
    (
      e: MouseEvent | TouchEvent,
      direction,
      ref: HTMLElement,
      delta: ResizableDelta,
      position: Position
    ) => {
      state.updateHomeHeight(ref.style.height)
      state.updateHomeWidth(ref.style.width)
    },
    [state]
  )

  return (
    <>
      <Flex
        style={{
          width: '100%',
          height: '100%',
          display: 'flex',
          flexDirection: 'row',
          gap: '10px',
          border: '1px solid #ff0000',
        }}
      >
        <HomePanelRef
          ref={panelRef}
          width={state.homeWidth}
          height={state.homeHeight}
          startWidth={state.homeStartWidth}
          startHeight={state.homeStartHeight}
          onResizeStart={onResizeStart}
          onResizeStop={onResizeStop}
        >
          <StyledHomeMenuBar>
            <HomeButton onClick={onHomeButtonClicked} />
          </StyledHomeMenuBar>
        </HomePanelRef>
      </Flex>
    </>
  )
}

interface HomeButtonProps {
  onClick: () => void
}

function HomeButton({ onClick }: HomeButtonProps) {
  const state = useContext(UIContext)

  const icon = <BoxIcon data-icon />

  return (
    <IconButton
      invertIcon={true}
      solidSelected={
        state.homeWidth !== state.homeStartWidth ||
        state.homeStartHeight !== state.homeHeight
      }
      solidColor="#ffffff"
      onClick={onClick}
    >
      {icon}
    </IconButton>
  )
}

function onDoubleCick(
  direction: string,
  startWidth: string,
  startHeight: string,
  height: string,
  width: string
): { height: string; width: string } {
  if (direction === 'right') {
    if (width === startWidth) {
      return { height: height, width: '100%' }
    } else {
      return { height: height, width: startWidth }
    }
  } else if (direction === 'bottom') {
    if (height === startHeight) {
      return { height: '100%', width: width }
    } else {
      return { height: startHeight, width: width }
    }
  } else if (direction === 'bottomRight') {
    if (height === startHeight && width === startWidth) {
      return { height: '100%', width: '100%' }
    } else {
      return { height: startHeight, width: startWidth }
    }
  }
  return { height: height, width: width }
}

const StyledHomeMenuBar = styled(Flex)`
  height: 64px;
  width: 64px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: #000000;
`

interface HomePanelRefProps {
  width: string
  height: string
  startWidth: string
  startHeight: string
  children?: React.ReactNode
  onResizeStart: RndResizeCallback
  onResizeStop: RndResizeCallback
}

const HomePanelRef = forwardRef<Rnd, HomePanelRefProps>(
  (
    {
      width,
      height,
      startWidth,
      startHeight,
      children,
      onResizeStart,
      onResizeStop,
    },
    ref
  ) => (
    <StyledRnd
      isExpanded={width !== startWidth || height !== startHeight}
      size={{ width, height }}
      ref={ref}
      style={{
        position: 'relative',
        display: 'flex',
      }}
      bounds="parent"
      enableResizing={{
        // right: true,
        // bottom: true,
        bottomRight: true,
        bottomLeft: true,
        topRight: true,
        topLeft: true,
      }}
      minHeight={startHeight}
      // disableDragging={true}
      minWidth={startWidth}
      onResizeStart={onResizeStart}
      onResizeStop={onResizeStop}
    >
      {children}
    </StyledRnd>
  )
)

HomePanelRef.displayName = 'HomePanelRef'

interface StyledHomePanelRefProps {
  isExpanded: boolean
}

const StyledRnd = styled(Rnd)<StyledHomePanelRefProps>`
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

interface IconButtonProps {
  backgroundColor?: string
  hoverColor?: string
  solidSelected?: boolean
  solidColor?: string
  invertIcon?: boolean
}

const IconButton = styled(Flex)<IconButtonProps>`
  padding: 10px;
  width: 36px;
  border-radius: 8px;
  display: flex;
  cursor: pointer;
  ${(props) =>
    props.solidSelected
      ? `background-color: ${props.solidColor || '#ffffff'};`
      : `background-color: ${props.backgroundColor || '#000000'};`}
  ${(props) =>
    props.solidSelected
      ? `
  [data-icon] {
    filter: invert(100%);
  }
  `
      : ``}
  &:hover {
    background-color: ${(props) => props.hoverColor || '#ffffff'};

    [data-icon] {
      filter: ${(props) => (props.invertIcon ? 'invert(100%)' : 'none')};
    }
  }
`
