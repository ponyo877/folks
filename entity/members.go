package entity

type Members struct {
	value map[UID]User
}

func (m *Members) HasUser(ID UID) bool {
	_, ok := m.value[ID]
	return ok
}
