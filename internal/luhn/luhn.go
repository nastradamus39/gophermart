package luhn

import (
	"strconv"
)

// Validate Checks whether the given number satisfies the Luhn condition
func Validate(num string) bool {
	number := reverse(num)
	n := len(number)
	sum := 0

	for i := 0; i < n; i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		// If even number multiply by two
		if (i+1)%2 == 0 {
			digit = digit * 2
		}

		// Sum of digits if it's greater than 9
		// 18 => 9  : 18 - 9 => 9
		// Subtracting by 9 gets the sum of the digits
		if digit > 9 {
			digit -= 9
		}
		sum += digit
	}
	return sum%10 == 0
}

// GenerateNumber and return only the checksum
func GenerateNumber(num string) (string, error) {
	n := len(num)
	number := reverse(num)
	sum := 0
	for i := 0; i < n; i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return "", err
		}
		if i%2 == 0 {
			digit *= 2
		}
		if digit > 9 {
			digit -= 9
		}
		sum += digit
	}
	if sum%10 == 0 {
		return "0", nil
	}
	return strconv.Itoa(10 - (sum % 10)), nil
}

// Generate the checksum and return the whole number
func Generate(num string) (string, error) {
	number, err := GenerateNumber(num)
	if err != nil {
		return "", err
	}
	return num + number, nil
}

// Reverse a string
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
