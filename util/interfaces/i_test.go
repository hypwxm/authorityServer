package interfaces

import "testing"

type a struct{}

func (i a) GetID() string {
	return ""
}
func (i a) GetValue(key string) interface{} {
	return ""
}
func (i a) Set(key string, value interface{}) {

}
func (i a) GetValueWithDefault(key string, df interface{}) interface{} {
	return ""
}
func (i a) ToCamelKey() ModelInterface {
	return new(a)
}

func TestI(t *testing.T) {
	var ab ModelInterface
	ab = new(a)
	ab = make(ModelMap)
	t.Fatal(ab)
}
