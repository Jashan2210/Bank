package main

import (
	"fmt"
	"os"
	"time"
)

type Account struct {
	Name        string
	AccNo       int
	Balance     int
	Active      bool   // Indicates if the account is active or deactivated
	Contact     string // Contact information (e.g., phone number)
	AccountType string // Type of account (e.g., Savings, Checking)
}

var accounts []Account
var transactionID int // Counter to generate unique transaction IDs

func main() {
	loadAccounts()
	for {
		fmt.Println("\n1. Create account\n2. Deposit\n3. Withdraw\n4. Balance\n5. Account details\n6. Deactivate account\n7. Exit")
		var choice int
		fmt.Print("Choice: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			createAccount()
		case 2:
			depositMoney()
		case 3:
			withdrawMoney()
		case 4:
			checkBalance()
		case 5:
			accountDetails()
		case 6:
			deactivateAccount()
		case 7:
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func loadAccounts() {
	file, _ := os.Open("accounts.txt")
	defer file.Close()
	for {
		var acc Account
		var status string
		if _, err := fmt.Fscanf(file, "Name: %s\nAccount Number: %d\nBalance: %d\nStatus: %s\nContact: %s\nAccountType: %s\n\n", &acc.Name, &acc.AccNo, &acc.Balance, &status, &acc.Contact, &acc.AccountType); err != nil {
			break
		}
		// Set the active status based on the string
		if status == "Active" {
			acc.Active = true
		} else {
			acc.Active = false
		}
		accounts = append(accounts, acc)
	}
}

func saveAccounts() {
	file, _ := os.Create("accounts.txt")
	defer file.Close()
	for _, acc := range accounts {
		status := "Deactivated"
		if acc.Active {
			status = "Active"
		}
		fmt.Fprintf(file, "Name: %s\nAccount Number: %d\nBalance: %d\nStatus: %s\nContact: %s\nAccountType: %s\n\n", acc.Name, acc.AccNo, acc.Balance, status, acc.Contact, acc.AccountType)
	}
}

func createAccount() {
	var name, contact, accountType string
	var accNo int
	fmt.Print("Name: ")
	fmt.Scan(&name)
	fmt.Print("AccNo: ")
	fmt.Scan(&accNo)
	fmt.Print("Contact (e.g., phone number): ")
	fmt.Scan(&contact)
	fmt.Print("Account Type (e.g., Savings, Checking): ")
	fmt.Scan(&accountType)

	// Check if account already exists
	for _, acc := range accounts {
		if acc.AccNo == accNo {
			fmt.Println("Account with this account number already exists!")
			return
		}
	}

	// If account does not exist, create it
	accounts = append(accounts, Account{name, accNo, 10, true, contact, accountType}) // Default to active
	saveAccounts()
	createTransactionFile(accNo)
	fmt.Println("Account created successfully!")
}

func depositMoney() {
	var accNo, amount int
	fmt.Print("AccNo: ")
	fmt.Scan(&accNo)
	fmt.Print("Amount: ")
	fmt.Scan(&amount)
	for i := range accounts {
		if accounts[i].AccNo == accNo && accounts[i].Active {
			Before := accounts[i].Balance // Capture the balance before the transaction
			accounts[i].Balance += amount
			saveAccounts()
			logTransaction(accNo, "Credit", amount, Before, accounts[i].Balance)
			return
		}
	}
	fmt.Println("Account not found or deactivated!")
}

func withdrawMoney() {
	var accNo, amount int
	fmt.Print("AccNo: ")
	fmt.Scan(&accNo)
	fmt.Print("Amount: ")
	fmt.Scan(&amount)
	for i := range accounts {
		if accounts[i].AccNo == accNo && accounts[i].Active {
			if accounts[i].Balance >= amount {
				Before := accounts[i].Balance // Capture the balance before the transaction
				accounts[i].Balance -= amount
				saveAccounts()
				logTransaction(accNo, "Debit", amount, Before, accounts[i].Balance)
			} else {
				fmt.Println("Insufficient balance!")
			}
			return
		}
	}
	fmt.Println("Account not found or deactivated!")
}

func checkBalance() {
	var accNo int
	fmt.Print("AccNo: ")
	fmt.Scan(&accNo)
	for i := range accounts {
		if accounts[i].AccNo == accNo && accounts[i].Active {
			fmt.Printf("Balance: %d\n", accounts[i].Balance)
			return
		}
	}
	fmt.Println("Account not found or deactivated!")
}

func accountDetails() {
	var accNo int
	fmt.Print("AccNo: ")
	fmt.Scan(&accNo)
	for i := range accounts {
		if accounts[i].AccNo == accNo {
			status := "Deactivated"
			if accounts[i].Active {
				status = "Active"
			}
			fmt.Printf("Name: %s\nAccNo: %d\nBalance: %d\nStatus: %s\nContact: %s\nAccount Type: %s\n", accounts[i].Name, accounts[i].AccNo, accounts[i].Balance, status, accounts[i].Contact, accounts[i].AccountType)
			return
		}
	}
	fmt.Println("Account not found!")
}

func deactivateAccount() {
	var accNo int
	fmt.Print("AccNo to deactivate: ")
	fmt.Scan(&accNo)
	for i := range accounts {
		if accounts[i].AccNo == accNo {
			accounts[i].Active = false // Deactivate the account
			saveAccounts()
			fmt.Println("Account deactivated!")
			return
		}
	}
	fmt.Println("Account not found!")
}

func createTransactionFile(accNo int) {
	os.Create(fmt.Sprintf("%d_transactions.txt", accNo))
}

func logTransaction(accNo int, transType string, amount, Before, totalBalance int) {
	transactionID++ // Increment the transaction ID for each new transaction
	file, _ := os.OpenFile(fmt.Sprintf("%d_transactions.txt", accNo), os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()
	timestamp := time.Now().Format(time.RFC3339)
	// Format: [Timestamp] AccountNo TransactionID Type Amount Before TotalBalance
	file.WriteString(fmt.Sprintf("[%s] AccountNo: %d, TransactionID: %d, Type: %s, Amount: %d, Before: %d, TotalBalance: %d\n", timestamp, accNo, transactionID, transType, amount, Before, totalBalance))
}
