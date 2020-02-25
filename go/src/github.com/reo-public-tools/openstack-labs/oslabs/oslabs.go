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
    capCommand    := flag.NewFlagSet("capacity", flag.ExitOnError)
    scratchCommand := flag.NewFlagSet("scratch", flag.ExitOnError)
    geninvCommand := flag.NewFlagSet("geninv", flag.ExitOnError)
    //checkCommand  := flag.NewFlagSet("check", flag.ExitOnError)

    // Create subcommand flag pointers
    configFilePtr := createCommand.String("c", "", "Lab yaml config file")

    // Delete subcommand flag pointers
    deleteLabNamePtr := deleteCommand.String("l", "", "Lab name to clean up")

    // Show subcommand flag pointers
    showLabNamePtr := showCommand.String("l", "", "Lab name pull details on")

    // Create subcommand flag pointers
    scratchConfigFilePtr := scratchCommand.String("c", "", "Lab yaml config file")

    // Geninventory subcommand flag pointers
    geninvLabNamePtr := geninvCommand.String("l", "", "Lab name pull inventory for")


    // Verify that the subcommand is provided
    // os.Arg[0] is the main command
    // os.Arg[1] is the subcommand
    if len(os.Args) < 2 {
        displaySyntax(os.Args[0])
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
    case "geninv":
        geninvCommand.Parse(os.Args[2:])
    case "capacity":
        capCommand.Parse(os.Args[2:])
    case "scratch":
        scratchCommand.Parse(os.Args[2:])
    case "help":
        displaySyntax(os.Args[0])
        os.Exit(1)
    default:
        displaySyntax(os.Args[0])
        flag.PrintDefaults()
        os.Exit(1)
    }


    // Run create if parsed
    if createCommand.Parsed() {

        // Exist out if arg is empty.
        if *configFilePtr == "" {
            displaySyntax(os.Args[0])
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
            displaySyntax(os.Args[0])
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
            displaySyntax(os.Args[0])
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Display details for the specified lab
        err := Show(*showLabNamePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Run geninv if parsed
    if geninvCommand.Parsed() {

        // Exist out if arg is empty.
        if *geninvLabNamePtr == "" {
            displaySyntax(os.Args[0])
            flag.PrintDefaults()
            os.Exit(1)
        }

        // Generate ansible inventory for the specified lab
        err := GenerateInventory(*geninvLabNamePtr)
        if err != nil {
            log.Fatal(err)
        }
    }

    // Run list if parsed
    if capCommand.Parsed() {
        err := Capacity()
        if err != nil {
            log.Fatal(err)
        }
    }

    // Scratch action used for tested & development
    if scratchCommand.Parsed() {

        // Exist out if arg is empty.
        if *scratchConfigFilePtr == "" {
            displaySyntax(os.Args[0])
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

func displaySyntax(progName string) {
    fmt.Printf("\nSyntax:\n\n")
    fmt.Printf("    Create a new lab:\n")
    fmt.Printf("        %s create -c ./path/to/lab_config.yaml\n\n", progName)
    fmt.Printf("    List labs:\n")
    fmt.Printf("        %s list\n\n", progName)
    fmt.Printf("    Show details on existing lab:\n")
    fmt.Printf("        %s show -l labname\n\n", progName)
    fmt.Printf("    Generate ansible inventory for an existing lab:\n")
    fmt.Printf("        %s geninv -l labname\n\n", progName)
    fmt.Printf("    Delete or relase an existing lab:\n")
    fmt.Printf("        %s delete -l labname\n\n", progName)
    fmt.Printf("    Get Ironic Capacity:\n")
    fmt.Printf("        %s capacity\n\n", progName)
    fmt.Printf("    Display command syntax:\n")
    fmt.Printf("        %s help\n\n", progName)
}
