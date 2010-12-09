// MACHINE GENERATED.

package _

type server struct {
	uuid         string
	access_log   string
	error_log    string
	chroot       string
	pid_File     string
	default_host int
	name         string
	port         int
}

type host struct {
	id          int
	server_id   int
	maintenance bool
	name        string
	matching    string
}
