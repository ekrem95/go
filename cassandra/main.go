package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var (
	keyspace = "keyspace_name"
)

func main() {
	// Provide the cassandra cluster instance here.
	cluster := gocql.NewCluster("127.0.0.1")

	// The authenticator is needed if password authentication is
	// enabled for your Cassandra installation..
	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: "username",
	// 	Password: "password",
	// }

	// gocql requires the keyspace to be provided before the session is created.
	// In future there might be provisions to do this later.
	cluster.Keyspace = "system"

	// This is time after which the creation of session call would timeout.
	// This can be customised as needed.
	cluster.Timeout = 5 * time.Second

	cluster.ProtoVersion = 4
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to cassandra cluster: %v", err)
	}

	if err = session.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
		WITH replication = {'class' : 'SimpleStrategy',	'replication_factor' : %d}`,
		keyspace, 1)).Exec(); err != nil {
		log.Fatal(err)
	}
	session.Close()

	cluster.Keyspace = keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to cassandra cluster: %v", err)
	}
	defer session.Close()

	// Check if the table already exists. Create if table does not exist
	keySpaceMeta, _ := session.KeyspaceMetadata(keyspace)

	if _, exists := keySpaceMeta.Tables["person"]; !exists {
		session.Query("CREATE TABLE person (" +
			"id text, name text, phone text, " +
			"PRIMARY KEY (id))").Exec()
	}

	// Insert records into table using prepared statements
	session.Query("INSERT INTO person (id, name, phone) VALUES (?, ?, ?)",
		"ekrem", "Ekrem", "1234567890").Exec()
	session.Query("INSERT INTO person (id, name, phone) VALUES (?, ?, ?)",
		"eko", "Karatas", "1234567890").Exec()

	// Update a record
	session.Query("UPDATE person SET phone = ? WHERE id = ?", "123456", "eko").Exec()

	// Select record and run some process on data fetched
	var name, phone string
	if err := session.Query(
		"SELECT name, phone FROM person WHERE id='ekrem'").Scan(
		&name, &phone); err != nil {
		if err != gocql.ErrNotFound {
			log.Fatalf("Query failed: %v", err)
		}
	}
	log.Printf("Name: %v", name)
	log.Printf("Phone: %v", phone)

	// Fetch multiple rows and run process over them
	iter := session.Query("SELECT name, phone FROM person").Iter()
	for iter.Scan(&name, &phone) {
		log.Printf("Iter Name: %v", name)
		log.Printf("Iter Phone: %v", phone)
	}

	// Delete a record
	if err = session.Query("DELETE FROM person WHERE id = ?", "eko").Exec(); err != nil {
		log.Fatal(err)
	}
}
