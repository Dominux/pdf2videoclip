package common

import (
	"fmt"
	"strings"
)

var (
	ones     = []string{"", "один", "два", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять"}
	teens    = []string{"десять", "одиннадцать", "двенадцать", "тринадцать", "четырнадцать", "пятнадцать", "шестнадцать", "семнадцать", "восемнадцать", "девятнадцать"}
	tens     = []string{"", "", "двадцать", "тридцать", "сорок", "пятьдесят", "шестьдесят", "семьдесят", "восемьдесят", "девяносто"}
	hundreds = []string{"", "сто", "двести", "триста", "четыреста", "пятьсот", "шестьсот", "семьсот", "восемьсот", "девятьсот"}
)

// Триады: [0] форма для 1 (ед.ч), [1] форма для 2-4 (род.ч ед.ч), [2] форма для 5+ (род.ч мн.ч), [3] женский род?
type unit struct {
	forms [3]string
	isFem bool
}

var units = []unit{
	{[3]string{"", "", ""}, false},                 // Единицы
	{[3]string{"тысяча", "тысячи", "тысяч"}, true}, // Тысячи
	{[3]string{"миллион", "миллиона", "миллионов"}, false},
	{[3]string{"миллиард", "миллиарда", "миллиардов"}, false},
	{[3]string{"триллион", "триллиона", "триллионов"}, false},
}

func NumberToRussian(n int) string {
	if n == 0 {
		return "ноль"
	}
	if n < 0 {
		return "минус " + NumberToRussian(-n)
	}

	var parts []string
	unitIdx := 0

	for n > 0 {
		triad := n % 1000
		if triad > 0 {
			parts = append([]string{triadToText(triad, units[unitIdx])}, parts...)
		}
		n /= 1000
		unitIdx++
	}

	return strings.Join(parts, " ")
}

func triadToText(n int, u unit) string {
	var res []string

	h, t, o := n/100, (n%100)/10, n%10

	if h > 0 {
		res = append(res, hundreds[h])
	}

	lastTwo := n % 100
	if lastTwo >= 10 && lastTwo <= 19 {
		res = append(res, teens[lastTwo-10])
	} else {
		if t >= 2 {
			res = append(res, tens[t])
		}
		if o > 0 {
			if u.isFem {
				if o == 1 {
					res = append(res, "одна")
				} else if o == 2 {
					res = append(res, "две")
				} else {
					res = append(res, ones[o])
				}
			} else {
				res = append(res, ones[o])
			}
		}
	}

	// Добавляем само название разряда (тысячи, миллионы...)
	if u.forms[0] != "" {
		res = append(res, getForm(n, u.forms))
	}

	return strings.Join(res, " ")
}

func getForm(n int, forms [3]string) string {
	n = n % 100
	n1 := n % 10
	if n > 10 && n < 20 {
		return forms[2]
	}
	if n1 > 1 && n1 < 5 {
		return forms[1]
	}
	if n1 == 1 {
		return forms[0]
	}
	return forms[2]
}

func main() {
	numbers := []int{1, 12, 125, 1000, 2001, 105532, 1000000}
	for _, num := range numbers {
		fmt.Printf("%d: %s\n", num, NumberToRussian(num))
	}
}
