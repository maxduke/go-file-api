# Go RESTful API Service

This is a simple RESTful API service implemented in Golang. The service accepts API requests, verifies an authentication key, and writes the request body to a specified file.

## Features

- Accepts command line arguments: `auth_key`, `port`, and `file`
- Verifies the authentication key from the request header
- Writes the request body to the specified file
- Provides colored log outputs for easy debugging
- Includes error handling for critical steps

## Usage

### Prerequisites

- Go 1.16 or higher

### Building the Project

```bash
git clone https://github.com/maxduke/go-file-api.git
cd go-file-api
go build -o api_service
```

### Running the Service

```bash
./api_service -auth_key=<your_auth_key> -port=<port_number> -file=<file_path>
```

- `auth_key`: The authentication key required for API calls
- `port`: The port number the service will listen on
- `file`: The file path where the request body will be written

### Example

```bash
./api_service -auth_key=123456 -port=8080 -file=./request_body.txt
```

### Testing the Service

You can test the service using `curl`:

```bash
curl -X POST http://localhost:8080 -H "Authorization: Bearer 123456" -d "Hello, World!"
```

## Logging

The service uses colored output to log key information:

- **Cyan**: Server start information
- **Yellow**: Unauthorized access attempts
- **Red**: Errors
- **Green**: Successful request processing

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
