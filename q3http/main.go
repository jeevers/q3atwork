package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
)

// heavily adapteed from https://github.com/AntoineMeresse/quake3-rcon-go
const BufferSize = 8192

var CommandPrefix = []byte{'\xff', '\xff', '\xff', '\xff'}

type Q3Client struct {
	Host       string
	Port       string
	Password   string
	Connection net.Conn
}

func (c *Q3Client) Connect() error {
	addr := net.JoinHostPort(c.Host, c.Port)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		slog.Error("connect failed", "error", err, "addr", addr)
		return err
	}
	slog.Info("connected", "addr", addr)
	c.Connection = conn
	return nil
}

func (c Q3Client) Send(cmd string) error {
	commandBytes := []byte(cmd)
	fullMessage := append(CommandPrefix, commandBytes...)
	_, err := c.Connection.Write(fullMessage)
	if err != nil {
		slog.Error("send failed", "error", err)
		return err
	}
	return nil
}

func (c Q3Client) Read() (response string, err error) {
	buffer := make([]byte, BufferSize)
	bytesRead, err := c.Connection.Read(buffer)
	if err != nil {
		slog.Error("read failed", "error", err)
		return "", err
	}

	if bytesRead >= 4 {
		slog.Info("non-empty message", "bytes-read", bytesRead)
		message := string(buffer[4:bytesRead])
		return message, nil
	} else {
		slog.Info("empty message", "bytes-read", bytesRead)
		return "", nil
	}
}

func (c Q3Client) Command(cmd string) (response string, err error) {
	sendErr := c.Send(cmd)
	if sendErr != nil {
		slog.Error("error while sending command", "error", sendErr)
		return "", sendErr
	}
	resp, readErr := c.Read()
	if readErr != nil {
		slog.Error("error while reading command response", "error", sendErr)
		return "", readErr
	}
	return resp, nil
}

func (c Q3Client) RconCommand(cmd string) (response string, err error) {
	rconCmd := fmt.Sprintf("rcon %s %s", c.Password, cmd)
	return c.Command(rconCmd)
}

func (c *Q3Client) Disconnect() {
	slog.Info("closing connection")
	err := c.Connection.Close()
	if err != nil {
		slog.Warn("connection already closed!?!?!")
	} else {
		slog.Info("connection closed")
		c.Connection = nil
	}
}

func main() {
	mux := http.NewServeMux()

	rconPassword, passwdInEnv := os.LookupEnv("IOQ3_RCON_PASSWORD")
	if !passwdInEnv {
		slog.Error("IOQ3_RCON_PASSWORD not set, nothing will work! exiting")
		os.Exit(78)
	}
	client := Q3Client{Host: "127.0.0.1", Port: "27960", Password: rconPassword, Connection: nil}
	_ = client.Connect()
	defer client.Disconnect()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := client.RconCommand("serverinfo")
		if err != nil {
			slog.Warn("serverinfo command failed", "error", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			io.WriteString(w, fmt.Sprintf("error: %s", err))
		} else {
			slog.Info("serviceinfo command successful")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "success")
		}
	})

	http.ListenAndServe("0.0.0.0:8088", mux)
}
