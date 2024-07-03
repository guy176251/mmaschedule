package scraper

type Event struct {
	Url          string  `json:"url"`
	Slug         string  `json:"slug"`
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Image        string  `json:"image"`
	Date         int     `json:"date"`
	Fights       []Fight `json:"fights"`
}

type Fight struct {
	Weight   string   `json:"weight"`
	FighterA *Fighter `json:"fighter_a"`
	FighterB *Fighter `json:"fighter_b"`
}

type Fighter struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Country string `json:"country"`
	Link    string `json:"link"`
}

type EventCallback func(e *Event)
type EventScraper func(c EventCallback)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"

var scrapers []EventScraper = []EventScraper{
	scrape_ufc,
	scrape_one,
}

func ScrapeEvents(callback EventCallback) {
	get_tapology := tapology_getter()

	for _, scraper := range scrapers {
		scraper(func(e *Event) {
			for _, fight := range e.Fights {
				fight.FighterA.Link = get_tapology(fight.FighterA.Name)
			}
            callback(e)
		})
	}
}
