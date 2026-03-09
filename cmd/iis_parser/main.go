package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"bzy/deployer/pkg/lms"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// IISSite represents each IIS site
type IISSite struct {
	Name         string      `json:"name"`
	State        string      `json:"state"`
	PhysicalPath string      `json:"physicalPath"`
	Bindings     BindingsRaw `json:"Bindings"`
}

// BindingsRaw handles both string or object
type BindingsRaw struct {
	Value []string
}

// UnmarshalJSON custom unmarshal
func (b *BindingsRaw) UnmarshalJSON(data []byte) error {
	// try object first
	var obj struct {
		Value []string `json:"value"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		b.Value = obj.Value
		return nil
	}

	// fallback: if it's a string, put it as a single-element slice
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	b.Value = []string{s}
	return nil
}

// ParseLMS reads JSON file and returns LMS sites
func ParseLMS(path string) ([]IISSite, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var sites []IISSite
	if err := json.Unmarshal(data, &sites); err != nil {
		return nil, err
	}

	var lmsSites []IISSite
	for _, site := range sites {
		for _, binding := range site.Bindings.Value {
			if strings.Contains(binding, ".nuudelms.mn") {
				lmsSites = append(lmsSites, site)
				break
			}
		}
	}

	return lmsSites, nil
}

func main() {
	path := "sites_103.json" // your JSON file
	sites, err := ParseLMS(path)
	if err != nil {
		panic(err)
	}

	// 3️⃣ Convert to *LMS
	var lmsList []*lms.LMS
	for _, s := range sites {
		lmsList = append(lmsList, &lms.LMS{
			Name:   s.Name,
			Path:   s.PhysicalPath,
			Status: s.State,
		})
	}
	db, err := gorm.Open(sqlite.Open("file:lms.db?cache=shared&_fk=1"), &gorm.Config{})
	// db, err := gorm.Open(sqlite.Open("lms.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&lms.LMS{}); err != nil {
		log.Fatal(err)
	}

	repo := lms.NewGormLMSRepository(db)
	// 4️⃣ Bulk insert into DB
	if err := repo.CreateMany(lmsList); err != nil {
		log.Fatal(err)
	}

	// 5️⃣ Show inserted LMSes
	all, _ := repo.ListAll()
	for _, lms := range all {
		fmt.Printf("LMS: %s, Path: %s,  ID: %s\n",
			lms.Name, lms.Path, lms.ID)
	}
}
