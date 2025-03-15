import { Rnd } from 'react-rnd'
import { RefObject, createContext, useState, ReactNode, useRef } from 'react'

export interface IUIState {
  mainHeight: string
  updateMainHeight: (height: string) => void

  leftPanel: RefObject<Rnd | null>
  updateLeftPanel: (leftPanel: RefObject<Rnd>) => void
  leftPanelStartWidth: string
  leftPanelMaxWidth: string
  leftPanelWidth: string
  updateLeftPanelWidth: (width: string) => void
  isLeftPanelExpanded: boolean
  updateLeftPanelExpanded: (toggle: boolean) => void
  lastLeftMenuButtonClicked: string | null
  updateLastLeftMenuButtonClicked: (button: string) => void
  lastLeftBorderClicked: string
  updateLastLeftBorderClicked: (direction: string) => void
  leftborderClickTime: number
  updateLeftBorderClickTime: (time: number) => void

  isRightPanelExpanded: boolean
  updateRightPanelExpanded: (toggle: boolean) => void
  lastRightMenuButtonClicked: string | null
  updateLastRightMenuButtonClicked: (button: string) => void

  homePanel: RefObject<Rnd | null>
  updateHomePanel: (homePanel: RefObject<Rnd>) => void
  homeStartWidth: string
  homeStartHeight: string
  homeHeight: string
  updateHomeHeight: (height: string) => void
  homeWidth: string
  updateHomeWidth: (width: string) => void
  lastHomeButtonClicked: string | null
  updateLastHomeButtonClicked: (button: string) => void
  lastHomeButtonClickedTime: number
  updateLastHomeButtonClickedTime: (time: number) => void
  lastHomeBorderClicked: string
  updateLastHomeBorderClicked: (direction: string) => void
  lastHomeBorderClickedTime: number
  updateLastHomeBorderClickedTime: (time: number) => void

  chatPanel: RefObject<Rnd | null>
  updateChatPanel: (chatPanel: RefObject<Rnd>) => void
  chatStartWidth: string
  chatStartHeight: string
  chatMaxHeight: string
  chatHeight: string
  updateChatHeight: (height: string) => void
  chatWidth: string
  updateChatWidth: (width: string) => void
  lastChatButtonClicked: string | null
  updateLastChatButtonClicked: (button: string) => void
  lastChatButtonClickedTime: number
  updateLastChatButtonClickedTime: (time: number) => void
  lastChatBorderClicked: string
  updateLastChatBorderClicked: (direction: string) => void
  lastChatBorderClickedTime: number
  updateLastChatBorderClickedTime: (time: number) => void
}

// Create a default state with no-op functions for the updates
const defaultUIState: IUIState = {
  mainHeight: '80vh',
  updateMainHeight: () => {},

  leftPanel: { current: null },
  updateLeftPanel: () => {},
  leftPanelStartWidth: '200px',
  leftPanelMaxWidth: '500px',
  leftPanelWidth: '200px',
  updateLeftPanelWidth: () => {},
  isLeftPanelExpanded: false,
  updateLeftPanelExpanded: () => {},
  lastLeftMenuButtonClicked: null,
  updateLastLeftMenuButtonClicked: () => {},
  lastLeftBorderClicked: '',
  updateLastLeftBorderClicked: () => {},
  leftborderClickTime: 0,
  updateLeftBorderClickTime: () => {},

  isRightPanelExpanded: false,
  updateRightPanelExpanded: () => {},
  lastRightMenuButtonClicked: null,
  updateLastRightMenuButtonClicked: () => {},

  homePanel: { current: null },
  updateHomePanel: () => {},
  homeStartWidth: '64px',
  homeStartHeight: '64px',
  homeHeight: '64px',
  updateHomeHeight: () => {},
  homeWidth: '64px',
  updateHomeWidth: () => {},
  lastHomeButtonClicked: null,
  updateLastHomeButtonClicked: () => {},
  lastHomeButtonClickedTime: 0,
  updateLastHomeButtonClickedTime: () => {},
  lastHomeBorderClicked: '',
  updateLastHomeBorderClicked: () => {},
  lastHomeBorderClickedTime: 0,
  updateLastHomeBorderClickedTime: () => {},

  chatPanel: { current: null },
  updateChatPanel: () => {},
  chatStartWidth: '100%',
  chatStartHeight: '32px',
  chatMaxHeight: '70vh',
  chatHeight: '32px',
  updateChatHeight: () => {},
  chatWidth: '100%',
  updateChatWidth: () => {},
  lastChatButtonClicked: null,
  updateLastChatButtonClicked: () => {},
  lastChatButtonClickedTime: 0,
  updateLastChatButtonClickedTime: () => {},
  lastChatBorderClicked: '',
  updateLastChatBorderClicked: () => {},
  lastChatBorderClickedTime: 0,
  updateLastChatBorderClickedTime: () => {},
}

// Create the context
export const UIContext = createContext<IUIState>(defaultUIState)

// Create a provider component
interface UIProviderProps {
  children: ReactNode
}

export const UIProvider = ({ children }: UIProviderProps) => {
  const [mainHeight, setMainHeight] = useState('80vh')

  const [leftPanel, setLeftPanel] = useState<RefObject<Rnd | null>>(
    useRef(null)
  )
  const [leftPanelWidth, setLeftPanelWidth] = useState('250px')
  const [isLeftPanelExpanded, setIsLeftPanelExpanded] = useState(false)
  const [lastLeftMenuButtonClicked, setLastLeftMenuButtonClicked] = useState<
    string | null
  >(null)
  const [lastLeftBorderClicked, setLastLeftBorderClicked] = useState('')
  const [leftborderClickTime, setLeftBorderClickTime] = useState(0)

  const [isRightPanelExpanded, setIsRightPanelExpanded] = useState(false)
  const [lastRightMenuButtonClicked, setLastRightMenuButtonClicked] = useState<
    string | null
  >(null)

  const [homePanel, setHomePanel] = useState<RefObject<Rnd | null>>(
    useRef(null)
  )
  const [homeHeight, setHomeHeight] = useState('64px')
  const [homeWidth, setHomeWidth] = useState('64px')
  const [lastHomeButtonClicked, setLastHomeButtonClicked] = useState<
    string | null
  >(null)
  const [lastHomeButtonClickedTime, setLastHomeButtonClickedTime] = useState(0)
  const [lastHomeBorderClicked, setLastHomeBorderClicked] = useState('')
  const [lastHomeBorderClickedTime, setLastHomeBorderClickedTime] = useState(0)

  const [chatPanel, setChatPanel] = useState<RefObject<Rnd | null>>(
    useRef(null)
  )
  const [chatHeight, setChatHeight] = useState('32px')
  const [chatWidth, setChatWidth] = useState('100%')
  const [lastChatButtonClicked, setLastChatButtonClicked] = useState<
    string | null
  >(null)
  const [lastChatButtonClickedTime, setLastChatButtonClickedTime] = useState(0)
  const [lastChatBorderClicked, setLastChatBorderClicked] = useState('')
  const [lastChatBorderClickedTime, setLastChatBorderClickedTime] = useState(0)

  const value: IUIState = {
    mainHeight,
    updateMainHeight: setMainHeight,

    leftPanel,
    updateLeftPanel: setLeftPanel,
    leftPanelStartWidth: '250px',
    leftPanelMaxWidth: '500px',
    leftPanelWidth,
    updateLeftPanelWidth: setLeftPanelWidth,
    isLeftPanelExpanded,
    updateLeftPanelExpanded: setIsLeftPanelExpanded,
    lastLeftMenuButtonClicked,
    updateLastLeftMenuButtonClicked: setLastLeftMenuButtonClicked,
    lastLeftBorderClicked,
    updateLastLeftBorderClicked: setLastLeftBorderClicked,
    leftborderClickTime,
    updateLeftBorderClickTime: setLeftBorderClickTime,

    isRightPanelExpanded,
    updateRightPanelExpanded: setIsRightPanelExpanded,
    lastRightMenuButtonClicked,
    updateLastRightMenuButtonClicked: setLastRightMenuButtonClicked,

    homePanel,
    updateHomePanel: setHomePanel,
    homeStartWidth: '64px',
    homeStartHeight: '64px',
    homeHeight,
    updateHomeHeight: setHomeHeight,
    homeWidth,
    updateHomeWidth: setHomeWidth,
    lastHomeButtonClicked,
    updateLastHomeButtonClicked: setLastHomeButtonClicked,
    lastHomeButtonClickedTime,
    updateLastHomeButtonClickedTime: setLastHomeButtonClickedTime,
    lastHomeBorderClicked,
    updateLastHomeBorderClicked: setLastHomeBorderClicked,
    lastHomeBorderClickedTime,
    updateLastHomeBorderClickedTime: setLastHomeBorderClickedTime,

    chatPanel,
    updateChatPanel: setChatPanel,
    chatStartWidth: '100%',
    chatStartHeight: '32px',
    chatMaxHeight: '70vh',
    chatHeight,
    updateChatHeight: setChatHeight,
    chatWidth,
    updateChatWidth: setChatWidth,
    lastChatButtonClicked,
    updateLastChatButtonClicked: setLastChatButtonClicked,
    lastChatButtonClickedTime,
    updateLastChatButtonClickedTime: setLastChatButtonClickedTime,
    lastChatBorderClicked,
    updateLastChatBorderClicked: setLastChatBorderClicked,
    lastChatBorderClickedTime,
    updateLastChatBorderClickedTime: setLastChatBorderClickedTime,
  }

  return <UIContext.Provider value={value}>{children}</UIContext.Provider>
}
