package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FileRepository struct {
	folderName string
	entityName string
}

func Initialize(folderName string, entityName string) (FileRepository, error) {
	p := FileRepository{folderName: folderName, entityName: entityName}
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		err := os.Mkdir(folderName, 0755)
		return p, err
	}
	return p, err
}

func (p FileRepository) getFileName() string {
	return fmt.Sprintf("%s/%s.json", p.folderName, p.entityName)
}

func (p FileRepository) Persist(data interface{}) error {
	bytes, _ := json.Marshal(data)
	err := ioutil.WriteFile(p.getFileName(), bytes, 0644)
	return err
}

func (p FileRepository) ReadAll(data interface{}) error {
	fileName := p.getFileName()
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		jsonFile, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer jsonFile.Close()
		err = json.NewDecoder(jsonFile).Decode(data)
		return err
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
