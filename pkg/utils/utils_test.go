package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPassword(t *testing.T) {
	tests := []struct {
		password     string
		specialSigns string
		expected     bool
	}{
		{"Password123", "", true},
		{"Password", "", true},         // Нет цифр.
		{"123456", "", true},           // Нет букв.
		{"", "", false},                // Пустая строка.
		{"Pass word123", "", false},    // Пробелы не допустимы.
		{"Password123!", "!", true},    // Допустимый специальный символ.
		{"Password123?", "?", true},    // Недопустимый специальный символ.
		{"Password!@#", "!@#", true},   // Несколько специальных символов.
		{"Password!@#%", "!@#", false}, // Несколько специальных символов.
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			got := IsPassword(test.password, test.specialSigns)
			if got != test.expected {
				t.Errorf("got %v, expected %v", got, test.expected)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}

func TestIsOneWord(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Test", true},
		{"Test Test", false},
		{"Test123", true},
		{"", true}, // Пустая строка.
		{"!@#$", true},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := IsOneWord(test.input)
			if got != test.expected {
				t.Errorf("got %v, expected %v", got, test.expected)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}

func TestIsLetter(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"a", true},
		{"A", true},
		{"1", false},
		{"!", false},
		{"a1", false},
		{"a!", false},
		{"A1", false},
		{"Test", true},
		{"Тест", true},
		{"1234", false},
		{"Test123", false},
		{"", false},
		{" ", false},
		{"Test Test", false},
		{"!@#$", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := IsLetter(test.input)
			if got != test.expected {
				t.Errorf("got %v, expected %v", got, test.expected)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}

func TestIsOneWordLetter(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Test", true},
		{"тест", true},    // Поддержка не-ASCII символов.
		{"Ярослав", true}, // Поддержка не-ASCII символов.
		{"1234", false},   // Цифры не допускаются.
		{"Test123", false},
		{"", true}, // Пустая строка.
		{"!@#$", false},
		{"Test Test", false}, // Пробелы не допускаются.
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := IsOneWordLetter(test.input)
			if got != test.expected {
				t.Errorf("got %v, expected %v", got, test.expected)
				assert.Equal(t, test.expected, got)
			}
		})
	}
}
