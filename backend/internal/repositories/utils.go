package repositories

type scanner interface {
	Scan(dest ...any) error
}
