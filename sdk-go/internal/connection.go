package internal

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Connection manages the TCP connection to a broker
type Connection struct {
	brokerAddr string
	conn       net.Conn
	mutex      sync.RWMutex
	closed     bool

	// Write channel for buffering writes
	writeChan chan []byte
	readChan  chan []byte

	// Error channel
	errChan chan error
}

// NewConnection creates a new connection to a broker
func NewConnection(brokerAddr string) (*Connection, error) {
	// Attempt to connect with timeout
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.Dial("tcp", brokerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", brokerAddr, err)
	}

	c := &Connection{
		brokerAddr: brokerAddr,
		conn:       conn,
		writeChan:  make(chan []byte, 100),
		readChan:   make(chan []byte, 100),
		errChan:    make(chan error, 10),
	}

	// Start read/write loops
	go c.readLoop()
	go c.writeLoop()

	return c, nil
}

// Write sends data to the broker
func (c *Connection) Write(data []byte) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return fmt.Errorf("connection is closed")
	}

	select {
	case c.writeChan <- data:
		return nil
	default:
		return fmt.Errorf("write buffer full")
	}
}

// Read receives data from the broker
func (c *Connection) Read() ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("connection is closed")
	}

	select {
	case data := <-c.readChan:
		return data, nil
	case err := <-c.errChan:
		return nil, err
	}
}

// readLoop continuously reads from the socket
func (c *Connection) readLoop() {
	buffer := make([]byte, 65536) // 64KB buffer

	for {
		c.mutex.RLock()
		if c.closed {
			c.mutex.RUnlock()
			return
		}
		conn := c.conn
		c.mutex.RUnlock()

		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		n, err := conn.Read(buffer)

		if err != nil {
			select {
			case c.errChan <- err:
			default:
			}
			return
		}

		if n > 0 {
			// Copy data to avoid buffer reuse issues
			data := make([]byte, n)
			copy(data, buffer[:n])

			select {
			case c.readChan <- data:
			default:
				// Channel full, skip
			}
		}
	}
}

// writeLoop continuously writes to the socket
func (c *Connection) writeLoop() {
	for {
		c.mutex.RLock()
		if c.closed {
			c.mutex.RUnlock()
			return
		}
		conn := c.conn
		c.mutex.RUnlock()

		select {
		case data := <-c.writeChan:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err := conn.Write(data)
			if err != nil {
				select {
				case c.errChan <- err:
				default:
				}
				return
			}
		}
	}
}

// Close closes the connection
func (c *Connection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	close(c.writeChan)

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

// IsConnected checks if the connection is active
func (c *Connection) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return !c.closed && c.conn != nil
}
