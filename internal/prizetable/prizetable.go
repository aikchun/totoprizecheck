package prizetable

func GetPrize(betType string, numbersMatched int, hasAdditionalNumber bool) string {
	switch betType {
	case "Ordinary":
		return getOrdinaryPrize(numbersMatched, hasAdditionalNumber)
	case "System 7":
		return getSystemSevenPrize(numbersMatched, hasAdditionalNumber)
	case "System 8":
		return getSystemEightPrize(numbersMatched, hasAdditionalNumber)
	case "System 9":
		return getSystemNinePrize(numbersMatched, hasAdditionalNumber)
	case "System 10":
		return getSystemTenPrize(numbersMatched, hasAdditionalNumber)
	case "System 11":
		return getSystemElevenPrize(numbersMatched, hasAdditionalNumber)
	case "System 12":
		return getSystemTwelvePrize(numbersMatched, hasAdditionalNumber)
	}
	return ""
}

func getOrdinaryPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$10"
		case 4:
			return "$50"
		case 5:
			return "Group 3"
		case 6:
			return "Group 1"

		}
	}

	switch numbersMatched {
	case 3:
		return "$25"
	case 4:
		return "Group 4"
	case 5:
		return "Group 2"
	}

	return "unknown"
}

func getSystemSevenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$40"
		case 4:
			return "$190"
		case 5:
			return "Group 3 + $250"
		case 6:
			return "Group 1 + 3"
		}
	}

	switch numbersMatched {
	case 3:
		return "$85"
	case 4:
		return "Group 4 + $150"
	case 5:
		return "Group 2 + 3 + 4"
	case 6:
		return "Group 1 + 2"
	}

	return "unknown"
}

func getSystemEightPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$100"
		case 4:
			return "$460"
		case 5:
			return "Group 3 + $850"
		case 6:
			return "Group 1 + 3 + $750"
		}
	}

	switch numbersMatched {
	case 3:
		return "$190"
	case 4:
		return "Group 4 + $490"
	case 5:
		return "Group 2 + 3 + 4 + $500"
	case 6:
		return "Group 1 + 2 + 3 + 4"
	}

	return "unknown"
}

func getSystemNinePrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$200"
		case 4:
			return "$900"
		case 5:
			return "Group 3 + $1,900"
		case 6:
			return "Group 1 + 3 + $2,450"
		}
	}

	switch numbersMatched {
	case 3:
		return "$350"
	case 4:
		return "Group 4 + $1,060"
	case 5:
		return "Group 2 + 3 + 4 + $1,600"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $1,250"
	}

	return "unknown"
}

func getSystemTenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$350"
		case 4:
			return "$1,550"
		case 5:
			return "Group 3 + $3,500"
		case 6:
			return "Group 1 + 3 + $5,300"
		}
	}

	switch numbersMatched {
	case 3:
		return "$575"
	case 4:
		return "Group 4 + $1,900"
	case 5:
		return "Group 2 + 3 + 4 + $3,400"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $3,950"
	}

	return "unknown"
}

func getSystemElevenPrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$560"
		case 4:
			return "$2,450"
		case 5:
			return "Group 3 + $5,750"
		case 6:
			return "Group 1 + 3 + $9,500"
		}
	}

	switch numbersMatched {
	case 3:
		return "$875"
	case 4:
		return "Group 4 + $3,050"
	case 5:
		return "Group 2 + 3 + 4 + $6,000"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $8,300"
	}

	return "unknown"
}

func getSystemTwelvePrize(numbersMatched int, hasAdditionalNumber bool) string {
	if !hasAdditionalNumber {
		switch numbersMatched {
		case 3:
			return "$840"
		case 4:
			return "$3,640"
		case 5:
			return "Group 3 + $8,750"
		case 6:
			return "Group 1 + 3 + $15,250"
		}
	}

	switch numbersMatched {
	case 3:
		return "$1,260"
	case 4:
		return "Group 4 + $4,550"
	case 5:
		return "Group 2 + 3 + 4 + $9,500"
	case 6:
		return "Group 1 + 2 + 3 + 4 + $14,500"
	}

	return "unknown"
}
