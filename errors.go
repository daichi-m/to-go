package togo

// import (
// 	"fmt"
// 	"reflect"
// )

// // FieldError indicates that the Field operation has failed due to some issue.
// type FieldError struct {
// 	field   Field
// 	message string
// }

// func (fe FieldError) Error() string {
// 	return fmt.Sprintf("Field %s errored due to %s", fe.field.name, fe.message)
// }

// // UnsupportedType is an error type whenever an unsupported data type is encountered
// // during creation of a Field or GoStruct
// type UnsupportedType struct {
// 	data interface{}
// }

// func (ut UnsupportedType) Error() string {
// 	tp := reflect.TypeOf(ut.data).Kind()
// 	return fmt.Sprintf("Data type %s is not recognized", tp)
// }

// // GoStructError indicates that operation on the GoStruct failed due to some issue.
// type GoStructError struct {
// 	gs      GoStruct
// 	message string
// }

// func (gse GoStructError) Error() string {
// 	return fmt.Sprintf("GoStruct %s errored due to %s", gse.gs.name, gse.message)
// }
