package error_log

import "log"

func ReturnValueIfErrorIsNotNilLogFatalIfError[Value any, Error error](value Value, error Error) Value {

	if any(error) != any(nil) {

		log.Fatal(error)
	}

	return value

}
