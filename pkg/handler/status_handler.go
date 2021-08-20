package handler

// StatusHandler check whetever the domain is running
// or already expired from Ngrok
type StatusHandler interface {
	Running(domain string) (bool, error)
}
