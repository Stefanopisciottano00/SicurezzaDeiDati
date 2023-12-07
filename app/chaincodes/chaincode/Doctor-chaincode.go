package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type DoctorContract struct {
	contractapi.Contract
}

type Doctor struct {
	DoctorID           string   `json:"doctorId"`
	FirstName          string   `json:"firstName"`
	LastName           string   `json:"lastName"`
	Password           string   `json:"password"`
	Age                int      `json:"age"`
	PhoneNumber        string   `json:"phoneNumber"`
	PatientsAssociated []string `json:"PatiensAssociated"`
}

func NewDoctor(doctorID, firstName, lastName, password string, age int, phoneNumber string) *Doctor {

	validateInputDoctor(doctorID, firstName, lastName, password, age, phoneNumber)

	// Hash della password usando SHA-256
	hashedPassword := hashPassword(password)

	return &Doctor{
		DoctorID:    doctorID,
		FirstName:   firstName,
		LastName:    lastName,
		Password:    hashedPassword,
		Age:         age,
		PhoneNumber: phoneNumber,
	}
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func validateInputDoctor(doctorID, firstName, lastName, password string, age int, phoneNumber string) error {

	// Controllo che i campi non siano vuoti
	if len(doctorID) == 0 || len(firstName) == 0 || len(lastName) == 0 || len(password) == 0 || len(phoneNumber) == 0 {
		return fmt.Errorf("I campi non possono essere vuoti...")
	}

	// Controllo che il campo patientID contenga solo numeri e lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", doctorID); !matched {
		return fmt.Errorf("Il campo patientID può contenere solo numeri o lettere")
	}

	// Controllo che il campo firstName contenga solo lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z]+$", firstName); !matched {
		return fmt.Errorf("Il campo firstName può contenere solo lettere")
	}

	// Controllo che il campo lastName contenga solo lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z]+$", lastName); !matched {
		return fmt.Errorf("Il campo firstName può contenere solo lettere")
	}

	// Controllo che il campo eta sia nel formato corretto
	if age <= 0 || age >= 150 {
		return fmt.Errorf("Inserire un'età corretta per il paziente!")
	}

	// Controllo che il campo PhoneNumber sia formato da 10 cifre
	if matched, _ := regexp.MatchString("^\\d{10}$", phoneNumber); !matched {
		return fmt.Errorf("Il campo phoneNumber deve contenere esattamente 10 cifre")
	}

	return nil
}

func (pc *DoctorContract) RegisterDoctor(ctx contractapi.TransactionContextInterface, doctorID string, firstname string, lastname string, password string, age int, phoneNumber string) error {

	// Verifica se il dottore esiste già
	existingDoctor, err := ctx.GetStub().GetState(doctorID)
	if err != nil {
		return fmt.Errorf("Errore durante la ricerca del dottore: %v", err)
	}

	if existingDoctor != nil {
		return fmt.Errorf("Il dottore con ID %s esiste già", doctorID)
	}

	// Crea un nuovo paziente
	newDoctor := NewDoctor(doctorID, firstname, lastname, password, age, phoneNumber)

	// Converti il paziente in formato JSON
	newDoctorJSON, err := json.Marshal(newDoctor)
	if err != nil {
		return fmt.Errorf("Errore durante la conversione del dottore in JSON: %v", err)
	}

	// Registra il paziente sulla blockchain
	err = ctx.GetStub().PutState(doctorID, newDoctorJSON)
	if err != nil {
		return fmt.Errorf("Errore durante la registrazione del dottore sulla blockchain: %v", err)
	}

	return nil

}

func (pc *DoctorContract) GetDoctor(ctx contractapi.TransactionContextInterface, patientID string) (*Doctor, error) {

	doctorJSON, err := ctx.GetStub().GetState(patientID)

	if err != nil {
		return nil, fmt.Errorf("Errore durante il recupero del dottore")
	}

	if doctorJSON == nil {
		return nil, fmt.Errorf("Errore, il dottore probabilmente non esiste!")
	}

	// Converti i dati JSON recuperati in un oggetto Patient
	var doctor Doctor

	err = json.Unmarshal(doctorJSON, &doctor)
	if err != nil {
		return nil, fmt.Errorf("Errore durante il parsing dei dati JSON del paziente: %s", err)
	}

	return &doctor, nil

}
