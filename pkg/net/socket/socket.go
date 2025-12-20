package socket

import (
	"net"
)

type TCPConn struct {
	conn *net.TCPConn
}

type TCPListener struct {
	ln *net.TCPListener
}

type UDPConn struct {
	conn *net.UDPConn
}

func DialTCP(addr string) (*TCPConn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &TCPConn{conn: conn}, nil
}

func ListenTCP(addr string) (*TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	return &TCPListener{ln: ln}, nil
}

func (l *TCPListener) Accept() (*TCPConn, error) {
	conn, err := l.ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return &TCPConn{conn: conn}, nil
}

func (l *TCPListener) AddrString() string {
	return l.ln.Addr().String()
}

func (l *TCPListener) Close() error {
	return l.ln.Close()
}

func (c *TCPConn) Read(buf []byte) (int, error) {
	return c.conn.Read(buf)
}

func (c *TCPConn) Write(buf []byte) (int, error) {
	return c.conn.Write(buf)
}

func (c *TCPConn) LocalAddrString() string {
	return c.conn.LocalAddr().String()
}

func (c *TCPConn) RemoteAddrString() string {
	return c.conn.RemoteAddr().String()
}

func (c *TCPConn) Close() error {
	return c.conn.Close()
}

func ListenUDP(addr string) (*UDPConn, error) {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	return &UDPConn{conn: conn}, nil
}

func DialUDP(addr string) (*UDPConn, error) {
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return &UDPConn{conn: conn}, nil
}

func (c *UDPConn) ReadFrom(buf []byte) (int, *net.UDPAddr, error) {
	return c.conn.ReadFromUDP(buf)
}

func (c *UDPConn) WriteTo(buf []byte, addr *net.UDPAddr) (int, error) {
	return c.conn.WriteToUDP(buf, addr)
}

func (c *UDPConn) Write(buf []byte) (int, error) {
	return c.conn.Write(buf)
}

func (c *UDPConn) LocalAddrString() string {
	return c.conn.LocalAddr().String()
}

func (c *UDPConn) RemoteAddrString() string {
	addr := c.conn.RemoteAddr()
	if addr == nil {
		return ""
	}
	return addr.String()
}

func (c *UDPConn) Close() error {
	return c.conn.Close()
}
