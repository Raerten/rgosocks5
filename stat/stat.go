package stat

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

type Stat struct {
	sync.RWMutex
	connCount uint64
	auth      string
}

type responseStat struct {
	Status    string `json:"status"`
	ConnCount uint64 `json:"connCount"`
}

func NewStat(enabled bool, address string, auth string) *Stat {
	stat := &Stat{
		auth: auth,
	}
	if enabled {
		slog.Info("Starting Status server", "address", address)
		go stat.listenAndServe(address)
	} else {
		slog.Info("Not starting Status server due to config")
	}
	return stat
}

func (s *Stat) Dial(ctx context.Context, network, address string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, network, address)
	if err != nil {
		return conn, err
	}

	s.OpenConn()

	return Conn{
		conn,
		s.CloseConn,
	}, nil
}

func (s *Stat) ConnCount() uint64 {
	return atomic.LoadUint64(&s.connCount)
}

func (s *Stat) OpenConn() {
	atomic.AddUint64(&s.connCount, 1)
}

func (s *Stat) CloseConn() {
	var delta uint64 = 1
	atomic.AddUint64(&s.connCount, ^(delta - 1))
}

func (s *Stat) listenAndServe(address string) {
	if err := http.ListenAndServe(address, s); err != nil {
		slog.Error("Status ListenAndServe", "err", err)
	}
}

func (s *Stat) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if s.auth != "" && ("Bearer "+s.auth) != request.Header.Get("Authorization") {
		writer.WriteHeader(401)
		return
	}

	if request.RequestURI == "/status" {
		writer.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(responseStat{
			Status:    "ok",
			ConnCount: s.ConnCount(),
		})

		if err == nil {
			_, _ = writer.Write(response)

			return
		}
	}

	writer.WriteHeader(404)
	_, _ = writer.Write([]byte("Host not found"))
}
