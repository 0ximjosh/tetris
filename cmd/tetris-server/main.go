package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"tetris"
)

const (
	protocolCellV1 = "ggp.cell.v1"

	capRenderCell    = "render.cell.v1"
	capInputKeyboard = "input.keyboard.v1"
	capScoreReport   = "score.report.v1"

	typeHello  = "hello"
	typeReady  = "ready"
	typeFrame  = "frame"
	typeInput  = "input"
	typeResize = "resize"
	typeScore  = "score"
	typeError  = "error"

	frameFull = "full"
)

type envelope struct {
	Type string `json:"type"`
}

type viewport struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

type hello struct {
	Type     string   `json:"type"`
	Protocol string   `json:"protocol"`
	Viewport viewport `json:"viewport"`
}

type ready struct {
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	TargetFPS    int      `json:"targetFps"`
	Capabilities []string `json:"capabilities"`
}

type frame struct {
	Type   string        `json:"type"`
	Seq    int           `json:"seq"`
	Mode   string        `json:"mode"`
	Status string        `json:"status,omitempty"`
	Cells  []tetris.Cell `json:"cells"`
}

type score struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

type input struct {
	Type string `json:"type"`
	Kind string `json:"kind"`
	Key  string `json:"key"`
}

type resize struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

type gameError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type session struct {
	conn      *websocket.Conn
	model     tetris.Model
	seq       int
	lastScore int64
	scoreSent bool
	mu        sync.Mutex
	writeMu   sync.Mutex
	done      chan struct{}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}

func main() {
	port := flag.String("port", env("PORT", "8080"), "Port for HTTP connections")
	host := flag.String("host", env("HOST", "0.0.0.0"), "Host address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/ggp", handleGGP)

	server := &http.Server{Addr: net.JoinHostPort(*host, *port), Handler: mux}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("tetris GGP server listening on %s", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server failed: %v", err)
			done <- nil
		}
	}()

	<-done
	log.Print("stopping server")
}

func handleGGP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	_, payload, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var hello hello
	if err := json.Unmarshal(payload, &hello); err != nil || hello.Type != typeHello || hello.Protocol != protocolCellV1 {
		_ = conn.WriteJSON(gameError{Type: typeError, Message: "expected ggp.cell.v1 hello"})
		return
	}

	s := &session{conn: conn, done: make(chan struct{})}
	s.model.UpdateDims(maxInt(hello.Viewport.Cols, 20), maxInt(hello.Viewport.Rows, 8))
	s.model.Reset()

	if err := s.write(ready{Type: typeReady, Title: "Tetris", TargetFPS: 4, Capabilities: []string{capRenderCell, capInputKeyboard, capScoreReport}}); err != nil {
		return
	}
	if err := s.sendFrame(); err != nil {
		return
	}
	if err := s.sendScoreIfChanged(); err != nil {
		return
	}

	go s.tickLoop()
	s.readLoop()
}

func (s *session) readLoop() {
	defer s.close()
	for {
		_, payload, err := s.conn.ReadMessage()
		if err != nil {
			return
		}

		var envelope envelope
		if err := json.Unmarshal(payload, &envelope); err != nil {
			continue
		}

		s.mu.Lock()
		sendFrame := true
		switch envelope.Type {
		case typeInput:
			var input input
			if err := json.Unmarshal(payload, &input); err == nil {
				s.model.HandleGGPInput(input.Key)
			}
		case typeResize:
			var resize resize
			if err := json.Unmarshal(payload, &resize); err == nil {
				s.model.UpdateDims(maxInt(resize.Cols, 20), maxInt(resize.Rows, 8))
			}
		default:
			sendFrame = false
		}
		s.mu.Unlock()

		if sendFrame {
			_ = s.sendFrame()
			_ = s.sendScoreIfChanged()
		}
	}
}

func (s *session) tickLoop() {
	for {
		s.mu.Lock()
		interval := s.model.FrameInterval()
		s.mu.Unlock()

		timer := time.NewTimer(interval)
		select {
		case <-timer.C:
			s.mu.Lock()
			s.model.Advance()
			s.mu.Unlock()
			if err := s.sendFrame(); err != nil {
				s.close()
				return
			}
			if err := s.sendScoreIfChanged(); err != nil {
				s.close()
				return
			}
		case <-s.done:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			return
		}
	}
}

func (s *session) sendFrame() error {
	s.mu.Lock()
	s.seq++
	frame := frame{Type: typeFrame, Seq: s.seq, Mode: frameFull, Cells: s.model.Cells()}
	s.mu.Unlock()
	return s.write(frame)
}

func (s *session) sendScoreIfChanged() error {
	s.mu.Lock()
	current := s.model.Score()
	if s.scoreSent && current == s.lastScore {
		s.mu.Unlock()
		return nil
	}
	s.scoreSent = true
	s.lastScore = current
	s.mu.Unlock()
	return s.write(score{Type: typeScore, Value: current})
}

func (s *session) write(value any) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	select {
	case <-s.done:
		return errors.New("session closed")
	default:
		return s.conn.WriteJSON(value)
	}
}

func (s *session) close() {
	select {
	case <-s.done:
	default:
		close(s.done)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
