package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type queryComponent struct {
	key      string
	value    []string
	operator string
}

type LogLevel struct {
	ID    int    `gorm:"primaryKey"`
	Level string `gorm:"unique;not null"`
}

type LogComponent struct {
	ID        uint   `gorm:"primaryKey"`
	Component string `gorm:"unique;not null"`
}

type LogHost struct {
	ID   uint   `gorm:"primaryKey"`
	Host string `gorm:"unique;not null"`
}

type Entry struct {
	gorm.Model
	TimeStamp   time.Time
	LevelID     uint
	ComponentID uint
	HostID      uint
	RequestID   string
	Message     string

	Level     LogLevel     `gorm:"foreignKey:LevelID;references:ID"`
	Component LogComponent `gorm:"foreignKey:ComponentID;references:ID"`
	Host      LogHost      `gorm:"foreignKey:HostID;references:ID"`
}

func (e Entry) String() string {
	if e.TimeStamp.IsZero() {
		return "Empty"
	} else {

		return fmt.Sprintf("%s : %s : %s : %s : %s : %s",
			e.TimeStamp.Format("2006-01-02 15:04:05"),
			e.Level.Level,
			e.Component.Component,
			e.Host.Host,
			e.RequestID,
			e.Message)
	}
}

func CreateDB(dbUrl string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, fmt.Errorf("couldn't open database %v", err)
	}
	return db, nil
}

func InitDb(db *gorm.DB) error {
	db.AutoMigrate(&LogLevel{}, &LogComponent{}, &LogHost{}, &Entry{})
	db.Create(&[]LogLevel{
		{Level: "INFO"},
		{Level: "DEBUG"},
		{Level: "WARN"},
		{Level: "ERROR"},
	})
	db.Create(&[]LogComponent{
		{Component: "auth"},
		{Component: "api-server"},
		{Component: "cache"},
		{Component: "database"},
		{Component: "worker"},
	})

	db.Create(&[]LogHost{
		{Host: "web01"},
		{Host: "web02"},
		{Host: "db01"},
		{Host: "cache01"},
		{Host: "worker01"},
	})

	return nil
}

func AddDb(db *gorm.DB, e Entry) error {

	ctx := context.Background()
	err := gorm.G[Entry](db).Create(ctx, &e)
	if err != nil {
		return err
	}
	return nil
}

func parseQuery(items []string) ([]queryComponent, error) {

	var ret []queryComponent
	pattern := `^(?P<key>[^\s=!<>]+)\s*(?P<operator>=|!=|>=|<=|>|<)\s*(?P<value>.+)$`
	r, _ := regexp.Compile(pattern)

	for _, item := range items {
		item = strings.TrimSpace(item)

		matches := r.FindStringSubmatch(item)
		if matches == nil {
			return nil, fmt.Errorf("inavlid condition %s", item)
		}

		values := strings.Split(matches[r.SubexpIndex("value")], ",")
		component := queryComponent{
			key:      matches[r.SubexpIndex("key")],
			operator: matches[r.SubexpIndex("operator")],
			value:    values,
		}

		ret = append(ret, component)
	}
	return ret, nil

}

func QueryDB(db *gorm.DB, query []string) ([]Entry, error) {
	var ret []Entry

	// Parse the query string
	parsed, err := parseQuery(query)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsed conditions:", parsed)

	q := db

	for _, c := range parsed {

		key := strings.ToLower(c.key)

		switch key {

		case "level":
			var ids []uint
			for _, v := range c.value {
				var l LogLevel
				if err := db.First(&l, "level = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown level : %s", v)
				}
				ids = append(ids, uint(l.ID))
			}

			c.key = "level_id"
			c.value = toStringSlice(ids)

		case "component":
			var ids []uint
			for _, v := range c.value {
				var comp LogComponent
				if err := db.First(&comp, "component = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown component '%s'", v)
				}
				ids = append(ids, comp.ID)
			}
			c.key = "component_id"
			c.value = toStringSlice(ids)

		case "host":
			var ids []uint
			for _, v := range c.value {
				var h LogHost
				if err := db.First(&h, "host = ?", v).Error; err != nil {
					return nil, fmt.Errorf("unknown host '%s'", v)
				}
				ids = append(ids, h.ID)
			}
			c.key = "host_id"
			c.value = toStringSlice(ids)

		}

		if len(c.value) == 1 {
			// single value
			fmt.Printf("Applying condition: %s %s %s\n", c.key, c.operator, c.value[0])
			q = q.Where(fmt.Sprintf("%s %s ?", c.key, c.operator), c.value[0])
		} else {
			//multiple values and operator is !=
			if c.operator == "!=" {
				fmt.Printf("Applying NOT IN condition: %s IN %v\n", c.key, c.value)
				q = q.Where(fmt.Sprintf("%s NOT IN ?", c.key), c.value)
			} else {
				// multi value and operator is =
				fmt.Printf("Applying IN condition: %s IN %v\n", c.key, c.value)
				q = q.Where(fmt.Sprintf("%s IN ?", c.key), c.value)
			}

		}
	}

	q = q.Preload("Level").Preload("Component").Preload("Host")
	// Execute final query
	if err := q.Find(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

// convert slice of foreign keys to string
func toStringSlice(nums []uint) []string {
	s := make([]string, len(nums))
	for i, n := range nums {
		s[i] = fmt.Sprint(n)
	}
	return s
}

func SplitUserFilter(input string) []string {
	var parts []string
	current := ""
	tokens := strings.Fields(input)

	for _, tok := range tokens {
		// If token contains an operator, then new condition
		if strings.Contains(tok, "=") ||
			strings.Contains(tok, ">=") ||
			strings.Contains(tok, "<=") ||
			strings.Contains(tok, ">") ||
			strings.Contains(tok, "<") {

			// Save previous condition
			if current != "" {
				parts = append(parts, current)
			}
			current = tok
		} else {
			// continuation (timestamps)
			current += " " + tok
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
