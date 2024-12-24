package main

import "testing"

func TestExamples(t *testing.T) {
	num := 123
	results := []int{
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}
	for i, secret := range results {
		num = NextSecret(num)
		if num != secret {
			t.Errorf("secret number %d not equal. num: %d, expected: %d", i, num, secret)
		}
	}
}

func TestAfter(t *testing.T) {
	results := []int{
		123,
		15887950,
		16495136,
		527345,
		704524,
		1553684,
		12683156,
		11100544,
		12249484,
		7753432,
		5908254,
	}

	for i, secret := range results {
		calculated := AfterNSecret(123, i)
		if calculated != secret {
			t.Errorf("secret number calculated through After%dSecret is wrong: num: %d, expected %d", i, calculated, secret)
		}
	}

	examples := map[int]int{1: 8685429, 10: 4700978, 100: 15273692, 2024: 8667524}

	for init, expected := range examples {
		calculated := AfterNSecret(init, 2000)

		if calculated != expected {
			t.Errorf("secret number calculated after 2000 is wrong: calculated %d, expected %d", calculated, expected)
		}
	}
}

func TestSequences(t *testing.T) {
	sequences := BuildSequences(123, 10)

	if sequences[Sequence{-3, 6, -1, -1}] != 4 {
		t.Errorf("first sequence of -3, 6, -1, -1: 4 is not found")
	}

	if sequences[Sequence{2, -2, 0, -2}] != 2 {
		t.Errorf("last sequence of 2, -2, 0, -2: 2 is not found")
	}
}
