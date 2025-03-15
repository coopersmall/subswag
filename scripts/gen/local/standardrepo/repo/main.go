package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	// ask user for name
	fmt.Println("Enter the name of the domain object you want to create the repo for:")
	var structName string
	fmt.Scanln(&structName)

	for !isTitleCamelCase(structName) {
		fmt.Println()
		fmt.Println("Domain object names must be in title camel case, for example: MyObj")
		fmt.Println("Enter the name of the domain object you want to create the repo for:")
		fmt.Scanln(&structName)
	}

	tmpl, err := template.New("standardRepo").Funcs(template.FuncMap{
		"varName":  variableNameFromStructName,
		"sqlcName": variableNameFromSQLCName,
	}).Parse(standardRepoTemplate)
	if err != nil {
		panic(err)
	}

	var rendered string
	writer := bytes.NewBufferString(rendered)
	tmpl.Execute(writer, TemplateData{StructName: structName})

	// capitalize the first letter of the struct name
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	repoFileName := strings.ToUpper(string(structName[0])) + structName[1:] + "Repo.go"
	repoFilePath := createPath(currentDir, RepoDir, repoFileName)
	err = os.WriteFile(repoFilePath, writer.Bytes(), 0644)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Repo created!")
	}

	currentSchemaBytes, err := os.ReadFile(createPath(currentDir, SQLDir, SchemaFile))
	if err != nil {
		panic(err)
	}

	currentSchema := string(currentSchemaBytes)

	// regex to check if the schema already exists
	re := regexp.MustCompile(`CREATE TABLE ` + dbNameFromStructName(structName))
	if re.MatchString(currentSchema) {
		fmt.Println("Schema already exists... skipping")
	} else {
		tmpl, err = template.New("schema").Funcs(template.FuncMap{
			"dbName": dbNameFromStructName,
		}).Parse(scehmaTemplate)
		if err != nil {
			panic(err)
		}
		writer = bytes.NewBufferString(rendered)
		tmpl.Execute(writer, TemplateData{StructName: structName})
		schemaFilePath := createPath(currentDir, SQLDir, SchemaFile)
		err = os.WriteFile(schemaFilePath, []byte(currentSchema+writer.String()), 0644)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Schema created!")
		}
	}

	currentQueriesBytes, err := os.ReadFile(createPath(currentDir, SQLDir, QueryFile))
	if err != nil {
		panic(err)
	}

	currentQueries := string(currentQueriesBytes)

	// regex to check if the queries already exist
	re = regexp.MustCompile(`-- name: Create` + structName)
	if re.MatchString(currentQueries) {
		fmt.Println("Queries already exist... skipping")
	} else {
		tmpl, err = template.New("queries").Funcs(template.FuncMap{
			"dbName": dbNameFromStructName,
		}).Parse(queriesTemplate)

		if err != nil {
			panic(err)
		}
		writer = bytes.NewBufferString(rendered)
		tmpl.Execute(writer, TemplateData{StructName: structName})
		queriesFilePath := createPath(currentDir, SQLDir, QueryFile)
		err = os.WriteFile(queriesFilePath, []byte(currentQueries+writer.String()), 0644)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("Queries created!")
		}
	}
}

func isTitleCamelCase(s string) bool {
	re := regexp.MustCompile(`^[A-Z][a-z]*([A-Z][a-z]*)*$`)
	return re.MatchString(s)
}

func createPath(strs ...string) string {
	return strings.Join(strs, "/")
}

type TemplateData struct {
	StructName string
}

const (
	RepoDir    = "repos"
	SQLDir     = "sql"
	SchemaFile = "schema.sql"
	QueryFile  = "query.sql"
)

const standardRepoTemplate = `package repos
import (
    "context"
    "database/sql"
    "encoding/json"

    "github.com/coopersmall/subswag/db"
    "github.com/coopersmall/subswag/domain"
)

type {{.StructName}}Repo struct {
    getQueries func() *db.Queries
    userId     domain.UserID
}

var New{{.StructName}}Repo = new{{.StructName}}Repo

func new{{.StructName}}Repo(
    getQueries func() *db.Queries,
    userId domain.UserID,
) *{{.StructName}}Repo {
    return &{{.StructName}}Repo{
        getQueries: getQueries,
        userId:     userId,
    }
}

func (repo *{{.StructName}}Repo) Get(
    ctx context.Context,
    {{varName .StructName}}Id domain.{{.StructName}}ID,
) (*domain.{{.StructName}}, error) {
    result, err := repo.getQueries().Get{{.StructName}}(ctx, db.Get{{.StructName}}Params{
        ID:     int64({{varName .StructName}}Id),
        UserID: int64(repo.userId),
    })
    if err != nil {
        return nil, err
    }
    return convertRowTo{{.StructName}}(result)
}

func (repo *{{.StructName}}Repo) All(
    ctx context.Context,
) ([]*domain.{{.StructName}}, error) {
    results, err := repo.getQueries().GetAll{{.StructName}}s(ctx)
    if err != nil {
        return nil, err
    }
    {{varName .StructName}}s := make([]*domain.{{.StructName}}, len(results))
    for i, result := range results {
        {{.StructName}}, err := convertRowTo{{.StructName}}(result)
        if err != nil {
            return nil, err
        }
        {{varName .StructName}}s[i] = {{.StructName}}
    }
    return {{varName .StructName}}s, nil
}

func (repo *{{.StructName}}Repo) Create(
    ctx context.Context,
    {{varName .StructName}} *domain.{{.StructName}},
) error {
    row, err := convert{{.StructName}}ToRow(repo.userId, {{varName .StructName}})
    if err != nil {
        return err
    }
    _, err = repo.getQueries().Create{{.StructName}}(ctx, db.Create{{.StructName}}Params{
        ID:        row.ID,
        UserID:    row.UserID,
        CreatedAt: row.CreatedAt,
        Data:      row.Data,
    })
    return err
}

func (repo *{{.StructName}}Repo) Update(
    ctx context.Context,
    {{varName .StructName}} *domain.{{.StructName}},
) error {
    row, err := convert{{.StructName}}ToRow(repo.userId, {{varName .StructName}})
    if err != nil {
        return err
    }
    _, err = repo.getQueries().Update{{.StructName}}(ctx, db.Update{{.StructName}}Params{
        ID:        row.ID,
        UserID:    row.UserID,
        UpdatedAt: row.UpdatedAt,
        Data:      row.Data,
    })
    return err
}

func (repo *{{.StructName}}Repo) Delete(
    ctx context.Context,
    {{varName .StructName}}Id domain.{{.StructName}}ID,
) error {
    _, err := repo.getQueries().Delete{{.StructName}}(ctx, db.Delete{{.StructName}}Params{
        ID:     int64({{varName .StructName}}Id),
        UserID: int64(repo.userId),
    })
    return err
}

func convertRowTo{{.StructName}}(result db.{{sqlcName .StructName}}) (*domain.{{.StructName}}, error) {
    var data domain.{{.StructName}}Data
    err := json.Unmarshal([]byte(result.Data), &data)
    if err != nil {
        return nil, err
    }
    return &domain.{{.StructName}}{
        ID:           domain.{{.StructName}}ID(result.ID),
        {{.StructName}}Data: data,
        Metadata: domain.Metadata{
            CreatedAt: result.CreatedAt,
            UpdatedAt: result.UpdatedAt.Int64,
        },
    }, nil
}

func convert{{.StructName}}ToRow(
    userId domain.UserID,
    {{.StructName}} *domain.{{.StructName}},
) (db.{{sqlcName .StructName}}, error) {
    data, err := json.Marshal({{.StructName}}.{{.StructName}}Data)
    if err != nil {
        return db.{{sqlcName .StructName}}{}, err
    }
    return db.{{sqlcName .StructName}}{
        ID:        int64({{.StructName}}.ID),
        UserID:    int64(userId),
        CreatedAt: {{.StructName}}.Metadata.CreatedAt,
        UpdatedAt: sql.NullInt64{
            Int64: {{.StructName}}.Metadata.UpdatedAt,
            Valid: {{.StructName}}.Metadata.UpdatedAt != 0,
        },
        Data: data,
    }, nil
}
`

const scehmaTemplate = `
DROP TABLE IF EXISTS {{dbName .StructName}};
CREATE TABLE {{dbName .StructName}} (
    ID BIGINT PRIMARY KEY,
    USER_ID BIGINT NOT NULL,
    CREATED_AT BIGINT NOT NULL,
    UPDATED_AT BIGINT,
    DATA JSONB NOT NULL,
    FOREIGN KEY (USER_ID) REFERENCES users(ID) ON DELETE CASCADE
);
`

const queriesTemplate = `
-- name: Create{{.StructName}} :execresult
INSERT INTO {{dbName .StructName}} (id, user_id, created_at, data)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, created_at, data;

-- name: Update{{.StructName}} :execresult
UPDATE {{dbName .StructName}}
SET updated_at = $3, data = $4
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: Delete{{.StructName}} :execresult
DELETE FROM {{dbName .StructName}}
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, created_at, updated_at, data;

-- name: Get{{.StructName}} :one
SELECT id, user_id, created_at, updated_at, data
FROM {{dbName .StructName}}
WHERE id = $1 AND user_id = $2;

-- name: GetAll{{.StructName}}s :many
SELECT id, user_id, created_at, updated_at, data
FROM {{dbName .StructName}}
ORDER BY created_at DESC;
`

func variableNameFromSQLCName(structName string) string {
	caser := cases.Title(language.AmericanEnglish)
	return caser.String(structName)
}

func variableNameFromStructName(structName string) string {
	return strings.ToLower(string(structName[0])) + structName[1:]
}

func dbNameFromStructName(structName string) string {
	return strings.ToLower(structName)
}
