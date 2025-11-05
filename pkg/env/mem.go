package env

import "github.com/bosley/slpx/pkg/object"

type memImpl struct {
	parent  *memImpl
	symbols map[object.Identifier]object.Obj
}

var _ MEM = &memImpl{}

func DefaultMEM() MEM {
	return &memImpl{
		symbols: make(map[object.Identifier]object.Obj),
	}
}

func (m *memImpl) Get(key object.Identifier, searchParent bool) (object.Obj, error) {
	value, exists := m.symbols[key]
	if exists {
		return value, nil
	}
	if searchParent && m.parent != nil {
		return m.parent.Get(key, searchParent)
	}
	return object.Obj{}, ErrUndefinedIdentifier
}

func (m *memImpl) Set(key object.Identifier, value object.Obj, searchParent bool) error {
	_, exists := m.symbols[key]
	if exists {
		m.symbols[key] = value
		return nil
	}

	if searchParent && m.parent != nil {
		_, err := m.parent.Get(key, true)
		if err == nil {
			return m.parent.Set(key, value, searchParent)
		}
	}

	m.symbols[key] = value
	return nil
}

func (m *memImpl) Delete(key object.Identifier, searchParent bool) error {
	_, exists := m.symbols[key]
	if exists {
		delete(m.symbols, key)
		return nil
	}
	if searchParent && m.parent != nil {
		return m.parent.Delete(key, searchParent)
	}
	return nil
}

func (m *memImpl) Keys() []object.Identifier {
	keys := make([]object.Identifier, 0, len(m.symbols))
	for key := range m.symbols {
		keys = append(keys, key)
	}
	return keys
}

func (m *memImpl) Values() []object.Obj {
	values := make([]object.Obj, 0, len(m.symbols))
	for _, value := range m.symbols {
		values = append(values, value)
	}
	return values
}

func (m *memImpl) Len() int {
	return len(m.symbols)
}

func (m *memImpl) IsEmpty() bool {
	return len(m.symbols) == 0
}

func (m *memImpl) Clear() {
	m.symbols = make(map[object.Identifier]object.Obj)
}

func (m *memImpl) GetAll() map[object.Identifier]object.Obj {
	return m.symbols
}

func (m *memImpl) Fork() MEM {
	return &memImpl{parent: m, symbols: make(map[object.Identifier]object.Obj)}
}
