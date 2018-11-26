package index

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"time"

	env "github.com/czerasz/go-lambda-sns-to-es/src/env"
)

// IdxVarName is the name of the environment variable
// which is used to overwrite the default
// ElasticSearch index name template
const IdxVarName = "ES_INDEX_TEMPLATE"
const defaultIdxTpl = "{{ .Prefix }}-{{ .Date.Year }}.{{ .Date.Month }}.{{ .Date.Day }}"

type tmpDate struct {
	Year  string
	Month string
	Day   string
}

type tmpData struct {
	Env    map[string]string
	Date   tmpDate
	Prefix string
}

// Generate ElasticSearch index name
func Generate() (string, error) {
	tpl := template.New("index template")

	tpl, err := tpl.Parse(env.GetEnv(IdxVarName, defaultIdxTpl))

	if err != nil {
		return "", err
	}

	allEnvVars := env.AllVars()
	// Do NOT include the INDEX_TEMPLATE
	// environment variable
	delete(allEnvVars, IdxVarName)

	now := time.Now()
	data := tmpData{
		Date: tmpDate{
			Year:  strconv.Itoa(now.Year()),
			Month: fmt.Sprintf("%02d", now.Month()),
			Day:   fmt.Sprintf("%02d", now.Day()),
		},
		Prefix: "sns",
		Env:    allEnvVars,
	}

	var output bytes.Buffer

	if err := tpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}
