package main

import (
	"backup_db/config"
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

var destinationServerConnection config.ConnectionParams
var backupSourceServerConnection config.ConnectionParams
var isDeleteTableAfterBackup bool
var targetTableName string

func proccessBackup(listBackup []config.BackupList) {
	startProcess := time.Now()
	totalBackup := len(listBackup)
	fmt.Printf("\n\n🤖  : Please wait while I help backup your %v table", totalBackup)
	fmt.Printf("\n\n🤖  : Start at 🕑 %s \n\n", startProcess.Format("2006-01-02 15:04:05"))

	maxProccess := make(chan struct{}, 8)
	var wg sync.WaitGroup

	for _, tbl := range listBackup {
		maxProccess <- struct{}{} // acquire a slot
		wg.Add(1)

		go func(tbl config.BackupList) {
			defer wg.Done()
			defer func() { <-maxProccess }() // release the slot

			if tbl.RowTotal == 0 {
				config.DropTable(backupSourceServerConnection, tbl.TableName)
			} else {
				config.BackupTable(destinationServerConnection, backupSourceServerConnection, tbl.TableName, isDeleteTableAfterBackup)
			}
		}(tbl)
	}

	wg.Wait() // Wait for all to finish

	elapsed := time.Since(startProcess) // Calculate duration
	fmt.Printf("\n\n🤖  : ✅ Done with execution time: 🕑 %s", elapsed)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("🤖 : Welcome to Backup {MS SQL SERVER} table program 👾👾👾...\n")
	backupSourceServerConnection = config.GetInputSourceServer(reader, false)
	listDatabaseBackupSourceServer := config.GetListDatabaseInServer(backupSourceServerConnection)
	backupSourceServerConnection.Database = config.GetInputSourceDB(reader, listDatabaseBackupSourceServer, true, false)

	fmt.Printf("\n 🌐  You're selected Source Backup: %s - %s \n\n", backupSourceServerConnection.Server, backupSourceServerConnection.Database)

	inputTableName, listTableForBackup := config.GetInputTableName(reader, backupSourceServerConnection)
	targetTableName = inputTableName

	if !config.AskingContinueBackupProcess(reader) {
		fmt.Println("\n🤖 : Thank you for interact with me ... ")
		fmt.Println("🤖 : Have A Nice Day ... ")
		return
	}

	destinationServerConnection = config.GetInputSourceServer(reader, true)
	listDatabaseDestinationServer := config.GetListDatabaseInServer(destinationServerConnection)
	destinationServerConnection.Database = config.GetInputSourceDB(reader, listDatabaseDestinationServer, true, true)
	isDeleteTableAfterBackup = config.GetInputIsDeleteAfterBackup(reader)
	fmt.Printf("\n ♻️  Deleting Table After Backup: %v \n\n", isDeleteTableAfterBackup)

	if !config.AskingContinueBackupProcess(reader) {
		fmt.Println("\n🤖 : Thank you for interact with me ... ")
		fmt.Println("🤖 : Have A Nice Day ... ")
		return
	}

	proccessBackup(listTableForBackup)
}
