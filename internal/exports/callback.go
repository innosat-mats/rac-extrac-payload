package exports

// Callback is the type of the callback function
type Callback func(data ExportablePackage)

// CallbackTeardown is a function to be called after last callback
type CallbackTeardown func()
