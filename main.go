package main

import (
	"fmt"
	"math"
	"strconv"
)

// Helper function to calculate Future Value of an annuity
func FV(rate float64, nper int, pmt float64, pv float64) float64 {
	return pv*math.Pow(1+rate, float64(nper)) + pmt*((math.Pow(1+rate, float64(nper))-1)/rate)
}

// Helper function to calculate Present Value of an annuity
func PV(rate float64, nper int, pmt float64, fv float64) float64 {
	return fv/math.Pow(1+rate, float64(nper)) + (pmt * (1 - math.Pow(1+rate, -float64(nper))) / rate)
}

// Helper function to calculate the PMT (payment) needed for a loan or investment
func PMT(rate float64, nper int, pv float64, fv float64) float64 {
	return (rate * (fv - pv*math.Pow(1+rate, float64(nper)))) / (math.Pow(1+rate, float64(nper)) - 1)
}

// Helper function to format numbers as currency with commas
func formatCurrency(value float64) string {
	roundedValue := int64(math.Round(value))
	return fmt.Sprintf("$%s", formatWithCommas(roundedValue))
}

// Helper function to add commas to a number
func formatWithCommas(value int64) string {
	// Convert the number to a string
	str := strconv.FormatInt(value, 10)

	// Handle negative numbers
	if value < 0 {
		str = str[1:]
	}

	// Find the length of the string
	n := len(str)

	// Add commas every three digits
	if n <= 3 {
		return str
	}

	// Compute where the first comma should be
	mod := n % 3
	commaStr := str[:mod]

	for i := mod; i < n; i += 3 {
		if i != 0 {
			commaStr += ","
		}
		commaStr += str[i : i+3]
	}

	// Add the negative sign back if needed
	if value < 0 {
		return "-" + commaStr
	}
	return commaStr
}

func main() {
	var startingAge, retirementAge, expiryAge int
	var desiredIncome, currentSavings, nominalReturnAccumulation, nominalReturnRetirement, inflationRate float64

	// Prompt user to enter values
	fmt.Print("Enter Current Age (Default: 32): ")
	fmt.Scanln(&startingAge)
	if startingAge == 0 {
		startingAge = 30
	}

	fmt.Print("Enter Retirement Age (Default: 60): ")
	fmt.Scanln(&retirementAge)
	if retirementAge == 0 {
		retirementAge = 60
	}

	fmt.Print("Enter Expiry Age (Default: 90): ")
	fmt.Scanln(&expiryAge)
	if expiryAge == 0 {
		expiryAge = 90
	}

	fmt.Print("Enter Desired Income in Today's Dollars (Default: 100000): ")
	fmt.Scanln(&desiredIncome)
	if desiredIncome == 0 {
		desiredIncome = 100000
	}

	fmt.Print("Enter Current Savings (Default: 100000): ")
	fmt.Scanln(&currentSavings)
	if currentSavings == 0 {
		currentSavings = 100000
	}

	fmt.Print("Enter Nominal Return in Accumulation (Default: 0.08): ")
	fmt.Scanln(&nominalReturnAccumulation)
	if nominalReturnAccumulation == 0 {
		nominalReturnAccumulation = 0.08
	}

	fmt.Print("Enter Nominal Return in Retirement (Default: 0.06): ")
	fmt.Scanln(&nominalReturnRetirement)
	if nominalReturnRetirement == 0 {
		nominalReturnRetirement = 0.06
	}

	fmt.Print("Enter Inflation Rate (Default: 0.025): ")
	fmt.Scanln(&inflationRate)
	if inflationRate == 0 {
		inflationRate = 0.025
	}

	// Calculate Desired Income at Retirement in future dollars
	yearsToRetirementAtStart := retirementAge - startingAge
	desiredIncomeRetirement := desiredIncome * math.Pow(1+inflationRate, float64(yearsToRetirementAtStart))

	fmt.Printf("\nDesired Income at Retirement: %s\n", formatCurrency(desiredIncomeRetirement))

	// Calculate Assets Needed at Retirement
	yearsInRetirement := expiryAge - retirementAge
	realReturnRetirement := nominalReturnRetirement - inflationRate
	assetsNeededAtRetirement := -PV(realReturnRetirement, yearsInRetirement, -desiredIncomeRetirement, 0)
	fmt.Printf("Assets Needed at Retirement: %s\n", formatCurrency(assetsNeededAtRetirement))

	// Calculate savings needed each year to hit Coast FIRE targets
	for age := startingAge + 1; age <= retirementAge-1; age++ {
		yearsAccumulating := age - startingAge
		yearsToRetirement := retirementAge - age
		assetsNeededAtAge := PV(nominalReturnAccumulation, yearsToRetirement, 0, assetsNeededAtRetirement)

		monthlyNominalReturnAccumulation := math.Pow(1+nominalReturnAccumulation, float64(1.0/12.0)) - 1
		savingsNeeded := PMT(monthlyNominalReturnAccumulation, yearsAccumulating*12, currentSavings, assetsNeededAtAge)
		totalContributions := savingsNeeded * float64(yearsAccumulating*12)

		fmt.Printf("\nAt Age %d:\n", age)
		fmt.Printf("  Assets Needed: %s\n", formatCurrency(assetsNeededAtAge))
		fmt.Printf("  Monthly Savings Needed: %s\n", formatCurrency(savingsNeeded))
		fmt.Printf("  Total Contributions: %s\n", formatCurrency(totalContributions))
	}
}
