package interfaces

type ModelInterface interface {
	GetID() string
}

type ModelMap map[string]interface{}

func (i ModelMap) GetID() string {
	if id, ok := i["id"].(string); ok {
		return id
	}
	return ""
}
