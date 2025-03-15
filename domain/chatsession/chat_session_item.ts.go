package chatsession

//tygo:emit
var _ = `import { Metadata } from "./domain.generated.ts";
import { UserID } from "./user.generated.ts";

export type ChatSessionItemType =
  | 'user'
  | 'assistant'

export type ChatSessionItem = 
  | UserChatSessionItem
  | AssistantChatSessionItem
`
