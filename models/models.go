package models

import "time"

// Customer schema of the customer table
type Customer struct {
    BankruptcyIndicatorFlag  bool       `json:"bankruptcy_indicator_flag"`
    CompanyName              string     `json:"company_name"`
    CreatedDate              time.Time  `json:"created_date"`
    DateOfBirth              time.Time  `json:"date_of_birth"`
    FirstName                string     `json:"first_name"`
    LastName                 string     `json:"last_name"`
    LegalEntityId            int64      `json:"legal_entity_id"`
    LegalEntityState         string     `json:"legal_entity_state"`
    LegalEntityType          string     `json:"legal_entity_type"`
}

/*
// User schema of the user table
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Location string `json:"location"`
    Age      int64  `json:"age"`
}
*/
