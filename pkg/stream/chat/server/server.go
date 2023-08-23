package pkg

// SOLID Principle: Single Responsibility Principle (SRP)
type MessageProcessor interface {
	ProcessMessage(msg Message)
}
