package zorkbot

import "fmt"

// Server ...
type Server struct {
	URL  string
	Port int64
}

// NewServer ...
func NewServer(url string, port int64) (*Server, error) {
	return &Server{
		URL:  url,
		Port: port,
	}, nil
}

// String ...
func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.URL, s.Port)
}
