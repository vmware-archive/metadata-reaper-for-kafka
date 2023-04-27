package deleter

type Deleter interface {
	DeleteTopics(topics []string) error
}
