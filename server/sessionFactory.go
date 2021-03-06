package server

import (
	"blabu/c2cService/parser"
	"blabu/c2cService/stat"
	"net"
	"time"
)

const minHeaderSize = 128

// StartNewSession - инициализирует все и стартует сессию
func StartNewSession(conn net.Conn, dT time.Duration, st *stat.Statistics) {
	req := make([]byte, minHeaderSize)
	conn.SetReadDeadline(time.Now().Add(dT))
	if n, err := conn.Read(req); err == nil {
		if p, err := parser.InitParser(req[:n]); err == nil && p != nil {
			st.NewConnection() // Регистрируем новое соединение
			start := time.Now()
			s := BidirectSession{
				Duration: dT,
				Tm:       time.NewTimer(dT),
				netReq:   req,
				logic:    CreateReadWriteMainLogic(p, dT),
			}
			s.Run(conn, p)
			s.logic.Close()
			st.CloseConnection()
			st.SetConnectionTime(time.Since(start))
		}
	}
	conn.Close()
}
