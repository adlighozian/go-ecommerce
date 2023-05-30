//go:generate mockery --output=../mocks --name Repositorier
package publisher

type PublisherInterface interface {
	Publish(body interface{}, queueName string) (err error)
}