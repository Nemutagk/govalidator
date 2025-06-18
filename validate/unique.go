package validate

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Nemutagk/godb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

/**
* Validates that the indicated value does not exists in the database
* Params: db connection, table name, column name
* Example: unique:princial,users,email
 */
func Unique(input string, payload map[string]interface{}, options []string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, dbConn *godb.ConnectionManager) map[string]interface{} {
	if len(options) < 3 || len(options) > 4 {
		panic("the options for connection is invalid")
	}

	if _, exists_input := payload[input]; !exists_input {
		fmt.Println("validate unique:input not exists")
		return list_errors
	}

	email := payload[input]
	var row map[string]interface{}

	conn, _ := dbConn.GetConnection(options[0])
	raw_conn := conn.GetRawConnection()

	var err error
	_, ok := raw_conn.(*gorm.DB)

	if ok {
		table := options[1]
		column := options[2]
		if len(options) == 3 {
			err = raw_conn.(*gorm.DB).Table(table).Select(column).Where(column+" = ?", email).Limit(1).Take(&row).Error
		} else if len(options) == 4 {
			id := options[3]
			err = raw_conn.(*gorm.DB).Table(table).Select(column).Where(column+" = ? AND id != ?", email, id).Limit(1).Take(&row).Error
		}
	} else {
		table := options[1]
		column := options[2]
		if len(options) == 3 {
			err = raw_conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email}).Decode(&row)
		} else if len(options) == 4 {
			id := options[3]
			err = raw_conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email, "_id": bson.M{"$ne": id}}).Decode(&row)
		}

		if err != nil {
			log.Println("error in unique validation:", err)
			if !errors.Is(err, mongo.ErrNoDocuments) && !errors.Is(err, mongo.ErrNilDocument) {
				list_errors = addError(input, "unique", list_errors, "Error al validar el valor")
			}
		}
	}

	return list_errors
}
