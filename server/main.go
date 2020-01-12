package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	jwt "github.com/dgrijalva/jwt-go"
	"strconv"
	
)

var mySigningKeyone = []byte("password1")  //User1's password 
var mySigningKeytwo = []byte("password2")  //User2's password 

// Users ID
var myclientidone string = "1"
var myclientidtwo string = "2"

var idlist []string // List of users ID
var clientnum string //User ID collected after signing in and used to authorize or block actions required by the user

// Declaration of the fields used throughout the project

type Certificate struct {
    Id       string   `json:"id,omitempty"`
    Title 	  string   `json:"title,omitempty"`
    Createdat  string   `json:"createdat,omitempty"`
	Ownerid    string    `json:"ownerid,omitempty"`
    Year   		string `json:"year,omitempty"`
	Note        string  `json:"note,omitempty"`
	Transfer    string   `json:"transfer,omitempty"`
}
type Transfer struct {
	Id	string	`json:"id,omitempty"`
    To  string `json:"to,omitempty"`
    Status string `json:"status,omitempty"`
}

type User struct {
    Userid  string  `json:"userid,omitempty"`
	Email	string	`json:"email,omitempty"`
    Name    string 	`json:"name,omitempty"`
}

var certif []Certificate
var trans []Transfer
var owner []User
var allcertif []Certificate //Creation of the slice which contains the user's certificates in the 'GetAllCertif' function

// Creation of the home page the user accesses after signing in
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
	fmt.Println("Endpoint Hit: homePage")
}

// Authorization function to allow the user to sign in
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
	
	var N = len(idlist) //Size of the users list
	
	if r.Header["Clientid"] != nil { //Let's first check the user filled the 'Clientid' header key
		//We verify whether or not the ID submitted is valid
		var verif string = "not ok"
		for i:=0;i<N;i++{
			if r.Header["Clientid"][0]==idlist[i]{
				verif="ok" 
			}
		}
		
		//Verification of the encrypted password submitted by the user
		if verif=="ok"{ //The Client ID is valid, let's verify the password is correct and matches the client ID
			if r.Header["Clientid"][0]=="1" { //First case: the user signed in as user 1
				if r.Header["Token"] != nil { //Once again, a token is required to continue the authentication process
				
						token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
							if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
								return nil, fmt.Errorf("There was an error")
							}
								return mySigningKeyone, nil	
						})
						
						if err != nil {
							fmt.Fprintf(w, err.Error())
						}
						
						//if the 'token' is valid, user 1 has now signed in
						if token.Valid {
							clientnum=r.Header["Clientid"][0] 
							endpoint(w, r)
						}
					}	else {
							fmt.Fprintf(w, "Not Authorized")
					}
			}	else { //Second case: the user signed in as user 2
				if r.Header["Token"] != nil { //Once again, a token is required to continue the authentication process
				
						token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
							if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
								return nil, fmt.Errorf("There was an error")
							}
								return mySigningKeytwo, nil	
						})
						
						if err != nil {
							fmt.Fprintf(w, err.Error())
						}
						//if the 'token' is valid, user 2 has now signed in
						if token.Valid {
							clientnum=r.Header["Clientid"][0] 
							endpoint(w, r)
						}
					}	else {
				
							fmt.Fprintf(w, "Not Authorized")
					}
			}
		}	else { 
				fmt.Fprintf(w, "This ID does not exist")
		}
					
		} else {
			fmt.Fprintf(w, "Insert ID")
		}
	})
}
	
// Function to get one certificate 
func GetCertif(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	verifi := "not ok"
	
	//We verify whether the required certificate exists or not
	for _, itemm := range certif {
		if itemm.Id==params["id"] {
			verifi="ok"
		}
	}
	
	if verifi=="ok" {
		for _, item := range certif {
			if item.Id == params["id"] && item.Ownerid==clientnum  { //Display information only if the owner ID of the certificate is the same as the connected user ID
				json.NewEncoder(w).Encode(item)
				return
			}	else if item.Id == params["id"] && item.Ownerid!=clientnum {
				fmt.Fprintf(w, "You can't access this certificate.") //If the certificate does not belong to the connected user: return this message
			}
		}
	}	else {
		fmt.Fprintf(w, "This certificate does not exist.")
	}
}

// Function to create a certificate
func CreateCertif(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    var certifi Certificate
	
    _ = json.NewDecoder(r.Body).Decode(&certifi) //Collect information entered in the 'Body' section by the user 
	
	//Whatever the certificate ID entered by the user, we create the certificate with the next available ID (ex.: if 5 certificates in the database: available ID = 6)
	var taille = len(certif) //Number of certificates
	taille=taille+1	//Available ID
	t:=strconv.Itoa(taille) //Convert an integer into string
	
	//Define Certificate ID and Owner ID, add it to the existing certificates
	certifi.Id = t
	certifi.Ownerid = clientnum
	certifi.Transfer = "Nil"
	certif=append(certif, certifi)
	//Then, we inform the user
	if params["id"] != t {
		fmt.Fprintf(w, "This certificate number was not available, we automatically created your new certificate with the ID:"+" "+t)
	}
}

//Function to display all the user's certificates
func GetAllCertif(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	allcertif=allcertif[:0]
	//Verification of the access authorization and collect of the user's certificates
	
	if params["userid"]==clientnum{
		for _, item := range certif {
			//We add the user's certifications in 'allcertif'
			if item.Ownerid==clientnum  {
				allcertif=append(allcertif, item)
			}
		}
		
		//We display the user's information for more clarity
		for _, itemm := range owner {
			if itemm.Userid==clientnum  {
			fmt.Fprintf(w, "Client information: ")
			json.NewEncoder(w).Encode(itemm)
			fmt.Fprintf(w, "\n")
			}
		}
		
		//We display the user's certificates
		fmt.Fprintf(w, "Certificates of the client: ")
		json.NewEncoder(w).Encode(allcertif)
		
	}	else {
	fmt.Fprintf(w, "You can't access this information.")
	}
}

 //Function to update a certificate
func UpdateCertif(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    for _, item := range certif {
        if item.Id == params["id"] && item.Ownerid==clientnum  {
			sav:=item.Transfer
            _ = json.NewDecoder(r.Body).Decode(&item) //We collect the new information provided by the user
			
			//The certificate and owner IDs as well as the Transfer status stay the same
			item.Id = params["id"] 
			item.Ownerid = clientnum
			item.Transfer=sav
			
			//The previous version of the certificate is erased 
			for index, itemm := range certif {
				if itemm.Id == params["id"] {
				certif = append(certif[:index], certif[index+1:]...)
				break
				}
			}
			
			//The new version is now added and displayed 
			certif=append(certif, item)
			json.NewEncoder(w).Encode(item)
            return
			
        }	else if item.Id == params["id"] && item.Ownerid!=clientnum {
			fmt.Fprintf(w, "You can't modify this certificate.") //Displayed message if attempt to modify a certificate which does not belong to the connected user
		}
    }
}

//Function to delete a certificate
func DeleteCertif(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, item := range certif {
        if item.Id == params["id"] && item.Ownerid==clientnum {
				certif = append(certif[:index], certif[index+1:]...)
				break
		}	else if item.Id == params["id"] && item.Ownerid!=clientnum { 
			fmt.Fprintf(w, "You are not allowed to delete this certificate.")
		}
	}
}

//Function to create a transfer
func CreateTransfer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var v = params["id"] //We use a simpler variable as the two 'Transfer' functions are fairly long
	
	for _, item := range certif {
		//Let's first check the user has the right to transfer this certificate
        if item.Id == v && item.Ownerid==clientnum  { 
			checkk:="ok"
			
			//We browse the transfers list to verify there is no transfer for this certificate
			for _, itemm := range trans {
				if itemm.Id == v {
					checkk="not ok"
				}
			}
			//If so, we inform the user
			if checkk=="not ok" {
				fmt.Fprintf(w, "There is already a transfer for this certificate which is pending.")
			} else {
				var transcertif Transfer //Creation of a new transfer
				
				//We collect the receiving party's email address 
				_ = json.NewDecoder(r.Body).Decode(&transcertif) 
				
				transcertif.Id=v //Certificate ID stays unchanged 
				transcertif.Status="Pending" //The transfer status is now changed until the answer of the receiving party
				
				//We verify if the email address of the receiver is correct and is not the current user's 
				check:="not ok" 
				for _, ite := range owner {
					if ite.Email==transcertif.To && ite.Userid!=clientnum {
						check="ok"
					}
				}
				
				if check=="ok" { //Then, if the email is valid, the transfer proposal is sent
						trans=append(trans, transcertif) //The transfer is added to the transfers list
						//The 'transfer' status of the certificate is changed
						for index, itp := range certif {
							if itp.Id==v {
								certif[index].Transfer="Transfer pending"
							}
						}
						fmt.Fprintf(w, "Transfer proposal sent.")
						
				}	else { //If the email is not valid, the user is informed
						fmt.Fprintf(w, "This email is not valid.")
				}
			}
		} else if item.Id == v && item.Ownerid!=clientnum {
			fmt.Fprintf(w, "You are not allowed to transfer this certificate.")
		}
	}
}

//Function to manage the transfer proposals
func ManageTransfer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var v = params["id"]
	verification := "not ok"
	var emailuser string
	
	//We collect the email address of the connected user
	for _, item := range owner {
		if item.Userid == clientnum {
			emailuser = item.Email
		}
	}
	
	//We browse the transfers list to make sure there is a transfer proposal for the certificate 
	for _, itemm := range trans {
		if itemm.Id == v && itemm.To == emailuser { //This verification is done thanks to the email address
			verification = "ok"
		}
	}
	
	//If so, the answer of the receiver is collected. If not, a message indicates there is no action needed
	if verification == "ok" {
		var response Transfer
		
		//Data collection of the answer: only the field "status" of 'response' will be modified here
		_ = json.NewDecoder(r.Body).Decode(&response) 
		response.Id=v
		response.To=emailuser
		
		//Two acceptable answers: 'Declined' or 'Accepted'. If anything else: the user is told he/she has to enter one of these two possibilities 
		if (response.Status=="Declined" || response.Status=="Accepted") {
		
			//If the transfer is accepted:
			if response.Status == "Accepted" {
			
				//Creation of a new certificate identical to the old one but containing the new owner ID
				var newcertif Certificate
				newcertif.Id = v
				newcertif.Transfer="Transferred" //'Transferred' to indicate it has been transferred at least once
				newcertif.Ownerid=clientnum
				for _, iitemm:=range certif {
						if iitemm.Id == v {
							newcertif.Title=iitemm.Title
							newcertif.Createdat=iitemm.Createdat
							newcertif.Year=iitemm.Year
							newcertif.Note=iitemm.Note	
						}
				}
				
				//Removal of the old certificate
				for index, iitem:= range certif {
					if iitem.Id == v {
						certif = append(certif[:index], certif[index+1:]...)
						break
					}
				}
				
				//The new certificate with the relevant owner ID is added and replaces the old one
				certif=append(certif, newcertif)
				
				//Return all user's certificates
				allcertif=allcertif[:0]
				for _, itte := range certif {
				//We add the user's certifications in 'allcertif'
					if itte.Ownerid==clientnum  {
						allcertif = append(allcertif, itte)
					}
				}
				json.NewEncoder(w).Encode(allcertif)
			}	else {
				//The 'transfer' status of the certificate is changed
				for index, itpp := range certif {
					if itpp.Id==v {
						certif[index].Transfer="Nil"
					}
				}
			}
			
			//Finally, the transfer proposal is removed from the transfers list
			for index, ite := range trans {
				if ite.Id == v {
						trans = append(trans[:index], trans[index+1:]...)
						break
				}
			}
				
		} else {
			fmt.Fprintf(w, "Please accept (type: 'Accepted') or decline (type: 'Declined') the transfer.")
		}
	} else {
		fmt.Fprintf(w, "No answer needed from you for this certificate.")
	}
}



func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	idlist=append(idlist, myclientidone, myclientidtwo) //It fills the 'Client ID' list
	
	//Creation of two users
	owner = append(owner, User{Userid: "1", Email: "gm@gmail.com", Name: "Guillaume Mendlevitch"})
	owner = append(owner, User{Userid: "2", Email: "so@gmail.com", Name: "Samuel Osina"})
	
	//Creation of two certificates each
    certif = append(certif, Certificate{Id: "1", Title: "Mona Lisa", Createdat: "15 feb 2017", Ownerid: "1", Year: "1503", Note: "Not damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "3", Title: "The Raft of the Medusa", Createdat: "5 oct 2016", Ownerid: "1", Year: "1818", Note: "A bit damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "2", Title: "The Wedding Feast at Cana", Createdat: "15 oct 2017", Ownerid: "2", Year: "1563", Note: "Not damaged", Transfer: "Nil"})
	certif = append(certif, Certificate{Id: "4", Title: "Sunflowers", Createdat: "6 Jan 2004", Ownerid: "2", Year: "1880", Note: "A bit damaged", Transfer: "Nil"})
		
	router.Handle("/", isAuthorized(homePage))
	
	router.HandleFunc("/certificates/{id}", GetCertif).Methods("GET")
	router.HandleFunc("/certificates/{id}", CreateCertif).Methods("POST")
	router.HandleFunc("/certificates/{id}", UpdateCertif).Methods("PATCH")
	router.HandleFunc("/certificates/{id}", DeleteCertif).Methods("DELETE")
	
	router.HandleFunc("/users/{userid}/certificates", GetAllCertif).Methods("GET")
	
	router.HandleFunc("/certificates/{id}/transfers", CreateTransfer).Methods("POST")
	router.HandleFunc("/certificates/{id}/transfers", ManageTransfer).Methods("PATCH")
	
	
	log.Fatal(http.ListenAndServe(":9000", router))
}


func main() {
	
	handleRequests()
}