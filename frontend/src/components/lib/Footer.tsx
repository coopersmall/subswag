import { useState, useEffect, useRef } from 'react'
import styled from 'styled-components'
import { Flex } from '@radix-ui/themes'
import {
  EnvelopeClosedIcon,
  InstagramLogoIcon,
  TwitterLogoIcon,
} from '@radix-ui/react-icons'
import { StyledIconButton } from './IconButton'

export default function Footer() {
  const [isVisible, setIsVisible] = useState(false)
  const footerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!footerRef.current) return
      const footerHeight = footerRef.current.clientHeight
      const bottomThreshold = window.innerHeight - footerHeight

      if (e.clientY >= bottomThreshold) {
        setIsVisible(true)
      }
      if (e.clientY <= window.innerHeight - footerHeight) {
        setIsVisible(false)
      }
    }

    window.addEventListener('mousemove', handleMouseMove)

    return () => {
      window.removeEventListener('mousemove', handleMouseMove)
    }
  }, [])

  return (
    <Flex
      style={{
        position: 'fixed',
        bottom: '0',
        display: 'flex',
        flexDirection: 'column',
        padding: '10px',
        alignItems: 'flex-start',
        width: '100%',
        height: '52px',
      }}
    >
      <FooterContainer isVisible={isVisible} ref={footerRef}>
        <Flex
          style={{
            display: 'flex',
            backgroundColor: '#000000',
            width: '20%',
          }}
        >
          <FooterIcon>
            <EnvelopeClosedIcon data-icon />
          </FooterIcon>
          <FooterIcon>
            <InstagramLogoIcon data-icon />
          </FooterIcon>
          <FooterIcon>
            <TwitterLogoIcon data-icon />
          </FooterIcon>
        </Flex>
      </FooterContainer>
    </Flex>
  )
}

interface FooterProps {
  isVisible: boolean
}

const FooterContainer = styled(Flex)<FooterProps>`
  visibility: ${({ isVisible }) => (isVisible ? 'visible' : 'hidden')};
  opacity: ${(props) => (props.isVisible ? 1 : 0)};
  transition:
    opacity 0.5s ease,
    visibility 0.5s ease;
`

const FooterIcon = styled(StyledIconButton)`
  margin-left: 5px;
  margin-right: 5px;
`
