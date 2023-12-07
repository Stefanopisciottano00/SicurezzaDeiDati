package chaincode

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PatientContract struct {
	contractapi.Contract
}

type Patient struct {
	PatientID        string  `json:"patientId"`
	FirstName        string  `json:"firstName"`
	LastName         string  `json:"lastName"`
	Password         string  `json:"password"`
	Age              int     `json:"age"`
	PhoneNumber      string  `json:"phoneNumber"`
	EmergPhoneNumber string  `json:"emergPhoneNumber"`
	Address          string  `json:"address"`
	BloodGroup       string  `json:"bloodGroup"`
	Allergies        string  `json:"allergies"`
	Visits           []Visit `json:"visits"`
}

type Visit struct {
	DoctorID  string `json:"doctorId"`
	VisitDate string `json:"visitdate"`
	Symptoms  string `json:"symptoms"`
	Diagnosis string `json:"diagnosis"`
}

func NewPatient(patientID, firstName, lastName, password string, age int, phoneNumber, emergPhoneNumber, address, bloodGroup string) *Patient {

	validateInputPatient(patientID, firstName, lastName, password, age, phoneNumber, emergPhoneNumber, address, bloodGroup)

	// Hash della password usando SHA-256
	hashedPassword := hashPassword2(password)

	return &Patient{
		PatientID:        patientID,
		FirstName:        firstName,
		LastName:         lastName,
		Password:         hashedPassword,
		Age:              age,
		PhoneNumber:      phoneNumber,
		EmergPhoneNumber: emergPhoneNumber,
		Address:          address,
		BloodGroup:       bloodGroup,
	}
}

//-------------------------------- Funzioni ausiliarie -----------------------

func hashPassword2(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func validateInputPatient(patientID, firstName, lastName, password string, age int, phoneNumber, emergPhoneNumber, address, bloodGroup string) error {

	// Controllo che i campi non siano vuoti
	if len(patientID) == 0 || len(firstName) == 0 || len(lastName) == 0 || len(password) == 0 || len(address) == 0 || len(bloodGroup) == 0 {
		return fmt.Errorf("I campi non possono essere vuoti...")
	}

	// Controllo che il campo patientID contenga solo numeri e lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", patientID); !matched {
		return fmt.Errorf("Il campo patientID può contenere solo numeri o lettere")
	}

	// Controllo che il campo firstName contenga solo lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z]+$", firstName); !matched {
		return fmt.Errorf("Il campo firstName può contenere solo lettere")
	}

	// Controllo che il campo lastName contenga solo lettere
	if matched, _ := regexp.MatchString("^[a-zA-Z]+$", lastName); !matched {
		return fmt.Errorf("Il campo lastName può contenere solo lettere")
	}

	// Controllo che il campo eta sia nel formato corretto
	if age <= 0 || age >= 150 {
		return fmt.Errorf("Inserire un'età corretta per il paziente!")
	}

	// Controllo che il campo PhoneNumber sia formato da 10 cifre
	if matched, _ := regexp.MatchString("^\\d{10}$", phoneNumber); !matched {
		return fmt.Errorf("Il campo phoneNumber deve contenere esattamente 10 cifre")
	}

	// Controllo che il campo emergPhoneNumber sia formato da 10 cifre
	if matched, _ := regexp.MatchString("^\\d{10}$", emergPhoneNumber); !matched {
		return fmt.Errorf("Il campo emergPhoneNumber deve contenere esattamente 10 cifre")
	}

	if matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", bloodGroup); !matched {
		return fmt.Errorf("Il campo bloodGroup può contenere solo numeri o lettere")
	}

	return nil
}

//-------------------------------- Funzioni ausiliarie -----------------------

func (pc *PatientContract) RegisterPatient(ctx contractapi.TransactionContextInterface, patientID string, firstName string, lastName string, password string, age int, phoneNumber string, emergPhoneNumber string, address string, bloodGroup string) error {

	// Verifica se il paziente esiste già
	existingPatient, err := ctx.GetStub().GetState(patientID)
	if err != nil {
		return fmt.Errorf("Errore durante la ricerca del paziente: %v", err)
	}

	if existingPatient != nil {
		return fmt.Errorf("Il paziente con ID %s esiste già", patientID)
	}

	// Crea un nuovo paziente
	newPatient := NewPatient(patientID, firstName, lastName, password, age, phoneNumber, emergPhoneNumber, address, bloodGroup)

	// Converti il paziente in formato JSON
	newPatientJSON, err := json.Marshal(newPatient)
	if err != nil {
		return fmt.Errorf("Errore durante la conversione del paziente in JSON: %v", err)
	}

	// Registra il paziente sulla blockchain
	err = ctx.GetStub().PutState(patientID, newPatientJSON)
	if err != nil {
		return fmt.Errorf("Errore durante la registrazione del paziente sulla blockchain: %v", err)
	}

	return nil
}

func (pc *PatientContract) GetPatient(ctx contractapi.TransactionContextInterface, patientID string) (*Patient, error) {

	patientJSON, err := ctx.GetStub().GetState(patientID)

	if err != nil {
		return nil, fmt.Errorf("Errore durante il recupero del paziente")
	}

	if patientJSON == nil {
		return nil, fmt.Errorf("Errore, il paziente probabilmente non esiste!")
	}

	// Converti i dati JSON recuperati in un oggetto Patient
	var patient Patient

	err = json.Unmarshal(patientJSON, &patient)
	if err != nil {
		return nil, fmt.Errorf("Errore durante il parsing dei dati JSON del paziente: %s", err)
	}

	return &patient, nil

}

func (pc *PatientContract) RecordVisit(ctx contractapi.TransactionContextInterface, patientID, doctorID, visitDate, symptoms, diagnosis string) error {

	// Controllo che i campi non siano vuoti
	if len(patientID) == 0 || len(doctorID) == 0 || len(visitDate) == 0 || len(symptoms) == 0 || len(diagnosis) == 0 {
		return fmt.Errorf("I campi non possono essere vuoti...")
	}

	// Ottieni il paziente
	patient, err := pc.GetPatient(ctx, patientID)

	if err != nil {
		return err
	}

	// Aggiungi la visita
	visit := Visit{
		DoctorID:  doctorID,
		VisitDate: visitDate,
		Symptoms:  symptoms,
		Diagnosis: diagnosis,
	}

	patient.Visits = append(patient.Visits, visit)

	// Aggiorna lo stato del paziente sulla blockchain
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(patientID, patientJSON)
}
