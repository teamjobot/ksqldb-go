This is a Go client for https://ksqldb.io/[ksqlDB]. It supports both pull and push queries, as well as command execution.

image::ksqldb-go.gif[Animation of the ksqlDB Golang client in action]

== Installation

Module install:

This client is a Go module, therefore you can have it simply by adding the following import to your code:

[source,golang]
----
import "github.com/teamjobot/ksqldb-go"
----

Then run a build to have this client automatically added to your go.mod file as a dependency.

Manual install:

[source,bash]
----
go get -u github.com/teamjobot/ksqldb-go
----

== Examples

See the link:test/environment.adoc[test environment here], and link:test/main.go[this sample code] which you can run with

[source,bash]
----
go run ./test/
----

Create a ksqlDB Client 

[source,go]
----
client := ksqldb.NewClient("http://ksqldb:8088","username","password").Debug()
----

For no authentication just use blank username and password values. 

=== Pull query

[source,go]
----
ctx, ctxCancel := context.WithTimeout(context.Background(), 10 * time.Second)
defer ctxCancel()

k := "SELECT TIMESTAMPTOSTRING(WINDOWSTART,'yyyy-MM-dd HH:mm:ss','Europe/London') AS WINDOW_START, TIMESTAMPTOSTRING(WINDOWEND,'HH:mm:ss','Europe/London') AS WINDOW_END, DOG_SIZE, DOGS_CT FROM DOGS_BY_SIZE WHERE DOG_SIZE='" + s + "';"
_, r, e := client.Pull(ctx, k, false)

if e != nil {
    // handle the error better here, e.g. check for no rows returned
    return fmt.Errorf("Error running Pull request against ksqlDB:\n%v", e)
}

var DOG_SIZE string
var DOGS_CT float64
for _, row := range r {
    if row != nil {
        // Should do some type assertions here
        DOG_SIZE = row[2].(string)
        DOGS_CT = row[3].(float64)
        fmt.Printf("🐶 There are %v dogs size %v\n", DOGS_CT, DOG_SIZE)
    }
}
----

=== Push query

[source,go]
----
rc := make(chan ksqldb.Row)
hc := make(chan ksqldb.Header, 1)

k := "SELECT ROWTIME, ID, NAME, DOGSIZE, AGE FROM DOGS EMIT CHANGES;"

// This Go routine will handle rows as and when they
// are sent to the channel
go func() {
    var NAME string
    var DOG_SIZE string
    for row := range rc {
        if row != nil {
            // Should do some type assertions here
            NAME = row[2].(string)
            DOG_SIZE = row[3].(string)

            fmt.Printf("🐾%v: %v\n",  NAME, DOG_SIZE)
        }
    }
}()

ctx, ctxCancel := context.WithTimeout(context.Background(), 10 * time.Second)
defer ctxCancel()

e := client.Push(ctx, k, rc, hc)

if e != nil {
    // handle the error better here, e.g. check for no rows returned
    return fmt.Errorf("Error running Push request against ksqlDB:\n%v", e)
}
----

=== Execute a command

[source,go]
----
if err := client.Execute(ctx, ksqlDBServer, `
	CREATE STREAM DOGS (ID STRING KEY, 
						NAME STRING, 
						DOGSIZE STRING, 
						AGE STRING) 
				  WITH (KAFKA_TOPIC='dogs', 
				  VALUE_FORMAT='JSON');
`); err != nil {
    return fmt.Errorf("Error creating the DOGS stream.\n%v", err)
}
----