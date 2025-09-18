
---

## Packet Structure

The application uses a custom binary packet format, `LRTPIDSPacket`, designed to be lightweight and efficient for transmitting PIDS data. This structure allows for precise control over the data payload and demonstrates the ability to define a protocol at the application layer.

- **TransactionId (uint16)**: A unique identifier for each packet transaction.
- **Flags (bit-field)**: A single byte with multiple flags to indicate the packet's purpose:
    - `IsAck`: Acknowledgment for a received packet.
    - `IsNewTrain`: Indicates a new train entry.
    - `IsUpdateTrain`: Indicates an update to an existing train.
    - `IsDeleteTrain`: Indicates a train's removal from the display.
    - `IsTrainArriving`: Notifies of a train's arrival.
    - `IsTrainDeparting`: Notifies of a train's departure.
- **TrainNumber (uint16)**: A unique identifier for a specific train.
- **Destination (string)**: The train's destination, with a variable length.

---

## Technical Highlights and Learning Outcomes

This project provided valuable experience in several key areas of network programming:

- **Implementing QUIC**: Hands-on experience with **QUIC's core features**, including stream multiplexing and built-in security. Unlike TCP, QUIC allows multiple independent streams of data to be sent over a single connection, avoiding head-of-line blocking and improving performance for a real-world application like a PIDS.
- **Custom Protocol Design**: Designed and implemented a custom, lightweight binary protocol, `LRTPIDSPacket`. This involved careful consideration of **packet serialization and deserialization** to ensure data integrity and interoperability between the publisher and subscriber.
- **Concurrency in Go**: The subscriber (server) was designed to handle multiple incoming QUIC streams concurrently using Go routines. This efficiently manages simultaneous connections and ensures the system remains responsive, even under high load.
- **Error Handling and Reliability**: The system includes a basic acknowledgment mechanism to confirm packet delivery, a fundamental aspect of reliable network communication. This ensures that critical train information is not lost during transmission.

---

## How to Run

### Prerequisites

- **Go (v1.21+)**: The project is built using the Go programming language.
- **QUIC-Go Library**: The `github.com/quic-go/quic-go` library is used for the QUIC implementation.

---

### Subscriber (Server / Layar PIDS)

1. Navigate to the project directory:

    ```bash
    cd h01-source-code
    ```

2. Run the subscriber:

    ```bash
    go run ./project/subscriber
    ```

The subscriber will start listening for incoming QUIC connections from publishers and display received messages, such as:

