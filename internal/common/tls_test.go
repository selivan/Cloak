package common

import (
	"net"
	"testing"
)

func BenchmarkTLSConn_Write(b *testing.B) {
	const bufSize = 16 * 1024

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Fatal(err)
	}
	defer listener.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			b.Error(err)
			return
		}
		defer conn.Close()

		readBuf := make([]byte, bufSize*2)
		for {
			_, err := conn.Read(readBuf)
			if err != nil {
				return
			}
		}
	}()

	defer func() {
		listener.Close()
		<-done
	}()

	data := make([]byte, bufSize)

	discardConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		b.Fatal(err)
	}

	tlsConn := NewTLSConn(discardConn)
	defer tlsConn.Close()

	b.SetBytes(bufSize)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := tlsConn.Write(data); err != nil {
				b.Error(err)
				return
			}
		}
	})
}
