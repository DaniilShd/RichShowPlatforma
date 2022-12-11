package models

type MasterClass struct {
	ID          int
	Name        string
	Duration    int
	Description string
}

type Show struct {
	ID          int
	Name        string
	Duration    int
	Description string
}

type Animation struct {
	ID          int
	Name        string
	Duration    int
	Description string
}

type PartyAndQuest struct {
	ID          int
	Name        string
	Duration    int
	Description string
}
type Other struct {
	ID          int
	Name        string
	Duration    int
	Description string
}

type LeadHero struct {
	ID              int
	ArtistID        int
	Description     string
	HeroID          int
	ArtistFirstName string
	ArtistLastName  string
}

type Child struct {
	ID             int
	Name           string
	Gender         int
	Artist         Artist
	DateOfBirthDay string
	Age            int
}

type Client struct {
	ID          int    `db:"id_client"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	PhoneNumber string `db:"phone_number"`
	Telegram    string `db:"telegram_client"`
}

type Lead struct {
	ID                   int    `db:"id_lead"`
	AmountOfChildren     int    `db:"amount_of_children"`
	AverageAgeOfChildren int    `db:"average_age_of_children"`
	Address              string `db:"address"`
	Date                 string `db:"date"`
	Time                 string `db:"time"`
	ActiveLead           bool   `db:"active_lead"`
	CheckArtists         bool   `db:"check_artists"`
	Confirmed            bool   `db:"confirmed"`
	CheckAssistants      bool   `db:"check_assistants"`
	Description          string `db:"description"`
	Duration             int
	MasterClasses        []MasterClass
	Animations           []Animation
	PartyAndQuests       []PartyAndQuest
	Others               []Other
	Shows                []Show
	Heroes               []LeadHero
	Child                Child
	Client               Client
	Assistants           []Assistant
}
