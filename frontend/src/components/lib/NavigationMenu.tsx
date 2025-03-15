import { useState } from 'react'
import styled from 'styled-components'
import {
  MoonIcon,
  RowsIcon,
  SunIcon,
  Cross1Icon,
  InstagramLogoIcon,
  TwitterLogoIcon,
  EnvelopeClosedIcon,
} from '@radix-ui/react-icons'
import { DropdownMenu } from '@radix-ui/themes'

export function NavigationMenu() {
  const [menuOpen, setMenuOpen] = useState(false)
  const [selected, setSelected] = useState(false)

  const handleMenuOpen = () => {
    setMenuOpen(!menuOpen)
    setSelected(!selected)
  }

  const setTriggerIcon = () => {
    if (selected) {
      return <Cross1Icon data-icon />
    } else {
      return <RowsIcon data-icon />
    }
  }

  return (
    <DropdownMenu.Root onOpenChange={handleMenuOpen}>
      <StyledTrigger isSelected={selected}>{setTriggerIcon()}</StyledTrigger>
      {menuOpen && (
        <StyledContent side="bottom" align="start">
          <StyledGroup>
            <StyledItem onSelect={() => console.log('Icon 1 clicked')}>
              <SunIcon data-icon />
            </StyledItem>
            <StyledItem onSelect={() => console.log('Icon 2 clicked')}>
              <MoonIcon data-icon />
            </StyledItem>
            <StyledItem onSelect={() => console.log('Icon 3 clicked')}>
              <EnvelopeClosedIcon data-icon />
            </StyledItem>
            <StyledItem onSelect={() => console.log('Icon 4 clicked')}>
              <InstagramLogoIcon data-icon />
            </StyledItem>
            <StyledItem onSelect={() => console.log('Icon 5 clicked')}>
              <TwitterLogoIcon data-icon />
            </StyledItem>
          </StyledGroup>
        </StyledContent>
      )}
    </DropdownMenu.Root>
  )
}

// Styled components
const StyledTrigger = styled(DropdownMenu.Trigger)<{ isSelected: boolean }>`
  padding: 10px;
  border: none; // Remove border
  border-radius: 8px;
  cursor: pointer; // Change cursor to pointer on hover
  ${({ isSelected }) =>
    isSelected &&
    `
    background-color: #ffffff; // Change as needed for your design
    [data-icon] {
      filter: invert(100%);
    }
  `}
  &:hover {
    background-color: #ffffff; // Example hover styling
    [data-icon] {
      filter: invert(100%); // Invert icon colors on hover of tooltip
    }
  }
  &:focus {
    outline: none; // Remove focus outline
  }
`

const StyledContent = styled(DropdownMenu.Content)``

const StyledItem = styled(DropdownMenu.Item)`
  padding: 10px;
  margin-bottom: 10px;
  margin-top: 10px;
  border: none; // Remove border
  border-radius: 8px;
  cursor: pointer; // Change cursor to pointer on hover
  &:hover {
    background-color: #ffffff; // Example hover styling
    [data-icon] {
      filter: invert(100%); // Invert icon colors on hover of tooltip
    }
  }
  &:focus {
    outline: none; // Remove focus outline
  }
`

const StyledGroup = styled(DropdownMenu.Group)`
  /* Your group styling */
`
