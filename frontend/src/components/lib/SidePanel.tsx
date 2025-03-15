import styled from 'styled-components'
import { Rnd } from 'react-rnd'

interface SidePanelProps {
  x: number
  y: number
  height: string
  width: string
  maxWidth: string
  isExpanded: boolean
  children?: React.ReactNode
  side: 'left' | 'right'
  onResizeStart?: (e, direction, ref, delta, position) => void
  onResizeStop?: (e, direction, ref, delta, position) => void
}

export function SidePanel({
  x,
  y,
  height,
  width,
  maxWidth,
  isExpanded,
  children,
  side,
  onResizeStart,
  onResizeStop,
}: SidePanelProps) {
  if (!isExpanded) {
    return null // Return null to render nothing when the panel is not expanded
  }
  return (
    <StyledSidePanel
      isExpanded={isExpanded}
      expandedWidth={width}
      side={side}
      position={{
        x: x,
        y: y,
      }}
      size={{
        width: width,
        height: height,
      }}
      disableDragging={true}
      enableResizing={{ right: true }}
      minHeight={height}
      minWidth={width}
      maxWidth={maxWidth}
      bounds="parent"
      style={{
        position: 'relative',
        display: 'flex',
      }}
      onResizeStart={onResizeStart}
      onResizeStop={onResizeStop}
    >
      {children}
    </StyledSidePanel>
  )
}

interface StyledSidePanelProps {
  isExpanded: boolean
  expandedWidth: string
}

const StyledSidePanel = styled(Rnd)<StyledSidePanelProps>`
  display: flex;
  flex: ${(props) =>
    props.isExpanded ? `0 0 ${props.expandedWidth}` : '0 0 0'};
  flex-direction: column;
  transition:
    opacity 0.2s ease-in-out,
    visibility 0.2s ease-in-out;
  opacity: ${(props) => (props.isExpanded ? 1 : 0)};
  visibility: ${(props) => (props.isExpanded ? 'visible' : 'hidden')};
  overflow: hidden;
  border: ${(props) =>
    props.isExpanded ? '1px solid #ff0000' : '0px solid #ff0000'};
  box-sizing: border-box; // Include border in the element's dimensions
`
