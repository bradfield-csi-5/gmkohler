package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type UserId int
type UserMap map[UserId]*User
type CentsAmount uint64

type User struct {
	id       UserId
	age      int
	payments []CentsAmount
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
		for _, p := range u.payments {
			totalCents += p
		}
	}
	return float64(totalCents) / 100 / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(users UserMap) float64 {
	mean := AveragePaymentAmount(users)
	squaredDiffs := 0.0
	count := 0
	for _, u := range users {
		count += len(u.payments)
		for _, p := range u.payments {
			// mean is a public API in dollars, so we have to convert here.
			diff := float64(p)/100 - mean
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
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	users := make(UserMap, len(userLines))
	for _, line := range userLines {
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
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	for _, line := range paymentLines {
		userId, _ := strconv.Atoi(line[2])
		paymentCents, _ := strconv.Atoi(line[0])
		users[UserId(userId)].payments = append(
			users[UserId(userId)].payments,
			CentsAmount(paymentCents),
		)
	}

	return users
}
