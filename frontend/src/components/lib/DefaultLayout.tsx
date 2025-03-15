import { Container, Flex, Section } from '@radix-ui/themes'
import { Theme } from '@radix-ui/themes'
import { Header } from './Header'
import { UIProvider } from '../../contexts/UIStateContext'
import Footer from './Footer'

export default function DefaultLayout({
  children,
}: {
  children: React.ReactNode
  backUrl?: string
}) {
  return (
    <>
      <html lang="en">
        <head>
          <meta charSet="UTF-8" />
          <title>Document</title>
        </head>
        <body style={{ minHeight: '100vh' }}>
          <UIProvider>
            <Theme appearance="dark" style={{ minHeight: '100vh' }}>
              <Flex direction="column">
                <Header />
                <Flex>
                  <Container
                    data-aid="main-content"
                    style={{
                      margin: '0 16px',
                      paddingTop: '16px',
                    }}
                  >
                    <Section size="1">{children}</Section>
                  </Container>
                </Flex>
                <Footer />
              </Flex>
            </Theme>
          </UIProvider>
        </body>
      </html>
    </>
  )
}
