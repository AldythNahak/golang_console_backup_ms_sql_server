package config

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func getInputUsernameServer(reader *bufio.Reader) {
	fmt.Println("\n🤖 : ")
}

func GetInputSourceServer(reader *bufio.Reader, isDestinationBackup bool) ConnectionParams {
	strQuestion := []string{"What Server you want to backup ?", "Input Source Server (target Server you want to backup): "}

	if isDestinationBackup {
		strQuestion[0] = "What Server you want to stored those data ?"
		strQuestion[1] = "Input Destination Server for Backup: "
	}

	fmt.Println("\n🤖 :", strQuestion[0])
	fmt.Print("➡️  ", strQuestion[1])
	inputSourceServer, _ := reader.ReadString('\n')
	inputSourceServer = strings.TrimSpace(inputSourceServer)
	fmt.Print("🔒 Input Login ID: ")
	inputLoginId, _ := reader.ReadString('\n')
	fmt.Print("🔐 Input Password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("Failed to read password:", err)
	}
	inputPassword := string(passwordBytes)

	connParams := ConnectionParams{
		Server:   inputSourceServer,
		LoginID:  inputLoginId,
		Password: inputPassword,
		Database: "master",
	}

	fmt.Printf("\n🕐 Please Wait while looking for Server [%s] \n", inputSourceServer)
	if !CheckServer(connParams) {
		fmt.Printf("\n⛔️ SERVER [%s] NOT FOUND ! OR Incorect Login ID & Password ! \n", inputSourceServer)
		return GetInputSourceServer(reader, isDestinationBackup)
	}

	return connParams
}

func GetInputSourceDB(reader *bufio.Reader, listBackup []DatabaseList, printListBackup bool, isDestinationBackup bool) string {
	strQuestion := []string{"What Database you want to backup ?", "Input Source Database (target Database you want to backup): "}

	if isDestinationBackup {
		strQuestion[0] = "What Database you want to stored those data ?"
		strQuestion[1] = "Input Destination Database (target Database you want to stored data): "
	}

	fmt.Println("\n\n🤖 :", strQuestion[0])

	if printListBackup {
		fmt.Println("***************************")
		for _, db := range listBackup {
			fmt.Println(fmt.Sprintf("%v. %s", db.Row, db.Database))
		}
		fmt.Println("***************************")
	}

	fmt.Println("❕ Input Number Only ")
	fmt.Print("➡️  ", strQuestion[1])
	inputDatabase, _ := reader.ReadString('\n')
	inputDatabase = strings.TrimSpace(inputDatabase)
	intInputDatabase, err := strconv.Atoi(inputDatabase)
	if err != nil {
		fmt.Println("⛔️ Invalid Input: ", err, "\n")
		return GetInputSourceDB(reader, listBackup, false, isDestinationBackup)
	}

	var selectedDatabase string
	isExistDatabase := false
	for _, db := range listBackup {
		if db.Row != intInputDatabase {
			continue
		}

		selectedDatabase = db.Database
		isExistDatabase = true
	}

	if !isExistDatabase {
		fmt.Println("\n⛔️ Database Selected Not Found ")
		return GetInputSourceDB(reader, listBackup, false, isDestinationBackup)
	}

	return selectedDatabase
}

func GetInputIsDeleteAfterBackup(reader *bufio.Reader) bool {
	fmt.Println("\n🤖 : Do you want to delete source table after Backup ?! ")
	fmt.Print("➡️  Type your answer (y/n): ")
	inputIsDeleteTableAfterBackup, _ := reader.ReadString('\n')
	inputIsDeleteTableAfterBackup = strings.TrimSpace(inputIsDeleteTableAfterBackup)
	inputIsDeleteTableAfterBackup = strings.ToLower(inputIsDeleteTableAfterBackup)

	if inputIsDeleteTableAfterBackup != "y" && inputIsDeleteTableAfterBackup != "n" {
		fmt.Println("⛔️ Invalid Input ! \n")
		return GetInputIsDeleteAfterBackup(reader)
	}

	return inputIsDeleteTableAfterBackup == "y"
}

func GetInputTableName(reader *bufio.Reader, conn ConnectionParams) (string, []BackupList) {
	fmt.Println("\n🤖 : Do you have a specific table name you'd like to back up ? ")
	fmt.Println("❕ [e.g., %%table name%%, table%%name]")
	fmt.Print("➡️  Please enter the specific table name you want to back up (or leave blank to back up all tables): ")
	inputTableName, _ := reader.ReadString('\n')
	inputTableName = strings.TrimSpace(inputTableName)

	listTableForBackup := GetListBackup(conn, inputTableName)

	if len(listTableForBackup) == 0 {
		fmt.Printf("\n\n⚠️  No Table Found at %s - %s With Name LIKE '%s'", conn.Server, conn.Database, inputTableName)
		return GetInputTableName(reader, conn)
	}

	return GetPreviewListTableForBackup(reader, listTableForBackup, conn, inputTableName)
}

func GetPreviewListTableForBackup(reader *bufio.Reader, listTableForBackup []BackupList, conn ConnectionParams, tableName string) (string, []BackupList) {
	totalTable := len(listTableForBackup)
	fmt.Printf("\n\n🤖 : Do you want to preview %v list table ? \n", totalTable)
	fmt.Println("❔ y: preview, n: next process, x: change table name")
	fmt.Print("➡️  Type your answer (y/n/x): ")
	inputPreview, _ := reader.ReadString('\n')
	inputPreview = strings.TrimSpace(inputPreview)
	inputPreview = strings.ToLower(inputPreview)

	if inputPreview != "y" && inputPreview != "n" && inputPreview != "x" {
		fmt.Println("⛔️ Invalid Input ! \n")
		return GetPreviewListTableForBackup(reader, listTableForBackup, conn, tableName)
	}

	if inputPreview == "x" {
		return GetInputTableName(reader, conn)
	}

	if inputPreview == "y" {
		fmt.Println("📃 Showing ", totalTable, " Table ...")
		fmt.Print("\n💎💎💎💎💎💎💎💎💎💎💎💎💎")
		for _, tbl := range listTableForBackup {
			fmt.Printf("\n💎 %v .  %s | %v rows | %v MB | created: %s", tbl.Row, tbl.TableName, tbl.RowTotal, tbl.TotalSpacedMB, tbl.DateCreated)
		}
		fmt.Println("\n💎💎💎💎💎💎💎💎💎💎💎💎💎")
		fmt.Println("\n🤖 : Continue Proccess ? ")
		fmt.Println("❔ y: next process, x: change table name, t: select table from list")
		fmt.Print("➡️  Type your answer (y/x/t): ")
		inputPreview, _ := reader.ReadString('\n')
		inputPreview = strings.TrimSpace(inputPreview)
		inputPreview = strings.ToLower(inputPreview)

		if inputPreview != "y" && inputPreview != "x" && inputPreview != "t" {
			fmt.Println("⛔️ Invalid Input ! \n")
			return GetPreviewListTableForBackup(reader, listTableForBackup, conn, tableName)
		}

		if inputPreview == "x" {
			return GetInputTableName(reader, conn)
		}

		if inputPreview == "t" {
			newListTableForBackup := getSelectedTableFromList(reader, listTableForBackup)
			return GetPreviewListTableForBackup(reader, newListTableForBackup, conn, tableName)
		}
	}

	return tableName, listTableForBackup
}

func getSelectedTableFromList(reader *bufio.Reader, listTableForBackup []BackupList) []BackupList {
	var newListTableForBackup []BackupList

	fmt.Println("❔ 't' [type number only separate by ,] exp: 1,2,3,10,15")
	fmt.Print("➡️  Type your selected table: ")
	inputSelectedTable, _ := reader.ReadString('\n')
	inputSelectedTable = strings.TrimSpace(inputSelectedTable)

	if inputSelectedTable == "" {
		return listTableForBackup
	}

	listTableSelected := strings.Split(inputSelectedTable, ",")

	fmt.Println("🕘 Please Wait for set your selected table 🕘")

	selected := make(map[int]bool)

	for _, order := range listTableSelected {
		numOrder, err := strconv.Atoi(order)
		if err != nil {
			continue
		}
		selected[numOrder] = true
	}

	for _, tbl := range listTableForBackup {
		if selected[tbl.Row] {
			newListTableForBackup = append(newListTableForBackup, tbl)
		}
	}

	return newListTableForBackup
}

func AskingContinueBackupProcess(reader *bufio.Reader) bool {
	fmt.Println("\n\n🤖 : Are you sure to continue backup process ? ")
	fmt.Print("➡️  Type your answer (y/n): ")

	inputPermission, _ := reader.ReadString('\n')
	inputPermission = strings.TrimSpace(inputPermission)
	inputPermission = strings.ToLower(inputPermission)

	if inputPermission != "y" && inputPermission != "n" {
		fmt.Println("⛔️ Invalid Input ! \n")
		return AskingContinueBackupProcess(reader)
	}

	return inputPermission == "y"
}
