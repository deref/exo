package api

// Max length of log messages. Size chosen to fit within a UDP packet (max size
// 64k) and reserve space for Syslog event headers. The message itself always
// includes a newline terminator as well.
const MaxMessageSize = 48 * 1024
