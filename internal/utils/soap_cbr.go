package utils

// ValidateLuhn проверяет номер карты по алгоритму Луна
func ValidateLuhn(number string) bool {
	var sum int
	var alternate bool

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// GenerateCardNumber генерирует номер карты, соответствующий алгоритму Луна
func GenerateCardNumber(prefix string, length int) string {
	// Проверка входных данных
	if len(prefix) >= length {
		return prefix
	}

	// Заполнение номера карты до нужной длины (кроме последней цифры)
	number := prefix
	for i := len(prefix); i < length-1; i++ {
		number += "0"
	}

	// Вычисление последней цифры для соответствия алгоритму Луна
	var sum int
	var alternate bool

	for i := length - 2; i >= 0; i-- {
		digit := int(number[i] - '0')

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	lastDigit := (10 - (sum % 10)) % 10
	number += string('0' + lastDigit)

	return number
}
