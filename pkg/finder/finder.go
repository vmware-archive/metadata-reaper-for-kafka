package finder

type Finder interface {
	GetTopics() ([]string, error)
}
