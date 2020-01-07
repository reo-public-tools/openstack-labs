package main

import (
    "os"
    "fmt"
    "log"
    "flag"
    "log/syslog"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)


// Globals and constants
const globalConfigFile = "/etc/oslabs/oslabs.yml"
var sysLog *syslog.Writer


// Main func
func main() {

    // ############################
    // #### Initialize logging ####
    // ############################

    // Set up the syslog logger with some defaults
    var err error
    sysLog, err = syslog.New(syslog.LOG_EMERG | syslog.LOG_USER, "oslabs")
    if err != nil {
        log.Fatal(err)
    }

    // Override the writers in the related packages
    // to change the identity id.  Later on we will
    // want to change where its writing for things
    // like remote logging.
    theforeman.OverrideLogWriter(sysLog)
    osutils.OverrideLogWriter(sysLog)



    // ##########################
    // #### Argument Parsing ####
    // ##########################

    // Set up subcommands
    createCommand := flag.NewFlagSet("create", flag.ExitOnError)
    listCommand   := flag.NewFlagSet("list", flag.ExitOnError)
    deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
    showCommand   := flag.NewFlagSet("show", flag.ExitOnError)
    scratchCommand := flag.NewFlagSet("scratch", flag.ExitOnError)
    //checkCommand  := flag.NewFlagSet("check", flag.ExitOnError)

    // Create subcommand flag pointers
    configFilePtr := createCommand.String("c", "", "Lab yaml config file")

    // Delete subcommand flag pointers
    deleteLabNamePtr := deleteCommand.String("l", "", "Lab name to clean up")

    // Show subcommand flag pointers
    showLabNamePtr := showCommand.String("l", "", "Lab name pull details on")

    // Create subcommand flag pointers
    scratchConfigFilePtr := scratchCommand.String("c", "", "Lab yaml config file")


    // Verify that the subcommand is provided
    // os.Arg[0] is the main command
    // os.Arg[1] is the subcommand
    if len(os.Args) < 2 {
        fmt.Println("create, check, delete or list subcommand is required")
        os.Exit(1)
    }

    // Switch subcommand
    switch os.Args[1] {
    case "create":
        createCommand.Parse(os.Args[2:])
    case "list":
        listCommand.Parse(os.Args[2:])
    case "delete":
        deleteCommand.Parse(os.Args[2:])
    case "show":
        showCommand.Parse(os.Args[2:])
    case "scratch":
        scratchCommand.Parse(os.Args[2:])
    default:
        flag.PrintDefaults()
        os.Exit(1)
    }


    // Run create if parsed
    if createCommand.Parsed() {

        // Exist out if arg is empty.
        if *configFilePtr == "" {
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Run the create function on the config
        err := Create(*configFilePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Run list if parsed
    if listCommand.Parsed() {
        err := List()
        if err != nil {
            log.Fatal(err)
        }
    }

    // Run delete if parsed
    if deleteCommand.Parsed() {

        // Exist out if arg is empty.
        if *deleteLabNamePtr == "" {
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Run function to clean up and release|delete a lab
        err := Delete(*deleteLabNamePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Run show if parsed
    if showCommand.Parsed() {

        // Exist out if arg is empty.
        if *showLabNamePtr == "" {
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Display details for the specified lab
        err := Show(*showLabNamePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Scratch action used for tested & development
    if scratchCommand.Parsed() {

        // Exist out if arg is empty.
        if *scratchConfigFilePtr == "" {
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Run the create function on the config
        err := Scratch(*scratchConfigFilePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

}
