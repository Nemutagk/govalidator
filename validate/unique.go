package validate

import (
	"context"
	"errors"
	"log"

	"github.com/Nemutagk/godb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

/**
* Valida que el valor indicado no exista en la base de datos
* Parámetros: conexión db, nombre de tabla, nombre de columna
* Ejemplo: unique:princial,users,email
 */
func Unique(input string, payload map[string]interface{}, options []string, list_errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, dbConn *godb.ConnectionManager) map[string]interface{} {
	if len(options) < 3 || len(options) > 4 {
		panic("the options for connection is invalid")
	}

	if _, exists_input := payload[input]; !exists_input {
		log.Println("validate unique:input not exists")
		return list_errors
	}

	email := payload[input]
	var row map[string]interface{}

	conn, _ := dbConn.GetRawConnection(options[0])

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
			log.Println("find one in mongodb 1")
			err = conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email}).Decode(&row)
		} else if len(options) == 4 {
			log.Println("find one in mongodb 2")
			id := options[3]
			err = conn.(*mongo.Database).Collection(table).FindOne(context.TODO(), bson.M{column: email, "_id": bson.M{"$ne": id}}).Decode(&row)
		}
	}
	log.Println("db err: ", err)
	log.Println("row: ", row)

	// para las validaciones
	// if ok == mysql
	// if !ok == mongodb
	if (ok && err == nil) || (!ok && err == nil && row != nil) {
		list_errors = addError(input, "unique", list_errors, "El valor ya ha sido registrado en la tabla")
	} else if err != nil && ((ok && !errors.Is(err, gorm.ErrRecordNotFound)) || (!ok && err != mongo.ErrNoDocuments)) {
		list_errors = addError(input, "unique", list_errors, "Error al conectar con la base de datos")
	}

	return list_errors
}
