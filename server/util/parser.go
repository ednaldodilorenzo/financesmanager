package util

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	s "strings"
	"time"
)

type FileImportType int

const (
	BBCA FileImportType = iota
	C6CC
	CUAL
)

type TransactionImportedData struct {
	Description     string
	AccountName     *string
	CategoryName    *string
	Value           int32
	PaymentDate     time.Time
	TransactionDate time.Time
}

// ParserFunction defines the signature of a parser function.
type ParserFunction func(io.Reader, time.Time) ([]TransactionImportedData, error)

// ParserFactory holds a registry of parser functions.
type ParserFactory struct {
	parsers map[FileImportType]ParserFunction
}

// NewParserFactory initializes a new parser factory.
func NewParserFactory() *ParserFactory {
	parserFactory := ParserFactory{
		parsers: make(map[FileImportType]ParserFunction),
	}
	parserFactory.RegisterParser(BBCA, parseBBCurrentAccount)
	parserFactory.RegisterParser(C6CC, parseCreditCardData)
	parserFactory.RegisterParser(CUAL, parseCustomData)
	return &parserFactory
}

// RegisterParser registers a parser function with a name.
func (pf *ParserFactory) RegisterParser(fileType FileImportType, parser ParserFunction) {
	pf.parsers[fileType] = parser
}

// GetParser retrieves a parser function by name.
func (pf *ParserFactory) GetParser(fileType FileImportType) (ParserFunction, error) {
	parser, exists := pf.parsers[fileType]
	if !exists {
		return nil, errors.New("parser not found")
	}
	return parser, nil
}

func parseBBCurrentAccount(fileReader io.Reader, _ time.Time) ([]TransactionImportedData, error) {
	csvReader := csv.NewReader(fileReader)
	csvReader.Comma = ','       // Adjust if the delimiter is different
	csvReader.LazyQuotes = true // Handle embedded quotes

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]TransactionImportedData, 0, 100)
	for i, record := range records {
		if (i == 0) || (i == 1) || (i == len(records)-1) {
			// Skip header row saldo and total
			continue
		}

		// Parse date
		transactionDate, err := time.Parse("02/01/2006", record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid date format at line %d: %v", i+1, err)
		}

		// Parse value
		value, err := strconv.ParseFloat(record[5], 32)
		if err != nil {
			return nil, fmt.Errorf("invalid value format at line %d: %v", i+1, err)
		}
		intValue := int32(math.Round(value * 100)) // Convert to cents

		// Create transaction schema
		transaction := TransactionImportedData{
			Description:     record[2],
			Value:           intValue,
			PaymentDate:     transactionDate,
			TransactionDate: transactionDate,
		}

		result = append(result, transaction)
	}

	return result, nil
}

func parseCreditCardData(fileReader io.Reader, date time.Time) ([]TransactionImportedData, error) {

	// Parse the CSV file
	csvReader := csv.NewReader(fileReader)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []TransactionImportedData
	// Iterate through the records and save them to the database
	for i, record := range records {

		if i == 0 {
			continue
		}

		value, err := strconv.ParseFloat(record[8], 32)
		if err != nil {
			return nil, err
		}

		intValue := int32(math.Round(-value * 100))

		transactionDate, err := time.Parse("02/01/2006", record[0]) // Adjust format as needed
		if err != nil {
			return nil, err
		}

		stallments := s.Split(record[5], "/")
		currentStallment := 0
		lastStallment := 0
		if len(stallments) > 1 {
			currentStallment, _ = strconv.Atoi(stallments[0])
			lastStallment, _ = strconv.Atoi(stallments[1])
		}

		monthCounter := 0
		for i := currentStallment; i <= lastStallment; i++ {
			currentDate := date.AddDate(0, monthCounter, 0)

			description := record[4]
			if currentStallment > 0 {
				description = fmt.Sprintf("%s (%d/%d)", description, currentStallment+monthCounter, lastStallment)
			}

			dbRecord := TransactionImportedData{
				Description:     description,
				Value:           intValue,
				PaymentDate:     currentDate,
				TransactionDate: transactionDate,
			}
			monthCounter += 1
			result = append(result, dbRecord)
		}
	}

	return result, nil
}

func parseCustomData(fileReader io.Reader, _ time.Time) ([]TransactionImportedData, error) {
	//Parse the CSV file
	csvReader := csv.NewReader(fileReader)
	csvReader.Comma = '\t'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []TransactionImportedData
	// Iterate through the records and save them to the database
	for _, record := range records {

		value, err := strconv.ParseFloat(record[2], 32)
		if err != nil {
			return nil, err
		}

		intValue := int32(math.Round(value * 100))

		paymentDate, err := time.Parse("02/01/2006", record[0]) // Adjust format as needed
		if err != nil {
			return nil, err
		}

		transactionDate, err := time.Parse("02/01/2006", record[5]) // Adjust format as needed
		if err != nil {
			return nil, err
		}

		categoryName := record[3]
		accountName := record[4]

		// Create a new Record instance
		dbRecord := TransactionImportedData{
			Description:     record[1],
			Value:           intValue,
			PaymentDate:     paymentDate,
			TransactionDate: transactionDate,
			CategoryName:    &categoryName,
			AccountName:     &accountName,
		}

		transactions = append(transactions, dbRecord)

	}

	return transactions, nil
}
