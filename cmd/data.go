package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var DataFolderPath = fmt.Sprintf("%s/data", os.Getenv("PWD"))

const ReleaseFileFormat = "%s/%s-%s"

func ensureDataFolder(folderPath string) error {
	fileInfo, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"directory": folderPath,
		}).Info("Data directory does not exist, creating one now...")
		err = os.Mkdir(folderPath, 0755)
		if err != nil {
			return err
		}
		return nil
	} else {
		if fileInfo.IsDir() {
			log.WithFields(log.Fields{
				"directory": folderPath,
			}).Info("Using existing data folder")
			return nil
		} else {
			return fmt.Errorf("\"data\" folder cannot be created due to existing\"data\" file")
		}
	}
}

func appendStringToFile(string string, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	lineDelimitedString := fmt.Sprintf("%s\n", string)
	_, err = file.WriteString(lineDelimitedString)
	if err != nil {
		return err
	}
	return nil
}

// each line in the file will be interpreted as a "true" entry in the resulting map[string]bool
func readMapFromFile(filePath string) (map[string]bool, error) {
	resultingMap := make(map[string]bool)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		resultingMap[line] = true
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return resultingMap, nil
}

//lint:ignore U1000 Ignore unused function
func writeMapToFile(mapToWrite map[string]bool, filePath string) error {
	for key, val := range mapToWrite {
		if val {
			err := appendStringToFile(key, filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
