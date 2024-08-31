# Twenty Six Cloud Go Client

This package provides a Go client for interacting with the Twenty Six Cloud API (formerly known as Aleph.im). It allows developers to easily integrate Twenty Six Cloud functionality into their Go applications.

## Features

- Account management
- Message creation and signing
- File storage
- Aggregate, Post, and Program message handling
- Instance management

## Installation

To install the package, use the following command:

```
go get github.com/bliiitz/go-twentysixcloud
```

## Usage

### Initializing the client

```go
import "github.com/bliiitz/go-twentysixcloud/client"

// Create a new account from a private key
account, err := client.NewTwentySixAccountFromPrivateKey("your_private_key")
if err != nil {
    // Handle error
}

// Initialize the client
twentySixClient := client.NewTwentySixClient(account, "your_channel", "api_url")
```

### Sending messages

```go
// Create and send an aggregate message
content := client.AggregateMessageContent{
    Key: "your_key",
    Content: map[string]string{"Hello": "World"},
}
message, response, err := twentySixClient.CreateAggregate(content)
if err != nil {
    // Handle error
}
```

### Storing files

```go
message, hash, err := twentySixClient.StoreFile("path/to/your/file")
if err != nil {
    // Handle error
}
```

### Creating instances

```go
instanceContent := client.InstanceMessageContent{
    // Fill in instance details
}
message, response, err := twentySixClient.CreateInstance(instanceContent)
if err != nil {
    // Handle error
}
```

## API Reference

For detailed information on available methods and structures, please refer to the Go documentation comments in the source code.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License