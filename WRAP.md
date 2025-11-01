# WRAP.md - Vado Client

## Project Overview
**vado-client** is a Go-based desktop GUI application that communicates with a backend server via gRPC. It provides a visual interface for monitoring server status and managing client-side operations.

## Technology Stack
- **GUI Framework**: [Fyne v2](https://fyne.io/) - Cross-platform GUI toolkit
- **RPC**: gRPC with Protocol Buffers
- **Logging**: Uber Zap
- **Message Queue**: Kafka (kafka-go)
- **Language**: Go 1.25

## Project Structure
```
vado-client/
├── api/              # API definitions and generated code
│   └── proto/        # Protocol Buffer definitions
├── cmd/
│   └── gui/          # Main GUI application entry point
├── internal/         # Private application code
│   ├── app/          # Application context and core logic
│   ├── component/    # UI components
│   ├── config/       # Configuration (colors, etc.)
│   ├── grpc/         # gRPC client implementations
│   └── utils/        # Utility functions
├── go.mod            # Go module dependencies
├── Makefile          # Build and generation commands
└── README.md         # Basic setup instructions
```

## Key Features
- Desktop GUI built with Fyne framework
- Real-time server status monitoring with visual indicator
- gRPC-based client-server communication
- Multi-instance support via APP_ID environment variable
- Structured logging with Zap

## Development Commands

### Generate gRPC Code
```shell
make go-proto
```
Generates Go code from `.proto` files in `api/proto/`.

### Clean Dependencies Cache
```shell
go clean -modcache
```

### View Dependency Tree
```shell
go install github.com/PaulXu-cn/go-mod-graph-chart/gmchart@latest
go mod graph | gmchart
```

## Running the Application

### Single Instance
```shell
go run cmd/gui/main.go
```

### Multiple Instances
Use the `APP_ID` environment variable to run multiple instances:
```shell
APP_ID=instance1 go run cmd/gui/main.go
APP_ID=instance2 go run cmd/gui/main.go
```

## Configuration
- Window size: 450x703 pixels
- Default app ID: `vado-client`
- Custom app ID: `vado-client-{APP_ID}` when APP_ID is set

## Server Connection
The client connects to a gRPC server and performs health checks via the Ping service. Server status is displayed with a color-coded indicator:
- 🟢 Green: Server is running
- 🔴 Red: Server is down
- 🟠 Orange: Status unknown/checking

## Related Projects
- **vado-server**: Backend gRPC server (located at `/home/vadmark/GolandProjects/vado-server`)

## Notes
- GUI text and labels are in Russian
- The application uses structured logging for debugging and monitoring
- gRPC client is initialized on application startup
