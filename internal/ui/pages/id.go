package pages

type ID int

const (
	PageOverview ID = iota
	PageImages
	PageContainers
)

type Meta struct {
	ID    ID
	Title string
	Key   string
}

func Metas() []Meta {
	return []Meta{
		{PageOverview, "Overview", "1"},
		{PageImages, "Images", "2"},
		{PageContainers, "Containers", "3"},
	}
}

func (p ID) Title() string {
	for _, m := range Metas() {
		if m.ID == p {
			return m.Title
		}
	}
	panic("unknown page id: Title not found")
}

func (p ID) Key() string {
	for _, m := range Metas() {
		if m.ID == p {
			return m.Key
		}
	}
	panic("unknown page id: Key not found")
}

func FromTitle(title string) ID {
	for _, m := range Metas() {
		if m.Title == title {
			return m.ID
		}
	}
	panic("unknown page title: " + title)
}

func FromKey(key string) ID {
	for _, m := range Metas() {
		if m.Key == key {
			return m.ID
		}
	}
	panic("unknown page key: " + key)
}
