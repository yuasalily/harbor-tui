package pages

type ID int

const (
	PageOverview ID = iota
	PageImages
)

func (p ID) Title() string {
	switch p {
	case PageOverview:
		return "Overview"
	case PageImages:
		return "Images"
	default:
		return "Unknown"
	}
}

func All() []ID { return []ID{PageOverview, PageImages} }

func FromTitle(title string) ID {
	for _, id := range All() {
		if id.Title() == title {
			return id
		}
	}
	panic("unknown page title: " + title)
}
