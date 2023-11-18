package sqlite

type PageType byte

const (
	InteriorIndexPageType PageType = 2
	InteriorTablePageType PageType = 5
	LeafIndexPageType     PageType = 10
	LeafTablePageType     PageType = 13
)

type Header struct {
}

type Page struct {
	Type PageType
}

type Database struct {
	Header Header
	Pages  []Page
}

func ParseDatabase(filename string) (*Database, error) {

	return nil, nil
}
