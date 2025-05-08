package config

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type BackupList struct {
	Row           int     `db:"orders"`
	TableName     string  `db:"table_name"`
	RowTotal      int     `db:"row_counts"`
	DateCreated   string  `db:"date_created"`
	TotalSpacedMB float64 `db:"total_spaced_MB"`
	UsedSpacedMB  float64 `db:"used_spaced_MB"`
	IsDoneBackup  int     `db:"is_done_backup"`
}

type DatabaseList struct {
	Row      int    `db:"order"`
	Database string `db:"db_name"`
}

func GetListBackup(conn ConnectionParams, filterTableNameLike string) []BackupList {
	fmt.Println("📢 log: Preapare query get list table for backup in ", conn.Server, "-", conn.Database)
	// serverDB := fmt.Sprintf("[%s].%s", SourceTargetServer, SourceTargetDB)
	serverDB := conn.Database

	if filterTableNameLike != "" {
		filterTableNameLike = fmt.Sprintf(`WHERE name LIKE '%%%s%%'`, filterTableNameLike)
	}

	queryListBackup := fmt.Sprintf(`SELECT DISTINCT 0 AS orders,
		A.name AS table_name,
		B.rows AS row_counts,
		CONVERT(VARCHAR, A.create_date, 23) AS date_created,
		CAST(ROUND(((SUM(D.total_pages) * 8) / 1024.00), 2) AS NUMERIC(36, 2)) AS total_spaced_MB,
		CAST(ROUND(((SUM(D.used_pages) * 8) / 1024.00), 2) AS NUMERIC(36, 2)) AS used_spaced_MB,
		0 AS is_done_backup
	FROM %s.sys.tables AS A 
		LEFT JOIN %s.sys.partitions AS B ON A.object_id = B.object_id
		LEFT JOIN %s.sys.dm_db_index_usage_stats AS C ON C.object_id = A.object_id
		LEFT JOIN %s.sys.allocation_units AS D ON D.container_id = B.partition_id
	%s
	GROUP BY A.name, B.rows, CONVERT(VARCHAR, A.create_date, 23), CONVERT(VARCHAR, C.LAST_USER_UPDATE, 23), CONVERT(VARCHAR,C.LAST_USER_SCAN, 23)
	ORDER BY total_spaced_MB DESC, date_created ASC`, serverDB, serverDB, serverDB, serverDB, filterTableNameLike)

	fmt.Println("📢 log: Connect ", conn.Server, "-", conn.Database)
	dbConn := ConnectDB(conn)

	fmt.Println("🚀 Execute Query get list table for backup in ", conn.Server, "-", conn.Database)
	rowData := ExecQuery(dbConn, queryListBackup, true)

	var collectData []BackupList
	rowCount := 0
	for rowData.Next() {
		var structBck BackupList
		err := rowData.StructScan(&structBck) // Scan each row into dynamic struct
		if err != nil {
			log.Println("⛔️ Row scan error:", err)
			break
		}
		rowCount++
		structBck.Row = rowCount
		collectData = append(collectData, structBck)
	}

	dbConn.Close()
	rowData.Close()

	fmt.Println(fmt.Sprintf("❇️ total list backup: %s", strconv.Itoa(len(collectData))))

	return collectData
}

func BackupTable(destinationConn ConnectionParams, backupSourceConn ConnectionParams, targetTable string, isDropSourceTable bool) {
	targetTable = strings.TrimSpace(targetTable)

	if targetTable == "" {
		fmt.Println("⛔️ Target Table should not be empty ⛔️")
		return
	}

	fmt.Println("📢 log: Prepare Query Backup Table", targetTable)
	targetData := fmt.Sprintf("[%s].%s.dbo.%s", backupSourceConn.Server, backupSourceConn.Database, targetTable)
	queryBackupTable := fmt.Sprintf(`IF OBJECT_ID('%s') IS NOT NULL
	BEGIN
		DROP TABLE %s
	END
	
	SELECT * INTO %s FROM %s
	`, targetTable, targetTable, targetTable, targetData)

	fmt.Println("📢 log: Connect Backup Server ", destinationConn.Server, "-", destinationConn.Database)
	dbConn := ConnectDB(destinationConn)

	fmt.Println("🚀 Execute Query Backup Table", targetTable)
	ExecQuery(dbConn, queryBackupTable, false)

	fmt.Println("📢 log: ✅ Success Backup Table", targetTable)

	dbConn.Close()

	if isDropSourceTable {
		DropTable(backupSourceConn, targetTable)
	}
}

func DropTable(conn ConnectionParams, targetTable string) {
	targetTable = strings.TrimSpace(targetTable)

	if targetTable == "" {
		fmt.Println("⛔️ Target Table should not be empty ⛔️")
		return
	}

	fmt.Println("📢 log: Prepare Query Delete Table", targetTable)

	queryDeleteTable := fmt.Sprintf(`IF OBJECT_ID('%s') IS NOT NULL
	BEGIN
		DROP TABLE %s
	END`, targetTable, targetTable)

	fmt.Println("📢 log: Connect ", conn.Server, "-", conn.Database)
	dbConn := ConnectDB(conn)

	fmt.Println("🚀 Execute Query Delete Table", targetTable)
	ExecQuery(dbConn, queryDeleteTable, false)

	fmt.Println("📢 log: ✅ Success Delete Table", targetTable)
	dbConn.Close()
}

func GetListDatabaseInServer(conn ConnectionParams) []DatabaseList {
	fmt.Println("📢 log: Prepare Query Get List Database In Server", conn.Server)
	query := fmt.Sprintf(`SELECT [name] AS db_name
		FROM sys.sysdatabases
		WHERE [name] NOT IN ('master', 'tempdb', 'model', 'msdb', 'ReportServer', 'ReportServerTempDB')
		ORDER BY [name] ASC`)

	fmt.Println("📢 log: Connect ", conn.Server)
	conn.Database = "master"
	dbConn := ConnectDB(conn)

	fmt.Println("🚀 Execute Query get list Database In Server", conn.Server)
	rowData := ExecQuery(dbConn, query, true)

	var collectData []DatabaseList
	order := 0
	for rowData.Next() {
		var structDB DatabaseList
		err := rowData.StructScan(&structDB) // Scan each row into dynamic struct
		if err != nil {
			log.Println("⛔️ Row scan error:", err)
			break
		}
		order++
		structDB.Row = order
		collectData = append(collectData, structDB)
	}

	dbConn.Close()
	rowData.Close()

	fmt.Println(fmt.Sprintf("❇️ total Database in [%s] = %s", conn.Server, strconv.Itoa(len(collectData))))

	return collectData
}
