package chess

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Club represents a USCF-rated chess Club
type Club struct {
	abbreviation string
	affiliateId  string
	name         string
}

func NewClub(theAbbreviation string, theAffiliateId string, theName string) *Club {
	checkRequiredField("Abbreviation", theAbbreviation)
	checkRequiredField("AffiliateId", theAffiliateId)
	checkRequiredField("Name", theName)

	return &Club{
		abbreviation: strings.ToUpper(theAbbreviation),
		affiliateId:  theAffiliateId,
		name:         theName,
	}
}

func (c *Club) GetAbbreviation() string {
	return c.abbreviation
}

func (c *Club) GetAffiliateId() string {
	return c.affiliateId
}

func (c *Club) GetName() string {
	return c.name
}

func checkRequiredField(theFieldName string, theFieldValue string) {
	if strings.TrimSpace(theFieldValue) == "" {
		err := errors.New("Club " + theFieldName + " may not be empty")
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func (c *Club) String() string {
	return "affiliateId: " + c.affiliateId + "\n" +
		"abbreviation: " + c.abbreviation + "\n" +
		"name: " + c.name + "\n"
}

var WCC = NewClub("WCC", "A5008948", "Waukesha Chess Club")
var SWCC = NewClub("SWCC", "A6011047", "Southwest Chess Club")

func GetClub(abbrev string) (*Club, error) {
	if WCC.abbreviation == strings.ToUpper(abbrev) {
		return WCC, nil
	}
	if SWCC.abbreviation == strings.ToUpper(abbrev) {
		return SWCC, nil
	}

	err := errors.New("No Club found for abbreviation: " + abbrev)
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
	return nil, err
}
