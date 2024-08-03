package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPersonInfo(hi *gin.Context) {
	personID := hi.Param("person_id")
	db := getDB()
	var name string
	err := db.QueryRow("SELECT name FROM person WHERE id = ?", personID).Scan(&name)
	if err != nil {
		hi.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	var phoneNumber string
	err = db.QueryRow("SELECT number FROM phone WHERE person_id = ?", personID).Scan(&phoneNumber)
	if err != nil {
		hi.JSON(http.StatusNotFound, gin.H{"error": "Phone not found"})
		return
	}
	var addressID int
	err = db.QueryRow("SELECT address_id FROM address_join WHERE person_id = ?", personID).Scan(&addressID)
	if err != nil {
		hi.JSON(http.StatusNotFound, gin.H{"error": "Address join not found"})
		return
	}

	var city, state, street1, street2, zipCode string
	err = db.QueryRow("SELECT city, state, street1, street2, zip_code FROM address WHERE id = ?", addressID).Scan(&city, &state, &street1, &street2, &zipCode)
	if err != nil {
		hi.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}
	person_det := gin.H{
		"name":         name,
		"phone_number": phoneNumber,
		"city":         city,
		"state":        state,
		"street1":      street1,
		"street2":      street2,
		"zip_code":     zipCode,
	}

	c.JSON(http.StatusOK, person_det)
}

func createPerson(gi *gin.Context) {
	a := getDB()

	var input struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		City        string `json:"city"`
		State       string `json:"state"`
		Street1     string `json:"street1"`
		Street2     string `json:"street2"`
		ZipCode     string `json:"zip_code"`
	}

	if err := gi.BindJSON(&input); err != nil {
		gi.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	tx, err := a.Begin()
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()
	r, err := tx.Exec("INSERT INTO person (name, age) VALUES (?, ?)", input.Name, 30) // Age default to 30 for simplicity
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert person"})
		return
	}

	personID, err := r.LastInsertId()
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get person ID"})
		return
	}
	_, err = tx.Exec("INSERT INTO phone (number, person_id) VALUES (?, ?)", input.PhoneNumber, personID)
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert phone"})
		return
	}
	r, err = tx.Exec("INSERT INTO address (city, state, street1, street2, zip_code) VALUES (?, ?, ?, ?, ?)",
		input.City, input.State, input.Street1, input.Street2, input.ZipCode)
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert address"})
		return
	}

	addressID, err := r.LastInsertId()
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get address ID"})
		return
	}
	_, err = tx.Exec("INSERT INTO address_join (person_id, address_id) VALUES (?, ?)", personID, addressID)
	if err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert address join"})
		return
	}
	if err := tx.Commit(); err != nil {
		gi.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}
	gi.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})
}
