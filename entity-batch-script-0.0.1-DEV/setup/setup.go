package setup

import (
	. "entity-batch-script/utils/constants"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey"
)

var questions = []*survey.Question{
	{
		Name:     "entitiyName",
		Prompt:   &survey.Input{Message: "Name of entity:"},
		Validate: survey.Required,
	},
	{
		Name: "db",
		Prompt: &survey.Select{
			Message: "What database are you using?",
			Options: []string{"mysql", "postgresql // not supported yet"},
		},
		Validate: survey.Required,
	},
	{
		Name:     "dbHost",
		Prompt:   &survey.Input{Message: "Database host:"},
		Validate: survey.Required,
	},
	{
		Name:     "dbPort",
		Prompt:   &survey.Input{Message: "Database port:"},
		Validate: survey.Required,
	},
	{
		Name:     "dbUser",
		Prompt:   &survey.Input{Message: "Database user:"},
		Validate: survey.Required,
	},
	{
		Name:     "dbPass",
		Prompt:   &survey.Input{Message: "Database password:"},
		Validate: survey.Required,
	},
	{
		Name:     "dbName",
		Prompt:   &survey.Input{Message: "Database name:"},
		Validate: survey.Required,
	},
	{
		Name:     "tableName",
		Prompt:   &survey.Input{Message: "Table name:"},
		Validate: survey.Required,
	},
}

type answers struct {
	EntitiyName string
	Db          string
	DbHost      string
	DbPort      string
	DbUser      string
	DbPass      string
	DbName      string
	TableName   string
}

// RunSetup ask the questions and return the answers
func RunSetup() (string, error) {
	var answers answers

	err := survey.Ask(questions, &answers)
	variables := []string{
		EntityName,
		Db,
		DbHost,
		DbPort,
		DbUser,
		DbPass,
		DbName,
		TableName,
		"",
	}

	envVars := fmt.Sprintf(
		strings.Join(variables, "=%v\n"),
		answers.EntitiyName,
		answers.Db,
		answers.DbHost,
		answers.DbPort,
		answers.DbUser,
		answers.DbPass,
		answers.DbName,
		answers.TableName,
	)

	return envVars, err
}

// EntityDotEnvDoesNotExist returns true if .env exists
func EntityDotEnvDoesNotExist() bool {
	_, err := os.Stat(EntityDotEnv)

	return os.IsNotExist(err)
}

// UpdateExistingSetup ask if user want to update existing .env.entity
func UpdateExistingSetup() bool {
	answer := false
	prompt := &survey.Confirm{
		Message: "Would you like to update",
	}
	survey.AskOne(prompt, &answer, survey.WithIcons(func(icons *survey.IconSet) {
		// you can set any icons
		icons.Question.Text = "#"
		// for more information on formatting the icons, see here: https://github.com/mgutz/ansi#style-format
		icons.Question.Format = "green+hb"
	}))

	return answer
}
