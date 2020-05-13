package common

// Callback is the type of the callback function
type Callback func(data DataRecord)

// CallbackTeardown is a function to be called after last callback
type CallbackTeardown func()
