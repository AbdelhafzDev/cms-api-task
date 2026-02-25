package messaging

const (
	EventUserCreated   = "user.created"
	EventUserUpdated   = "user.updated"
	EventUserDeleted   = "user.deleted"
	EventUserOnline    = "user.online"
	EventUserOffline   = "user.offline"

	EventUserLoggedIn  = "auth.login"
	EventUserLoggedOut = "auth.logout"
	EventTokenRefreshed = "auth.token_refreshed"

	EventConversationCreated = "conversation.created"
	EventConversationUpdated = "conversation.updated"
	EventConversationDeleted = "conversation.deleted"
	EventParticipantAdded    = "conversation.participant_added"
	EventParticipantRemoved  = "conversation.participant_removed"

	EventMessageSent     = "message.sent"
	EventMessageUpdated  = "message.updated"
	EventMessageDeleted  = "message.deleted"
	EventMessageRead     = "message.read"
	EventTypingStarted   = "message.typing_started"
	EventTypingStopped   = "message.typing_stopped"
)

const (
	ChannelUsers         = "users"
	ChannelAuth          = "auth"
	ChannelConversations = "conversations"
	ChannelMessages      = "messages"
)

type UserEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email,omitempty"`
	Action string `json:"action"`
}

type AuthEvent struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type ConversationEvent struct {
	ConversationID string   `json:"conversation_id"`
	UserID         string   `json:"user_id,omitempty"`
	Participants   []string `json:"participants,omitempty"`
}

type MessageEvent struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content,omitempty"`
}

type TypingEvent struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"user_id"`
	IsTyping       bool   `json:"is_typing"`
}
