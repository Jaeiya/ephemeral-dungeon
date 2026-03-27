package dialog

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/jaeiya/monster/internal/dictionary"
	"github.com/pelletier/go-toml/v2"
)

const (
	_dialogSqlPath  = "./assets/sql/dialog_tables.sqlite"
	_dialogTomlPath = "./assets/dialog/universal.toml"
)

type QueryCode uint8

const (
	Q_UNKNOWN = QueryCode(iota)
	Q_WHO
	Q_WHAT
	Q_WHEN
	Q_WHERE
	Q_WHY
	Q_HOW
	Q_DO
	Q_DOES
	Q_DID
)

var dict = func() *dictionary.Dictionary {
	d, err := dictionary.NewDictionary(dictionary.NewFileStorage("./data/dictionary.db"))
	if err != nil {
		panic(err)
	}
	return d
}()

type UniversalDialog struct {
	Questions []struct {
		Question []string `toml:"question"`
		Answer   string   `toml:"answer"`
	} `toml:"questions"`
}

type Question struct {
	Query  map[string]struct{}
	Answer string
}

func BuildDialog(db *sql.DB) error {
	var d UniversalDialog
	data, err := os.ReadFile(_dialogTomlPath)
	if err != nil {
		return fmt.Errorf("failed to read TOML dialog file: %w", err)
	}

	err = toml.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	sqlSchema, err := os.ReadFile(_dialogSqlPath)
	if err != nil {
		return fmt.Errorf("failed to read dialog sql file: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start dialog universal dialog transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(strings.TrimSpace(string(sqlSchema)))
	if err != nil {
		return fmt.Errorf("failed to create universal dialog tables: %w", err)
	}

	ansStmt, err := tx.Prepare("INSERT INTO answers (content) VALUES (?)")
	if err != nil {
		return fmt.Errorf("failed to prepare sql statement: %w", err)
	}
	defer ansStmt.Close()

	queryStmt, err := tx.Prepare(
		"INSERT INTO universal_dialog_queries (id, answer_id) VALUES (?, ?)",
	)
	if err != nil {
		return fmt.Errorf("failed to prepare sql statement: %w", err)
	}
	defer queryStmt.Close()

	for _, entry := range d.Questions {
		res, err := ansStmt.Exec(entry.Answer)
		if err != nil {
			return fmt.Errorf("failed to insert answer: %w", err)
		}
		ansID, _ := res.LastInsertId()
		for _, q := range entry.Question {
			queryBytes, err := EncodeQuestion(q)
			if err != nil {
				return fmt.Errorf("failed to encode question: %w", err)
			}
			if _, err = queryStmt.Exec(queryBytes, ansID); err != nil {
				return fmt.Errorf("failed to add query: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit dialog entries to database: %w", err)
	}

	return nil
}

func EncodeQuestion(s string) ([]byte, error) {
	s = strings.ToLower(s)

	firstSpace := strings.IndexByte(s, ' ')
	if firstSpace == -1 {
		return nil, fmt.Errorf("input too short")
	}

	qCode := encodeQueryWords(s[:firstSpace])
	if qCode == Q_UNKNOWN {
		return nil, fmt.Errorf("not a valid question")
	}

	contentCodes := make([]uint16, 0, 13)

	for word := range strings.FieldsSeq(s[firstSpace:]) {
		if val, exists := dict.WordMap[word]; exists {
			contentCodes = append(contentCodes, val)
		}
	}

	if len(contentCodes) == 0 {
		return nil, fmt.Errorf("no valid dictionary words referenced")
	}

	buf := make([]byte, 1, 1+len(contentCodes)*2)
	buf[0] = uint8(qCode)
	return binary.Append(buf, binary.BigEndian, contentCodes)
}

func encodeQueryWords(word string) QueryCode {
	switch word {
	case "who":
		return Q_WHO
	case "what":
		return Q_WHAT
	case "when":
		return Q_WHEN
	case "where":
		return Q_WHERE
	case "why":
		return Q_WHY
	case "how":
		return Q_HOW
	case "do":
		return Q_DO
	case "does":
		return Q_DOES
	case "did":
		return Q_DID

	default:
		return Q_UNKNOWN
	}
}
