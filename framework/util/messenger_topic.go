package util

// Callback function type for MessengerTopic
type EventCallbackFunction func(data interface{}, topic *MessengerTopic)

// Subscribable messaging bus
// Allows data traffic between distinct parts of code
type MessengerTopic struct {
	// Subscribers, mapping names to callback functions
	subscribers map[string]EventCallbackFunction
	// Last data value passed to the topic
	LastData interface{}
}

// Subscribe to topic with a given an identifier and a callback function
func (topic *MessengerTopic) subscribe(subscriberId string, subscriber EventCallbackFunction) {
	topic.subscribers[subscriberId] = subscriber
}

// Unsubscribe a callback function from the topic, given a identifier
func (topic *MessengerTopic) unsubscribe(subscriberId string) {
	delete(topic.subscribers, subscriberId)
}

// Notify all subscribed functions with given data
// Subscribed functions are called concurrently in individual goroutines
func (topic *MessengerTopic) notify(data interface{}) {
	topic.LastData = data
	for _, f := range topic.subscribers {
		go f(data, topic)
	}
}
