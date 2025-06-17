package validate

import (
	"context"

	"github.com/Nemutagk/godb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func Exists(input string, payload map[string]interface{}, options []string, errors map[string]interface{}, addError func(string, string, map[string]interface{}, string) map[string]interface{}, dbManager *godb.ConnectionManager) map[string]interface{} {
	if len(options) != 3 {
		errors = addError(input, "exists", errors, "La configuración de conexiones no es correcta")
		return errors
	}

	conn, _ := dbManager.GetRawConnection(options[0])

	if dbConn, ok := conn.(*gorm.DB); ok {
		var exists_row struct{}
		err := dbConn.Table(options[1]).Where(options[2]+" = ?", payload[input]).Limit(1).Find(&exists_row).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				errors = addError(input, "exists", errors, "El valor no existe")
			}
		}
	} else if dbConn, ok := conn.(*mongo.Database); ok {
		coll := dbConn.Collection(options[1])

		count_rows, err_count := coll.CountDocuments(context.TODO(), bson.M{options[2]: payload[input]})

		if err_count != nil {
			errors = addError(input, "exists", errors, "El valor no existe")
		}

		if count_rows == 0 {
			errors = addError(input, "exists", errors, "El valor no existe")
		}
	}

	return errors
}
