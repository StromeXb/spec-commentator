package test

import (
	"fmt"
	"os"
	"strconv"
)

// срабатывает после импорта пакета, поэтому не работает
func SetTesting() {
	os.Setenv(getTestEnvName(), "true")
}

// IsTest получает информацию о том является ли runtime запуском тестов.
// Эта информация хранится в переменной среды %s_TEST, где %s - значение переменной PREFIX,
// устанавливаемой vioper.
// Переменная должна быть установлена до запуска, потому что импорты
// срабатывают раньше, чем выполняются какие-либо инструкции, которые могут установить переменные среды,
// как это работает в случае с viper и как это победить, - я не понял.
func IsTest() bool {
	envName := getTestEnvName()
	value, ok := os.LookupEnv(envName)
	if ok {
		result, err := strconv.ParseBool(value)
		if err != nil {
			return false
		}

		return result
	}

	return false
}

func getTestEnvName() string {
	envPrefix := "PROJECT"
	// PREFIX - эта переменная устанавливается viper
	// метод SetEnvPrefix
	if value, ok := os.LookupEnv("PREFIX"); ok {
		envPrefix = value
	}
	return fmt.Sprintf("%s_TEST", envPrefix)
}
