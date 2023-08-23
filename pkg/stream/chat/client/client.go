package pkg

// SOLID Principle: Single Responsibility Principle (SRP)
type MessageHandler interface {
	HandleMessage(msg Message)
}
