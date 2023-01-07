package contract

const IDKey = "ygo:id"

type IDService interface {
	NewID() string
}
