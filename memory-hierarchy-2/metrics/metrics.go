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
type centsAmount uint64

const dollarsToCents float64 = 100

type User struct {
	age      int
	payments []centsAmount
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
	return averagePaymentAmountCents(users) / dollarsToCents
}

// StdDevPaymentAmount computes the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := averagePaymentAmountCents(users)
	squaredDiffs := 0.0
	count := 0
	// users is a map, don't use indexing
	for _, u := range users {
		payments := u.payments
		count += len(payments)
		for _, p := range payments {
			// mean is a public API in dollars, so we have to convert here.
			diff := float64(p) - mean
			squaredDiffs += diff * diff
		}
	}
	return math.Sqrt(squaredDiffs/float64(count)) / dollarsToCents
}

// averagePaymentAmountCents calculates the average in cents so
// StdDevPaymentAmount can use it without converting each of its entries to
// dollars via division (instead dividing once at the end of that function)
func averagePaymentAmountCents(users UserMap) float64 {
	totalCents := centsAmount(0)
	count := 0
	for _, u := range users {
		payments := u.payments
		count += len(payments)
		for _, p := range payments {
			totalCents += p
		}
	}
	return float64(totalCents) / float64(count)
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
			payments: []centsAmount{},
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
			centsAmount(paymentCents),
		)
	}

	return users
}
