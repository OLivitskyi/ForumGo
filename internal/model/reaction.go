package model

type ReactionType int

const (
	Like ReactionType = iota
	Dislike
)

type Reaction struct {
	UserID string
	PostID string
	Type   ReactionType
}
