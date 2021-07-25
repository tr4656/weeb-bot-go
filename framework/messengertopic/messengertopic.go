package messengertopic

// Callback function type for MessengerTopic
type EventCallbackFunction func(data *Message, topic *messengerTopic)

// Subscribable messaging bus
// Allows data traffic between distinct parts of code
type messengerTopic struct {
	// Subscribers, mapping names to callback functions
	subscribers map[string]EventCallbackFunction
	// Last data value passed to the topic
	LastData interface{}
}

// Generic message data interface
type Message interface{}

// Create new instance of messengerTopic
func New() *messengerTopic {
	topic := messengerTopic{
		subscribers: make(map[string]EventCallbackFunction),
	}
	return &topic
}

// Subscribe to topic with a given an identifier and a callback function
func (topic *messengerTopic) Subscribe(subscriberId string, subscriber EventCallbackFunction) {
	topic.subscribers[subscriberId] = subscriber
}

// Unsubscribe a callback function from the topic, given a identifier
func (topic *messengerTopic) Unsubscribe(subscriberId string) {
	delete(topic.subscribers, subscriberId)
}

// Notify all subscribed functions with given data
// Subscribed functions are called concurrently in individual goroutines
func (topic *messengerTopic) Notify(data *Message) {
	topic.LastData = data
	for _, f := range topic.subscribers {
		go f(data, topic)
	}
}
