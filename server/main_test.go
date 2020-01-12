package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
	//"fmt"
	"encoding/json"
	//"log"
	//jwt "github.com/dgrijalva/jwt-go"
	//"strconv"
	"bytes"
)

func init() {

}

func Router() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/certificates/{id}", GetCertif).Methods("GET")
	router.HandleFunc("/certificates/{id}", CreateCertif).Methods("POST")
	router.HandleFunc("/users/{userid}/certificates", GetAllCertif).Methods("GET")
	router.HandleFunc("/certificates/{id}", UpdateCertif).Methods("PATCH")
	router.HandleFunc("/certificates/{id}", DeleteCertif).Methods("DELETE")
	router.HandleFunc("/certificates/{id}/transfers", CreateTransfer).Methods("POST")
	router.HandleFunc("/certificates/{id}/transfers", ManageTransfer).Methods("PATCH")
    return router
}

//Test for 'GetCertif' function 
func TestGetCertif(t *testing.T) {
	clientnum="1" //Let's pretend user 1 is connected
	
	//The certificates database contains only one certificate
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	
	//Execution of the 'Get Certificate 1' request
    request, _ := http.NewRequest("GET", "/certificates/1", nil)
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	oricertif:=certif[0] //'oricertif' is the only element of the 'certif' database, it is now the original certificate
	
	//Returned response saved in 'm'
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	
	//One compares the original certificate with the collected one
	

	if m["id"] != oricertif.Id {
        t.Errorf("Expected certificate Id to be '%v'. Got '%v'", oricertif.Id, m["id"])
    }
	if m["title"] != oricertif.Title {
        t.Errorf("Expected certificate title to be '%v'. Got '%v'", oricertif.Title, m["title"])
    }
	if m["createdat"] != oricertif.Createdat {
        t.Errorf("Expected creation date of the certificate to be '%v'. Got '%v'", oricertif.Createdat, m["createdat"])
    }
	if m["ownerid"] != oricertif.Ownerid {
        t.Errorf("Expected owner Id to be '%v'. Got '%v'", oricertif.Ownerid, m["ownerid"])
    }	
	if m["year"] != oricertif.Year {
        t.Errorf("Expected creation year of the artwork to be '%v'. Got '%v'", oricertif.Year, m["year"])
    }
	if m["note"] != oricertif.Note {
        t.Errorf("Expected note concerning the artwork to be '%v'. Got '%v'", oricertif.Note, m["note"])
    }
	if m["transfer"] != oricertif.Transfer {
        t.Errorf("Expected transfer status to be '%v'. Got '%v'", oricertif.Transfer, m["transfer"])
    }	
}


//Test for the 'CreateCertif' function
func TestCreateCertif(t *testing.T) {
	clientnum="1"
	
	//Creation of the certificate based on the 'Certificate' structure 
	certificate := &Certificate{
		Id: "1", 
		Title: "Mona Lisa", 
		Createdat: "15 feb 2017", 
		Ownerid: "1", 
		Year: "1503", 
		Note: "Not damaged", 
		Transfer: "Nil"}

	//POST request for certificate 1 with the above information implemented in the 'body' in order to create 'certificate 1'
	jsonCertificate, _ := json.Marshal(certificate)
    request, _ := http.NewRequest("POST", "/certificates/1", bytes.NewBuffer(jsonCertificate))
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	m := certif[0] //One collects the certificate which was just sent
	
	//Comparison of the 'input' data and 'output' data
	if m.Id != certificate.Id {
        t.Errorf("Expected certificate Id to be '%v'. Got '%v'", certificate.Id, m.Id)
    }
	
	if m.Title != certificate.Title {
        t.Errorf("Expected certificate title to be '%v'. Got '%v'", certificate.Title, m.Title)
    }
	if m.Createdat != certificate.Createdat {
        t.Errorf("Expected creation date of the certificate to be '%v'. Got '%v'", certificate.Createdat, m.Createdat)
    }
	if m.Ownerid != certificate.Ownerid {
        t.Errorf("Expected owner Id to be '%v'. Got '%v'", certificate.Ownerid, m.Ownerid)
    }	
	if m.Year != certificate.Year {
        t.Errorf("Expected creation year of the artwork to be '%v'. Got '%v'", certificate.Year, m.Year)
    }
	if m.Note != certificate.Note {
        t.Errorf("Expected note concerning the artwork to be '%v'. Got '%v'", certificate.Note, m.Note)
    }
	if m.Transfer != certificate.Transfer {
        t.Errorf("Expected transfer status to be '%v'. Got '%v'", certificate.Transfer, m.Transfer)
    }
}

//Test for 'GetAllcertif' function
func TestGetAllCertif(t *testing.T) {
	clientnum="1"
	
	//2 certificates for each user are added
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "3", Title: "The Raft of the Medusa", Createdat: "5 oct 2016", Ownerid: "1", Year: "1818", Note: "A bit damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "2", Title: "The Wedding Feast at Cana", Createdat: "15 oct 2017", Ownerid: "2", Year: "1563", Note: "Not damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "4", Title: "Sunflowers", Createdat: "6 Jan 2004", Ownerid: "2", Year: "1880", Note: "A bit damaged", Transfer: "Nil"})
	
	//The user 1's certificates are supposed to be saved in the 'allcertif' slice and displayed thanks to the following commands:
	request, _ := http.NewRequest("GET", "/users/1/certificates", nil)
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	//'m' and 'n' are respectively the first and second certificate of user 1 collected from 'allcertif' after the above request was sent 
	m:=allcertif[0]
	n:=allcertif[1]
	
	//'mbefore' and 'nbefore' are respectively the first and second original certificate of user 1, directly added in the 'certif' database 

	mbefore:=certif[0]
	nbefore:=certif[1]

	//We compare the data of each certificate to make sure the function returns the correct user 1's certificates
	if m.Id != mbefore.Id {
        t.Errorf("Expected certificate Id to be '%v'. Got '%v'", mbefore.Id, m.Id)
    }
	if m.Title != mbefore.Title {
        t.Errorf("Expected certificate title to be '%v'. Got '%v'", mbefore.Title, m.Title)
    }
	if m.Createdat != mbefore.Createdat {
        t.Errorf("Expected creation date of the certificate to be '%v'. Got '%v'", mbefore.Createdat, m.Createdat)
    }
	if m.Ownerid != mbefore.Ownerid {
        t.Errorf("Expected owner Id to be '%v'. Got '%v'", mbefore.Ownerid, m.Ownerid)
    }	
	if m.Year != mbefore.Year {
        t.Errorf("Expected creation year of the artwork to be '%v'. Got '%v'", mbefore.Year, m.Year)
    }
	if m.Note != mbefore.Note {
        t.Errorf("Expected note concerning the certificate to be '%v'. Got '%v'", mbefore.Note, m.Note)
    }
	if m.Transfer != mbefore.Transfer	{
        t.Errorf("Expected transfer status to be '%v'. Got '%v'", mbefore.Transfer, m.Transfer)
    }
	
	
	if n.Id != nbefore.Id {
        t.Errorf("Expected certificate Id to be '%v'. Got '%v'", nbefore.Id, n.Id)
    }
	if n.Title != nbefore.Title {
        t.Errorf("Expected certificate title to be '%v'. Got '%v'", nbefore.Title, n.Title)
    }
	if n.Createdat != nbefore.Createdat {
        t.Errorf("Expected creation date of the certificate to be '%v'. Got '%v'", nbefore.Createdat, n.Createdat)
    }
	if n.Ownerid != nbefore.Ownerid {
        t.Errorf("Expected owner Id to be '%v'. Got '%v'", nbefore.Ownerid, n.Ownerid)
    }	
	if n.Year != nbefore.Year {
        t.Errorf("Expected creation year of the artwork to be '%v'. Got '%v'", nbefore.Year, n.Year)
    }
	if n.Note != nbefore.Note {
        t.Errorf("Expected note concerning the certificate to be '%v'. Got '%v'", nbefore.Note, n.Note)
    }
	if n.Transfer != nbefore.Transfer	{
        t.Errorf("Expected transfer status to be '%v'. Got '%v'", nbefore.Transfer, n.Transfer)
    }
}


//Test for 'UpdateCertif' function
func TestUpdateCertif(t *testing.T) {
	clientnum="1"
	//One certificate saved and meant to be updated 
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	
	//The original certificate is saved
	originalcertif:=certif[0]
	
	//Updated information are given below
	certificate := &Certificate{
		Id: "1", 
		Title: "La Joconde", 
		Createdat: "15 feb 2017", 
		Ownerid: "1", 
		Year: "1503", 
		Note: "Some damages", 
		Transfer: "Transferred"}
	
	//Request made including the updated information
	jsonCertificate, _ := json.Marshal(certificate)
    request, _ := http.NewRequest("PATCH", "/certificates/1", bytes.NewBuffer(jsonCertificate))
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	//Collection of the new certificate 
	modifiedcertif:=certif[0]
	
	//Comparison between the original and the new certificate
	
	if modifiedcertif.Id != originalcertif.Id {
        t.Errorf("Expected the certificate ID to remain the same (%v). Got %v", originalcertif.Id, modifiedcertif.Id)
    }
	if modifiedcertif.Title != certificate.Title {
		t.Errorf("Expected the title to change from '%v' to '%v'. Got '%v'", originalcertif.Title, certificate.Title, modifiedcertif.Title)
    }
	if modifiedcertif.Createdat != certificate.Createdat {
        t.Errorf("Expected the creation date to change from '%v' to '%v'. Got '%v'", originalcertif.Createdat, certificate.Createdat, modifiedcertif.Createdat)
    }
	if modifiedcertif.Ownerid != originalcertif.Ownerid {
        t.Errorf("Expected the owner ID to remain the same (%v). Got %v", originalcertif.Ownerid, modifiedcertif.Ownerid)
    }	
	if modifiedcertif.Year != certificate.Year {
        t.Errorf("Expected the creation year of the artwork to change from '%v' to '%v'. Got '%v'", originalcertif.Year, certificate.Year, modifiedcertif.Year)
    }
	if modifiedcertif.Note != certificate.Note {
        t.Errorf("Expected the additional notes to change from '%v' to '%v'. Got '%v'", originalcertif.Note, certificate.Note, modifiedcertif.Note)
    }
	if modifiedcertif.Transfer != originalcertif.Transfer {
        t.Errorf("Expected the transfer status to remain the same (%v). Got %v", originalcertif.Transfer, modifiedcertif.Transfer)
    }
}

//Test for the 'DeleteCertif' function
func TestDeleteCertif(t *testing.T) {

	clientnum="1"
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	
	//First step: deleting the above certificate from the database
	request, _ := http.NewRequest("DELETE", "/certificates/1", nil)
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	//Second step: returning the above certificate
	request, _ = http.NewRequest("GET", "/certificates/1", nil)
	response = httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	
	//Verification: the certificate is not supposed to be collected but a message mentioning that the certificate does not exit is displayed
	assert.Equal(t, "This certificate does not exist.", response.Body.String(), "It is supposed to be deleted.")

}

//Test for the 'CreateTransfer' function
func TestCreateTransfer(t *testing.T) {
	clientnum="1"
	
	//Creation of one certificate and two users 
	
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	owner = append(owner, User{Userid: "1", Email: "gm@gmail.com", Name: "Guillaume Mendlevitch"})
	owner = append(owner, User{Userid: "2", Email: "so@gmail.com", Name: "Samuel Osina"})
	
	//Choose the receiver's email in order to send him the certificate
	transf := &Transfer{
		To: "so@gmail.com"}
	
	//Request containing the receiver's information is sent 
	jsonTransfer, _ := json.Marshal(transf)
	request, _ := http.NewRequest("POST", "/certificates/1/transfers", bytes.NewBuffer(jsonTransfer))
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	newtrans:=trans[0] //New transfer added and saved in 'newtrans'
	certifi:=certif[0] //The sent certificate is saved in 'certifi'
	
	//Let's check if the transfer proposal contains correct information
	
	if newtrans.Id != certifi.Id {
        t.Errorf("Expected the certificate ID to remain the same (%v). Got %v", certifi.Id, newtrans.Id)
    }
	if newtrans.To != transf.To {
		t.Errorf("Expected the receiver's email to be '%v'. Got '%v'", transf.To, newtrans.To)
    }
	if newtrans.Status != "Pending" {
		t.Errorf("Expected the status of the transfer to be '%v'. Got '%v'", "Pending", newtrans.Status)
    }
	
}


//Test for the 'ManageTransfer' function
func TestManageTransfer(t *testing.T) {

	clientnum="2" //Connected as user 2 to accept/decline transfers
	
	//Creation of the certificates transferred by user 1
	certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "3", Title: "The Raft of the Medusa", Createdat: "5 oct 2016", Ownerid: "1", Year: "1818", Note: "A bit damaged", Transfer: "Nil"})
	
	owner = append(owner, User{Userid: "1", Email: "gm@gmail.com", Name: "Guillaume Mendlevitch"})
	owner = append(owner, User{Userid: "2", Email: "so@gmail.com", Name: "Samuel Osina"})
	
	//Transfers creation
	trans = append(trans, Transfer{Id: "1", To: "so@gmail.com", Status: "Pending"})
	trans = append(trans, Transfer{Id: "3", To: "so@gmail.com", Status: "Pending"})
	
//-----------------First case: 'Mona Lisa' transfer accepted-------------------------------------
	
	//First step: accepting the transfer of certificate N°1
	
	transf := &Transfer{
		Status: "Accepted"}
	
	jsonTransfer, _ := json.Marshal(transf)
	
	request, _ := http.NewRequest("PATCH", "/certificates/1/transfers", bytes.NewBuffer(jsonTransfer))
    response := httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	//Second step: returning the certificate N°1 and verifying the owner ID has changed
	request, _ = http.NewRequest("GET", "/certificates/1", nil)
    response = httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	var firstcerti map[string]string
	
	firstuser:=owner[0] //User 1 (sender) information
	seconduser:=owner[1] //User 2 (receiver) information
	
	json.Unmarshal(response.Body.Bytes(), &firstcerti) //Collect certificate N°1 after user 2 gave an answer
		
	//Verify the owner's ID has changed
	if firstcerti["ownerid"]!=seconduser.Userid {
		t.Errorf("Expected the owner's ID to change from '%v' to '%v'. Got '%v'", firstuser.Userid, seconduser.Userid, firstcerti["ownerid"])
	}
	
//-----------------Second case: 'The Raft of the Medusa' transfer declined-------------------------------------

	//First step: declining the transfer of certificate N°3

	transf = &Transfer{
		Status: "Declined"}
	
	jsonTransfer, _ = json.Marshal(transf)
	
	request, _ = http.NewRequest("PATCH", "/certificates/3/transfers", bytes.NewBuffer(jsonTransfer))
    response = httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	//Second step: returning the certificate N°3 and verifying the owner ID hasn't changed

	clientnum="1" //If everything went well, the owner ID is still '1', thus, let's connect as user 1
	
	request, _ = http.NewRequest("GET", "/certificates/3", nil)
    response = httptest.NewRecorder()
    Router().ServeHTTP(response, request)
	
	var secondcerti map[string]string

	json.Unmarshal(response.Body.Bytes(), &secondcerti)
	
	//Verify the owner's ID hasn't changed
	if secondcerti["ownerid"]!=firstuser.Userid {
		t.Errorf("Expected the owner's ID to remain the same ('%v'). Got '%v'", firstuser.Userid, secondcerti["ownerid"])
	}
}


