package util

// Callback function type for MessengerTopic
type EventCallbackFunction func(data interface{}, topic *messengerTopic)

// Subscribable messaging bus
// Allows data traffic between distinct parts of code
type messengerTopic struct {
	// Subscribers, mapping names to callback functions
	subscribers map[string]EventCallbackFunction
	// Last data value passed to the topic
	LastData interface{}
}

func New() messengerTopic {
	topic := messengerTopic{
		subscribers: make(map[string]EventCallbackFunction),
	}
	return topic
}

// Subscribe to topic with a given an identifier and a callback function
func (topic *messengerTopic) subscribe(subscriberId string, subscriber EventCallbackFunction) {
	topic.subscribers[subscriberId] = subscriber
}

// Unsubscribe a callback function from the topic, given a identifier
func (topic *messengerTopic) unsubscribe(subscriberId string) {
	delete(topic.subscribers, subscriberId)
}

// Notify all subscribed functions with given data
// Subscribed functions are called concurrently in individual goroutines
func (topic *messengerTopic) notify(data interface{}) {
	topic.LastData = data
	for _, f := range topic.subscribers {
		go f(data, topic)
	}
}
