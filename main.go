package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/itltf512116/univ/resolve"
)

func main() {
	fmt.Println("resolve the list of all university in the world...")
	client := resolve.NewUnivClient()

	cs, err := client.ResolveAllCountries()
	if err != nil {
		panic(err)
	}

	err = dumpToJson(cs, "./data/country.json")
	if err != nil {
		panic(err)
	}

	us, err := client.ResolveUniversityInUnitedStates()
	if err != nil {
		panic(err)
	}
	dumpToJson(us, fmt.Sprintf("./data/university_%s.json", "UnitedStates"))

	for _, c := range cs {
		fmt.Printf("country:%s\n", c.Name)
		us, err := client.ResolveUniversityByCountry(c)
		if err != nil {
			fmt.Printf("ResolveUniversityByCountry:%s err:%v\n", c.Name, err)
			continue
		}

		time.Sleep(time.Duration(rand.Int31n(3)+1) * time.Second)

		dumpToJson(us, fmt.Sprintf("./data/university_%s.json", c.Name))
	}

}

// dumpToJson save v to fn with json
func dumpToJson(v interface{}, fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
