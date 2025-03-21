// Code generated by tygo. DO NOT EDIT.

//////////
// source: chat_session.go

export interface ChatSession extends ChatSessionData {
  id: ChatSessionID
  metadata?: Metadata
}
export interface ChatSessionData {
  chatsessionitemids: ChatSessionItemID[]
  userids: unknown /* user.UserID */[]
}

//////////
// source: chat_session_item.ts.go

import { Metadata } from './domain.generated.ts'
import { UserID } from './user.generated.ts'

export type ChatSessionItemType = 'user' | 'assistant'

export type ChatSessionItem = UserChatSessionItem | AssistantChatSessionItem

//////////
// source: chat_session_item_assistant.go

export interface AssistantChatSessionItemData {
  content: string
}
export interface AssistantChatSessionItem
  extends ChatSessionItemBase,
    AssistantChatSessionItemData {
  type: 'assistant'
}

//////////
// source: chat_session_item_base.go

export interface ChatSessionItemBase {
  id: ChatSessionItemID
  type: ChatSessionItemType

  sessionid: ChatSessionID
  metadata?: Metadata
}

//////////
// source: chat_session_item_id.go

export type ChatSessionItemID = number

//////////
// source: chat_session_item_user.go

export interface UserChatSessionItemData {
  /**
   * 	   @minLength 1
   */
  content: string
  userId: UserID
}
export interface UserChatSessionItem
  extends ChatSessionItemBase,
    UserChatSessionItemData {
  type: 'user'
}
