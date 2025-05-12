package validate

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nemutagk/govalidator/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

/**
* Validates that the indicated value does not exists in the database
* Params: db connection, table name, column name
* Example: unique:princial,users,email
 */
func Unique(input string, payload map[string]interface{}, options []string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, dbConn *db.ConnectionManager) map[string]interface{} {
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

	var err error
	_, ok := conn.(*gorm.DB)

	if ok {
		table := options[1]
		column := options[2]
		if len(options) == 3 {
			err = conn.(*gorm.DB).Table(table).Select(column).Where(column+" = ?", email).Limit(1).Take(&row).Error
		} else if len(options) == 4 {
			id := options[3]
			err = conn.(*gorm.DB).Table(table).Select(column).Where(column+" = ? AND id != ?", email, id).Limit(1).Take(&row).Error
		}
	} else {
		table := options[1]
		column := options[2]
		if len(options) == 3 {
			err = conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email}).Decode(&row)
		} else if len(options) == 4 {
			id := options[3]
			err = conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email, "_id": bson.M{"$ne": id}}).Decode(&row)
		}
	}
	fmt.Println("db err: ", err)
	fmt.Println("row: ", row)

	// for the validations
	// if ok == mysql
	// if !ok == mongodb
	if (ok && err == nil) || (!ok && err == nil && row != nil) {
		list_errors = addError(input, "unique", list_errors, "the value has been taken in the table")
	} else if err != nil && ((ok && !errors.Is(err, gorm.ErrRecordNotFound)) || (!ok && err != mongo.ErrNoDocuments)) {
		list_errors = addError(input, "unique", list_errors, "error to connect with db")
	}

	return list_errors
}
