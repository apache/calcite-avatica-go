/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calcitesql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/apache/calcite-avatica-go/v5"
	"github.com/olekukonko/tablewriter"
)

func ExecuteQuery(db *sql.DB, query string) {
	cmd := strings.TrimRight(query, ";")
	start := time.Now()
	// Execute the query
	rows, err := db.Query(cmd)
	duration := time.Since(start)
	if err != nil {
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Println("Error retrieving column names:", err)
		return
	}

	// Create a new table writer for each query
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(true)

	// Create a slice to store the query results
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch and print rows
	count := 0
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println("Error retrieving row data:", err)
			continue
		}

		// Prepare row data
		rowData := make([]string, len(columns))
		for i, v := range values {
			if v != nil {
				rowData[i] = fmt.Sprintf("%v", v)
			} else {
				rowData[i] = "NULL"
			}
		}

		// Add row to the table
		table.Append(rowData)
		count++
	}

	// Set the table headers
	table.SetHeader(columns)

	// Render the table
	table.Render()

	fmt.Printf("Rows: %d\nExecution Time: %s\n\n", count, duration)
}