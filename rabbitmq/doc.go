/*
NewRabbitMQ returns a new initialized connection struct.
It will manage the active connection in the background.
Connection should be closed in order to shut it down gracefully.

	func example() {
		user := "guest"
		pass := "guest"
		host := "localhost"
		port := "5672"
		consumer := "user-service" //  Unique name for each consumer used to sign messages.

		// You can specify your own config or use DefaultConfig() instead.
		config := Config{
			QueueSize:         100,
			MaxWorkers:        50,
			ReconnectInterval: time.Second * 2,
			MaxRequests:       5,
			ClearInterval:     time.Second * 5,
			ClosedTimeout:     time.Second * 5,
		}

		// Logger and tracer are optional.
		rabbit := rabbitmq.NewRabbitMQ(consumer, user, pass, host, port, config, WithLogger(customLogger))
		defer rabbit.Close()
	}
*/
package rabbitmq
