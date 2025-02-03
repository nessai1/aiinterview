package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomStringFromCharset(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			return "", fmt.Errorf("cannot read generated character by crypto/rand: %w", err)
		}
		b[i] = charset[randomByte[0]%byte(len(charset))]
	}

	return string(b), nil
}

// GenerateUUIDv7 генерирует UUIDv7 согласно RFC 9562
func GenerateUUIDv7() (string, error) {
	b := make([]byte, 16)

	// Берём текущий timestamp в миллисекундах
	timestamp := uint64(time.Now().UnixMilli())

	// Записываем timestamp (первые 48 бит)
	binary.BigEndian.PutUint64(b, timestamp<<16)

	// Заполняем оставшиеся байты случайными данными
	_, err := rand.Read(b[6:])
	if err != nil {
		return "", err
	}

	// Проставляем версию (версия 7 = 0111b = 0x7)
	b[6] = (b[6] & 0x0F) | 0x70

	// Проставляем флаги варианта (RFC 4122, вариант 1: 10xx xxxx)
	b[8] = (b[8] & 0x3F) | 0x80

	// Преобразуем в строку стандартного формата UUID
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(b[0:4]),  // timestamp high
		binary.BigEndian.Uint16(b[4:6]),  // timestamp low
		binary.BigEndian.Uint16(b[6:8]),  // version + random
		binary.BigEndian.Uint16(b[8:10]), // variant + random
		b[10:]), nil
}
