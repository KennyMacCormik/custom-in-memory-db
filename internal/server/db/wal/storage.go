package wal

import "fmt"

func (w *Wal) Get(key string) (string, error) {
	result, err := w.st.Get(key)
	if err != nil {
		return "", fmt.Errorf("get failed: %w", err)
	}

	return result, nil
}
func (w *Wal) Set(key, value string) error {
	err := w.write([]byte("SET " + key + " " + value + "\n"))
	if err != nil {
		return fmt.Errorf("wal failed: %w", err)
	}

	err = w.st.Set(key, value)
	if err != nil {
		return fmt.Errorf("set failed: %w", err)
	}

	return nil
}
func (w *Wal) Del(key string) error {
	err := w.write([]byte("DEL " + key + "\n"))
	if err != nil {
		return fmt.Errorf("wal failed: %w", err)
	}

	err = w.st.Del(key)
	if err != nil {
		return fmt.Errorf("set failed: %w", err)
	}

	return nil
}
