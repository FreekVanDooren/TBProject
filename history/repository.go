package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func setupRepository() Memories {
	memories, err := readFromFile()
	if err != nil {
		panic(err)
	}
	if memories != nil {
		return memories
	}
	_, err = os.Stat("data/")
	if os.IsNotExist(err) {
		os.Mkdir("data/", 0755)
	} else if err != nil {
		panic(err)
	}
	return make(map[int]*Memory)
}

func readFromFile() (Memories, error) {
	if _, err := os.Stat("data/history.json"); !os.IsNotExist(err) {
		jsonFile, err := os.Open("data/history.json")
		if err != nil {
			return nil, err
		}
		defer jsonFile.Close()
		var memories Memories
		err = json.NewDecoder(jsonFile).Decode(&memories)
		if err != nil {
			return nil, err
		}
		return memories, nil
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

/*
  Might lead to loss of data if concurrent requests are done.
*/
func persist(m Memories) {
	bytes, _ := json.Marshal(m)
	err := ioutil.WriteFile("data/history.json", bytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
