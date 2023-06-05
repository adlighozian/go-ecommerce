package publisher

type Publisher interface {
	Public(req any, queueName string) error
}
