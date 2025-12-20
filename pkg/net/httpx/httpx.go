package httpx

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	client *http.Client
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

type Request struct {
	Method  string
	Path    string
	Headers http.Header
	Body    []byte
}

type Server struct {
	listener net.Listener
	srv      *http.Server

	mu         sync.RWMutex
	respStatus int
	respBody   []byte
	respType   string

	reqCh chan *Request
}

func NewClient() *Client {
	return &Client{client: &http.Client{}}
}

func (c *Client) Do(req *http.Request) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	}, nil
}

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener:   ln,
		respStatus: http.StatusOK,
		respBody:   []byte{},
		respType:   "",
		reqCh:      make(chan *Request, 64),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)
	s.srv = &http.Server{Handler: mux}

	go func() {
		_ = s.srv.Serve(ln)
	}()

	return s, nil
}

func (s *Server) AddrString() string {
	return s.listener.Addr().String()
}

func (s *Server) Close() error {
	return s.srv.Close()
}

func (s *Server) SetStaticResponse(status int, body []byte, contentType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.respStatus = status
	s.respBody = append([]byte(nil), body...)
	s.respType = contentType
}

func (s *Server) NextRequest() *Request {
	return <-s.reqCh
}

func (s *Server) TryNextRequest() *Request {
	select {
	case req := <-s.reqCh:
		return req
	default:
		return nil
	}
}

func (s *Server) NextRequestTimeout(timeout time.Duration) (*Request, bool) {
	if timeout < 0 {
		req := <-s.reqCh
		return req, true
	}
	if timeout == 0 {
		req := s.TryNextRequest()
		if req == nil {
			return nil, false
		}
		return req, true
	}

	select {
	case req := <-s.reqCh:
		return req, true
	case <-time.After(timeout):
		return nil, false
	}
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()

	req := &Request{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: r.Header.Clone(),
		Body:    body,
	}

	s.enqueueRequest(req)

	s.mu.RLock()
	status := s.respStatus
	respBody := append([]byte(nil), s.respBody...)
	contentType := s.respType
	s.mu.RUnlock()

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(status)
	_, _ = w.Write(respBody)
}

func (s *Server) enqueueRequest(req *Request) {
	select {
	case s.reqCh <- req:
		return
	default:
		select {
		case <-s.reqCh:
		default:
		}
		s.reqCh <- req
	}
}
