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

package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/apache/calcite-avatica-go/v5"
	prompt "github.com/apache/calcite-avatica-go/cmd/cli/prompt"
	"github.com/spf13/cobra"
)

var (
	connectionURL    = "http://localhost:8080"
	serialization    = "protobuf"
	schema           string
	connectionParams string
	user             string
	passwd           string
	maxRowsTotal     string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cli",
		Short: "A calcite CLI prompt to execute queries",
		Run:   runSQLPrompt,
	}

	// Define flags for connection URL and additional parameters
	rootCmd.Flags().StringVar(&connectionURL, "url", connectionURL, "Connection URL")
	rootCmd.Flags().StringVar(&serialization, "serialization", "", "Serialization parameter")
	rootCmd.Flags().StringVar(&connectionParams, "params", "", "Extra parameters for avatica connection (ex: \"parameter1=value&...parameterN=value\")")
	rootCmd.Flags().StringVarP(&schema, "schema", "s", "", "The schema path sets the default schema to use for this connection.")
	rootCmd.Flags().StringVarP(&user, "username", "u", "", "The user to use when authenticating against Avatica")
	rootCmd.Flags().StringVarP(&passwd, "password", "p", "", "The password to use when authenticating against Avatica")
	rootCmd.MarkFlagsRequiredTogether("username", "password")
	rootCmd.Flags().StringVarP(&maxRowsTotal, "maxRowsTotal", "m", "", "The maxRowsTotal parameter sets the maximum number of rows to return for a given query")

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func runSQLPrompt(cmd *cobra.Command, args []string) {
	// Establish a connection to the calcite server
	db := establishConnection()
	defer db.Close()

	// Create and run the SQL prompt
	prompt.CreateAndRunPrompt(db)
}

func establishConnection() *sql.DB {
	parameters := buildConnectionURL()
	fmt.Println("Connecting to ", parameters)
	db, err := sql.Open("avatica", parameters)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected")
	return db
}

func buildConnectionURL() string {
	var url strings.Builder

	// Append the connection URL
	url.WriteString(connectionURL)

	var params []string

	// Add serialization parameter by default protobuf
	if serialization != "" {
		params = append(params, "serialization="+serialization)
	}

	// Add username and password as parameter
	if user != "" {
		params = append(params, "avaticaUser="+user)
		params = append(params, "avaticaPassword="+passwd)
	}

	// Add connection parameters
	if connectionParams != "" {
		params = append(params, connectionParams)
	}

	if maxRowsTotal != "" {
		params = append(params, "maxRowsTotal="+maxRowsTotal)
	}

	// Combine the connection URL and parameters
	if schema != "" {
		url.WriteString("/")
		url.WriteString(schema)
	}

	if len(params) > 0 {
		url.WriteString("?")
		url.WriteString(strings.Join(params, "&"))
	}

	return url.String()
}
