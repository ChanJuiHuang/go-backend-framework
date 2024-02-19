package clickhouse

type Config struct {
	Addr         []string
	Database     string
	Username     string
	Password     string
	MaxOpenConns int
	MaxIdleConns int
}
