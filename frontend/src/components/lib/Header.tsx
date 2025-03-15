import { MagnifyingGlassIcon } from '@radix-ui/react-icons'
import { Flex } from '@radix-ui/themes'
import styled from 'styled-components'

export function Header() {
  return (
    <Flex
      style={{
        alignItems: 'center',
        display: 'flex',
        flexShrink: 0,
        height: '36px',
        marginTop: '12px',
        marginBottom: '12px',
        marginLeft: '30px',
        justifyContent: 'center',
      }}
    >
      <a href="/">
        <LogoIcon data-icon />
      </a>
      <Flex
        style={{
          alignItems: 'center',
          display: 'flex',
          flexBasis: '10%',
          flexGrow: 0,
          flexShrink: 0,
        }}
      ></Flex>

      <Flex
        style={{
          alignItems: 'center',
          justifyContent: 'center',
          display: 'flex',
          flexBasis: '60%',
          height: '64px',
          marginTop: '40px',
          flexGrow: 0,
          flexShrink: 0,
        }}
      >
        <SearchContainer width="100%" height="10px">
          <SearchButton width="20%" height="100%">
            <MagnifyingGlassIcon data-icon />
          </SearchButton>
          <SearchBar width="80%" height="100%" placeholder="Search" />
        </SearchContainer>
      </Flex>
      <div
        style={{
          display: 'flex',
          flexBasis: '20%',
          flexGrow: 0,
          flexShrink: 0,
          justifyContent: 'right',
          alignItems: 'flex-end',
        }}
      ></div>
    </Flex>
  )
}

function LogoIcon() {
  return (
    <Flex
      className="logo-icon"
      style={{
        backgroundImage: 'url(/images/logo.svg)',
        backgroundRepeat: 'no-repeat',
        backgroundSize: 'contain',
        filter: 'invert(100%)',
        height: '64px',
        marginTop: '40px',
        width: '50px',
      }}
    />
  )
}

const SearchButton = styled(Flex)`
  justify-content: center;
  align-items: center;
  transition: opacity 0.3s ease-in-out;
  cursor: pointer;
`

const SearchContainer = styled.div<{
  width: string
  height: string
}>`
  display: flex;
  align-items: center;
  padding: 10px;
  width: ${(props) => props.width};
  height: ${(props) => props.height};
  gap: 2;
  opacity: 0.5;
  transition: opacity 0.3s ease-in-out;
  &:hover {
    opacity: 1;
  }
`

const SearchBar = styled.input<{
  width: string
  height: string
}>`
  width: ${(props) => props.width};
  height: ${(props) => props.height};
  padding: 10px;
  border-radius: 25px;
  background-color: #1e1e1e;
  transition:
    opacity 0.3s ease-in-out,
    box-shadow 0.3s ease-in-out;
  opacity: 0.8;
`
