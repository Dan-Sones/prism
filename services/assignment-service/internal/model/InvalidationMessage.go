package model

type InvalidationAction string

const (
	ActionRemove         InvalidationAction = "REMOVE"
	ActionUpdate         InvalidationAction = "UPDATE"
	ActionFullInvalidate InvalidationAction = "FULL_INVALIDATE"
)

func (a InvalidationAction) String() string {
	return string(a)
}

type InvalidationMessage struct {
	Bucket   int32              `json:"bucket"`
	Flag     string             `json:"flag"`
	Action   InvalidationAction `json:"action"`
	NewValue string             `json:"new_value,omitempty"`
}
