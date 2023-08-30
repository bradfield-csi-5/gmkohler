package metrics

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type UserAge int
type CentsAmount uint64

const (
	dollarsToCents float64 = 100
	unrollFactor           = 4
)

func AverageAge(ages []UserAge) float64 {
	acc0, acc1, acc2, acc3 := UserAge(0), UserAge(0), UserAge(0), UserAge(0)
	count := len(ages)
	loopEnd := count - unrollFactor
	var j int
	for j = 0; j < loopEnd; j += unrollFactor {
		acc0 += ages[j]
		acc1 += ages[j+1]
		acc2 += ages[j+2]
		acc3 += ages[j+3]
	}

	for ; j < count; j++ {
		acc0 += ages[j]
	}
	return float64(acc0+acc1+acc2+acc3) / float64(count)
}

func AveragePaymentAmount(payments []CentsAmount) float64 {
	return averagePaymentAmountCents(payments) / dollarsToCents
}

// StdDevPaymentAmount computes the standard deviation of payment amounts
func StdDevPaymentAmount(payments []CentsAmount) float64 {
	sqDiff0, sqDiff1, sqDiff2, sqDiff3 := 0.0, 0.0, 0.0, 0.0
	mean := averagePaymentAmountCents(payments)
	count := len(payments)
	loopEnd := count - unrollFactor
	var j int
	for j = 0; j < loopEnd; j += unrollFactor {
		diff0 := float64(payments[j]) - mean
		diff1 := float64(payments[j+1]) - mean
		diff2 := float64(payments[j+2]) - mean
		diff3 := float64(payments[j+3]) - mean

		sqDiff0 += diff0 * diff0
		sqDiff1 += diff1 * diff1
		sqDiff2 += diff2 * diff2
		sqDiff3 += diff3 * diff3
	}
	for ; j < count; j++ {
		diff := float64(payments[j]) - mean
		sqDiff0 += diff * diff
	}
	return math.Sqrt((sqDiff0+sqDiff1+sqDiff2+sqDiff3)/float64(count)) /
		dollarsToCents
}

// averagePaymentAmountCents calculates the average in cents so
// StdDevPaymentAmount can use it without converting each of its entries to
// dollars via division (instead dividing once at the end of that function)
func averagePaymentAmountCents(payments []CentsAmount) float64 {
	acc0, acc1, acc2, acc3 :=
		CentsAmount(0), CentsAmount(0), CentsAmount(0), CentsAmount(0)
	count := len(payments)
	loopEnd := count - unrollFactor
	var j int
	for j = 0; j < loopEnd; j += unrollFactor {
		acc0 += payments[j]
		acc1 += payments[j+1]
		acc2 += payments[j+2]
		acc3 += payments[j+3]
	}
	for ; j < count; j++ {
		acc0 += payments[j]
	}

	return float64(acc0+acc1+acc2+acc3) / float64(count)
}

func LoadData() ([]UserAge, []CentsAmount) {
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
	var userAges []UserAge
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Unable to parse users.csv as csv", err)
		}
		age, _ := strconv.Atoi(line[2])
		userAges = append(userAges, UserAge(age))
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
	var payments []CentsAmount
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Unable to parse payments.csv as csv", err)
		}
		paymentCents, _ := strconv.Atoi(line[0])
		payments = append(
			payments,
			CentsAmount(paymentCents),
		)
	}

	return userAges, payments
}
