package main

import (
	"encoding/json"
	"entity-batch-script/config"
	"entity-batch-script/models/entity"
	"entity-batch-script/models/identity"
	"entity-batch-script/setup"
	. "entity-batch-script/utils/constants"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/joho/godotenv"
	. "github.com/logrusorgru/aurora"
)

const topicsFile = "topics.json"

var entityDetails map[string]string
var identityDetails map[string]string

var topics = make(map[string]interface{})

func readTopics() error {
	file, err := ioutil.ReadFile(topicsFile)

	if err != nil {
		return fmt.Errorf("Error in reading %v: %v", topicsFile, err)
	}

	err = json.Unmarshal(file, &topics)

	if err != nil {
		return fmt.Errorf("Error in unmarshal: %v", err)
	}

	return nil
}

func createEntityDotEnv() error {
	envVars, setupError := setup.RunSetup()

	if setupError != nil {
		return fmt.Errorf("Error in setting up: %v", setupError)
	}

	env, unmarshalError := godotenv.Unmarshal(envVars)

	if unmarshalError != nil {
		return fmt.Errorf("Error in unmarshal: %v", setupError)
	}

	writtingError := godotenv.Write(env, EntityDotEnv)

	if writtingError != nil {
		return fmt.Errorf("Error in create/write to %v: %v", EntityDotEnv, setupError)
	}

	return nil
}

func init() {
	var answer = false

	if setup.EntityDotEnvDoesNotExist() == false {
		result, err := godotenv.Read(EntityDotEnv)

		if err != nil {
			log.Fatalln("Error reading env file:", err)
		}

		fmt.Println(Bold(Green("There is an existing setup.\n")))
		entityDetails = make(map[string]string)
		for key, value := range result {
			entityDetails[key] = value
			fmt.Printf("[%v]: %v\n", Yellow(key), Bold(BrightYellow(value)))
		}
		fmt.Println()
		answer = setup.UpdateExistingSetup()
		fmt.Println()
	}

	if answer || setup.EntityDotEnvDoesNotExist() {
		fmt.Println(Bold(Green("Creating a configuration file .. ..\n")))
		setupError := createEntityDotEnv()

		if setupError != nil {
			log.Fatalln(setupError)
		}

		fmt.Println(Bold(Green("\nSuccessfully created a configuration file!\n")))

		result, err := godotenv.Read(EntityDotEnv)

		if err != nil {
			log.Fatalln("Error reading env file:", err)
		}

		entityDetails = make(map[string]string)
		for key, value := range result {
			entityDetails[key] = value
		}
	}

	readTopicsError := readTopics()

	if readTopicsError != nil {
		log.Fatalln(readTopicsError)
	}
}

func main() {
	identityDetails, err := godotenv.Read()
	if err != nil {
		log.Fatalln("Error reading env file:", err)
	}

	identityDb, connectionErr := config.CreateSQLdatabase(identityDetails)

	if connectionErr != nil {
		log.Fatalln(connectionErr)
	}

	defer identityDb.Close()

	identityModel := identity.Model{
		Db:        identityDb,
		DbName:    identityDetails[DbName],
		TableName: identityDetails[TableName],
	}

	tableExist, existErr := identityModel.TableExist()

	if existErr != nil {
		log.Fatalln(existErr)
	}

	entityDb, connectionErr1 := config.CreateSQLdatabase(entityDetails)

	if connectionErr1 != nil {
		log.Fatalln(connectionErr1)
	}

	defer entityDb.Close()

	entityModel := entity.Model{
		Db:        entityDb,
		TopicMap:  topics,
		TableName: entityDetails[TableName],
	}

	fmt.Println(Bold(Green(fmt.Sprintf("Fetching rows from %s table .. ..", entityDetails[TableName]))))
	topicDataMatrix, getErr := entityModel.GetRows()

	if getErr != nil {
		log.Fatalln(getErr)
	}

	fmt.Println(Bold(Green("Successfully fetched rows!\n")))

	if tableExist {
		dropErr := identityModel.DropTable()

		if dropErr != nil {
			log.Fatalln(dropErr)
		}
	}

	firstElement := topicDataMatrix[0]
	var columns []string

	for i := 0; i < len(firstElement); i++ {
		columns = append(columns, firstElement[i].Topic)
	}
	createErr := identityModel.CreateTable(columns)

	if createErr != nil {
		log.Fatalln(createErr)
	}

	fmt.Println(Bold(Green("Hashing and copying rows .. ..")))

	bulkInsertErr := identityModel.BulkInsert(topicDataMatrix)

	if bulkInsertErr != nil {
		log.Fatalln(bulkInsertErr)
	}

	fmt.Println(Bold(Green("Successfully hashed and copied rows!")))
}
