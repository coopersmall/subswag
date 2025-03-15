import { Flex } from '@radix-ui/themes'
import styled from 'styled-components'

export const StyledIconButton = styled(Flex)<{
  backgroundColor?: string
  hoverColor?: string
  solidSelected?: boolean
  solidColor?: string
  invertIcon?: boolean
}>`
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
