package data

type DBConn interface {
}

func NewDBConn(dsn string) (DBConn, error) {

	return &postgresDbConn{}, nil
}

type postgresDbConn struct {
}
