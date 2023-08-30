package metrics

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type UserId int
type UserMap map[UserId]*User
type CentsAmount uint64

type User struct {
	age      int
	payments []CentsAmount
	id       UserId
}

func AverageAge(users UserMap) float64 {
	total := 0
	for _, u := range users {
		total += u.age
	}
	return float64(total) / float64(len(users))
}

func AveragePaymentAmount(users UserMap) float64 {
	totalCents := CentsAmount(0)
	count := 0
	for _, u := range users {
		count += len(u.payments)
		for j := range u.payments {
			totalCents += u.payments[j]
		}
	}
	return float64(totalCents) / 100 / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs := 0.0
	count := 0
	// users is a map, don't use indexing
	for _, u := range users {
		count += len(u.payments)
		// payments is an array, use indexing
		for j := range u.payments {
			// mean is a public API in dollars, so we have to convert here.
			diff := float64(u.payments[j])/100 - mean
			squaredDiffs += diff * diff
		}
	}
	return math.Sqrt(squaredDiffs / float64(count))
}

func LoadData() UserMap {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("Unable to close payments.csv", err)
		}
	}(f)
	reader := csv.NewReader(f)
	users := UserMap{}
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Unable to parse users.csv as csv", err)
		}
		id, _ := strconv.Atoi(line[0])
		age, _ := strconv.Atoi(line[2])
		users[UserId(id)] = &User{
			id:       UserId(id),
			age:      age,
			payments: []CentsAmount{},
		}
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln("Unable to close payments.csv", err)
		}
	}(f)
	reader = csv.NewReader(f)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Unable to parse payments.csv as csv", err)
		}
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		users[UserId(userId)].payments = append(
			users[UserId(userId)].payments,
			CentsAmount(paymentCents),
		)
	}

	return users
}
