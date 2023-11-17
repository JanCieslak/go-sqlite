package sqlite

type PageType byte

const (
	InteriorIndexPageType PageType = 2
	InteriorTablePageType PageType = 5
	LeafIndexPageType     PageType = 10
	LeafTablePageType     PageType = 13
)

type Header struct {
	Type PageType
}
