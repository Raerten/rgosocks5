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
	connCnt   uint64
	readBite  uint64
	writeBite uint64
	auth      string
}

type responseStat struct {
	Status    string `json:"status"`
	ConnCount uint64 `json:"connCount"`
	ReadBite  uint64 `json:"readBite"`
	WriteBite uint64 `json:"writeBite"`
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

	s.connOpen()

	return Conn{
		conn,
		s.connRead,
		s.connWrite,
		s.connClose,
	}, nil
}

func (s *Stat) connOpen() {
	atomic.AddUint64(&s.connCnt, 1)
}

func (s *Stat) connClose() {
	var delta uint64 = 1
	atomic.AddUint64(&s.connCnt, ^(delta - 1))
}

func (s *Stat) connWrite(cnt int) {
	atomic.AddUint64(&s.writeBite, uint64(cnt))
}

func (s *Stat) connRead(cnt int) {
	atomic.AddUint64(&s.readBite, uint64(cnt))
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
			ConnCount: atomic.LoadUint64(&s.connCnt),
			ReadBite:  atomic.LoadUint64(&s.readBite),
			WriteBite: atomic.LoadUint64(&s.writeBite),
		})

		if err == nil {
			_, _ = writer.Write(response)

			return
		}
	}

	writer.WriteHeader(404)
	_, _ = writer.Write([]byte("Host not found"))
}
