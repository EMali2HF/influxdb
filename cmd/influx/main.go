package main

import (
	"fmt"
	"os"

	"github.com/influxdata/influxdb/client"
	"github.com/influxdata/influxdb/cmd/influx/cli"
	flag "github.com/spf13/pflag"
)

// These variables are populated via the Go linker.
var (
	version string
)

const (
	// defaultFormat is the default format of the results when issuing queries
	defaultFormat = "column"

	// defaultPrecision is the default timestamp format of the results when issuing queries
	defaultPrecision = "ns"

	// defaultPPS is the default points per second that the import will throttle at
	// by default it's 0, which means it will not throttle
	defaultPPS = 0
)

func init() {
	// If version is not set, make that clear.
	if version == "" {
		version = "unknown"
	}
}

func main() {
	c := cli.New(version)

	fs := flag.NewFlagSet("InfluxDB shell version "+version, flag.ExitOnError)
	fs.StringVarP(&c.Host, "host", "H", client.DefaultHost, "Influxdb host to connect to.")
	fs.IntVarP(&c.Port, "port", "p", client.DefaultPort, "Influxdb port to connect to.")
	fs.StringVarP(&c.Username, "username", "u", c.Username, "Username to connect to the server.")
	fs.StringVar(&c.Password, "password", c.Password, `Password to connect to the server.  Leaving blank will prompt for password (--password="").`)
	fs.StringVarP(&c.Database, "database", "d", c.Database, "Database to connect to the server.")
	fs.BoolVar(&c.Ssl, "ssl", false, "Use https for connecting to cluster.")
	fs.BoolVar(&c.UnsafeSsl, "unsafeSsl", false, "Set this when connecting to the cluster using https and not use SSL verification.")
	fs.StringVar(&c.Format, "format", defaultFormat, "Format specifies the format of the server responses:  json, csv, or column.")
	fs.StringVar(&c.Precision, "precision", defaultPrecision, "Precision specifies the format of the timestamp:  rfc3339,h,m,s,ms,u or ns.")
	fs.StringVar(&c.WriteConsistency, "consistency", "all", "Set write consistency level: any, one, quorum, or all.")
	fs.BoolVar(&c.Pretty, "pretty", false, "Turns on pretty print for the json format.")
	fs.StringVarP(&c.Execute, "execute", "e", c.Execute, "Execute command and quit.")
	fs.BoolVarP(&c.ShowVersion, "version", "v", false, "Displays the InfluxDB version.")
	fs.BoolVar(&c.Import, "import", false, "Import a previous database.")
	fs.IntVar(&c.PPS, "pps", defaultPPS, "How many points per second the import will allow.  By default it is zero and will not throttle importing.")
	fs.StringVar(&c.Path, "path", "", "path to the file to import")
	fs.BoolVar(&c.Compressed, "compressed", false, "set to true if the import file is compressed")

	// Define our own custom usage to print
	fs.Usage = func() {
		fmt.Println(`Usage of influx:
  -v, --version
       Display the version and exit.
  -H, --host 'host name'
       Host to connect to.
  -p, --port 'port #'
       Port to connect to.
  -d, --database 'database name'
       Database to connect to the server.
  --password 'password'
      Password to connect to the server.  Leaving blank will prompt for password (--password '').
  -u, --username 'username'
       Username to connect to the server.
  --ssl
        Use https for requests.
  --unsafeSsl
        Set this when connecting to the cluster using https and not use SSL verification.
  -e, --execute 'command'
       Execute command and quit.
  --format 'json|csv|column'
       Format specifies the format of the server responses:  json, csv, or column.
  --precision 'rfc3339|h|m|s|ms|u|ns'
       Precision specifies the format of the timestamp:  rfc3339, h, m, s, ms, u or ns.
  --consistency 'any|one|quorum|all'
       Set write consistency level: any, one, quorum, or all
  --pretty
       Turns on pretty print for the json format.
  --import
       Import a previous database export from file
  --pps
       How many points per second the import will allow.  By default it is zero and will not throttle importing.
  --path
       Path to file to import
  --compressed
       Set to true if the import file is compressed

Examples:

    # Use influx in a non-interactive mode to query the database "metrics" and pretty print json:
    $ influx --database 'metrics' --execute 'select * from cpu' --format 'json' --pretty

    # Connect to a specific database on startup and set database context:
    $ influx --database 'metrics' --host 'localhost' --port '8086'
`)
	}
	fs.Parse(os.Args[1:])

	if c.ShowVersion {
		c.Version()
		os.Exit(0)
	}

	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
