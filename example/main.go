package main

import (
	"fmt"
	"github.com/metakeule/fat"
	"github.com/metakeule/fmtdate"
	"github.com/metakeule/templ"
	. "github.com/metakeule/templ/fat"
	"time"
)

type (
	Person struct {
		FirstName        *fat.Field `fat.type:"string" fat.enum:"Peter|Paul|Mary"`
		LastName         *fat.Field `fat.type:"string"`
		Age              *fat.Field `fat.type:"int" fat.default:"32"`
		FieldsOfInterest *fat.Field `fat.type:"slice.string"`
		Points           *fat.Field `fat.type:"slice.float"`
		Votes            *fat.Field `fat.type:"map.int"`
		Birthday         *fat.Field `fat.type:"time"`
		Datings          *fat.Field `fat.type:"slice.time"`
		Meetings         *fat.Field `fat.type:"map.time"`
	}
)

var (
	// prototype of a Person, used to get field infos, like placeholders
	// and to create new
	PERSON  = fat.Proto(&Person{}).(*Person)
	details = templ.New("details").MustAdd("\n--------------\nDETAILS",
		"\nThe first name: ", Placeholder(PERSON.FirstName),
		"\nThe last name: ", Placeholder(PERSON.LastName),
		"\nThe age: ", Placeholder(PERSON.Age),
		"\nThe fields of interest: ", Placeholder(PERSON.FieldsOfInterest),
		"\nThe points: ", Placeholder(PERSON.Points),
		"\nThe votes: ", Placeholder(PERSON.Votes),
		"\nThe birthday: ", Placeholder(PERSON.Birthday),
		"\nThe datings: ", Placeholder(PERSON.Datings),
		"\nThe meetings: ", Placeholder(PERSON.Meetings),
		"\n-------------\n\n").MustParse()
)

// use fat.New to create a new Person and have the field informations
// available via the prototype
func NewPerson() *Person { return fat.New(PERSON, &Person{}).(*Person) }

func main() {
	now := time.Now()
	bday, _ := fmtdate.Parse("DD.MM.YYYY", "02.01.1952")

	peter := NewPerson()
	peter.FirstName.MustSet("Peter")
	peter.Points.MustScan(`[2,3,4.5]`)
	peter.LastName.MustScan("Pan")
	peter.Votes.MustSet(fat.MapType("int", "Mary", 3, "Paul", 2))
	peter.Birthday.Set(bday)
	peter.Datings.Set(fat.Times(now, now.Add(2*time.Hour)))
	peter.Meetings.Set(fat.MapType("time", "Paul", now.Add(5*time.Hour)))
	peter.FieldsOfInterest.MustSet(fat.Strings("cooking", "swimming"))

	paul := NewPerson()
	paul.FirstName.MustSet("Paul")
	paul.LastName.MustSet("Panzer")
	paul.Age.MustSet(42)
	paul.Age.MustScan("53")
	paul.Points.MustSet(fat.Floats(1.0, 2.3))
	paul.Votes.MustScan(`{"Peter": 45}`)
	paul.Meetings.Set(fat.Map(map[string]time.Time{
		"Peter": now.Add(5 * time.Hour),
		"Mary":  now.Add(7 * time.Hour),
	}))

	fmt.Printf("%s: %s is set? %v\n", peter.FirstName.Name(), peter.FirstName, peter.FirstName.IsSet)
	fmt.Printf("%s: %s is set? %v\n", peter.Age.Name(), peter.Age, peter.Age.IsSet)

	fmt.Println(
		details.MustReplace(Setters(peter)...),
	)

	fmt.Printf("%s: %s is  set? %v\n", paul.FirstName.Name(), paul.FirstName, paul.FirstName.IsSet)
	fmt.Printf("%s: %s is set? %v\n", paul.Age.Name(), paul.Age, paul.Age.IsSet)
	fmt.Println(
		details.MustReplace(Setters(paul)...),
	)
}
