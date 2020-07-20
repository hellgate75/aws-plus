package commands

import (
	"flag"
	"fmt"
	"github.com/hellgate75/aws-plus/connect"
	"os"
	"strconv"
)
var commandsInitialized = false
func InitCommands() {
	if commandsInitialized {
		return
	}
	initHelpCommand()
	initListCommand()
	commandsInitialized = true
}

func printHelp(command *Command, fg *flag.FlagSet) {
		fmt.Println("aws-plus", command.String())
		fmt.Println("use aws-plus help <command> for details")
		fmt.Println("Options:")
		fg.PrintDefaults()
		fmt.Println()
		fmt.Println()
}

func initHelpCommand() {
	commands["help"] = Command{
		Name: "help",
		Sub: []string{"command"},
		Args: []string{},
		Action: func(c *Command, args ...Arg) Response {
			resp := Response{Content: NilContentText}
			if cd, ok := c.data["code"]; ok && cd != CodeOk {
				// Issues on parsing pre-conditions
				code, _ := strconv.Atoi(fmt.Sprintf("%v", cd))
				resp.Code = code
				resp.Error = "Error occurred during parsing"
			} else if err, ok := c.data["error"]; ok && err != nil {
				// Issues on parsing execution
				resp.Code = int(CodeParserError)
				resp.Error = fmt.Sprintf("%v", err)
			} else {
				subCommand := fmt.Sprintf("%v", c.data["command"])
				if subCommand == "" || subCommand == "nil" {
					resp.Code = int(CodeExecFatal)
					resp.Error = "Empty help sub-command"
				} else {
					hC, err := CommandByName(subCommand)
					if err != nil {
						resp.Code = int(CodeExecFatal)
						resp.Error = fmt.Sprintf("Unable to find help for command: %s", subCommand)
					} else if v, ok := c.data["flagSet"]; ok && v == nil {
						resp.Code = int(CodeExecFatal)
						resp.Error = "Unable to collect flag-set reference"
					} else {
						fg := c.data["flagSet"].(*flag.FlagSet)
						printHelp(hC, fg)
						os.Exit(0)
					}
				}
			}
			return resp
		},
		Parser: func(c *Command, args []string, fg *flag.FlagSet) {
			c.data["flagSet"]=fg
			c.data["error"]=nil
			if len(args) < 3 {
				c.data["code"]=CodeInsufficientArgs
			} else {
				c.data["command"]=args[2]
				c.data["code"]=CodeOk
				if len(args) >= 3 {
					c.data["error"]=fg.Parse(args[3:])
				} else if len(c.Args) > 0 {
					c.data["code"]=CodeInsufficientArgs
				}
			}
		},
		data: make(map[string]interface{}),
	}
}


func initListCommand() {
	commands["list"] = Command{
		Name: "list",
		Sub: []string{"subject"},
		Args: []string{},
		Action: func(c *Command, args ...Arg) Response {
			resp := Response{Content: NilContentText}
			if cd, ok := c.data["code"]; ok && cd != CodeOk {
				// Issues on parsing pre-conditions
				code, _ := strconv.Atoi(fmt.Sprintf("%v", cd))
				resp.Code = code
				resp.Error = "Error occurred during parsing"
			} else if err, ok := c.data["error"]; ok && err != nil {
				// Issues on parsing execution
				resp.Code = int(CodeParserError)
				resp.Error = fmt.Sprintf("%v", err)
			} else {
				subCommand := fmt.Sprintf("%v", c.data["subject"])
				if subCommand == "" || subCommand == "nil" {
					resp.Code = int(CodeExecFatal)
					resp.Error = "Empty help sub-command"
				} else {
					var partition = ""
					var silent = false
					for _, arg := range args {
						if arg.Name == "partition" {
							partition = arg.String()
						} else if arg.Name == "silent" {
							silent = arg.Bool()
						}
					}
					if ! silent {
						fmt.Println("Subject:", subCommand)
					}
					resp.Code = int(CodeOk)
					resp.Error = "-- --"
					if subCommand == "partitions" {
						resp.Content = connect.GetPartitionIds()
					} else if subCommand == "regions" {
						part, err := connect.GetAwsPartition(partition)
						if err != nil {
							resp.Code = int(CodeExecError)
							resp.Error = fmt.Sprintf("Cannot access partition: %s", partition)
						} else {
							resp.Content = part.GetRegionNames()
						}
					} else if subCommand == "services" {
						var partition = ""
						for _, arg := range args {
							if arg.Name == "partition" {
								partition = fmt.Sprintf("%v", arg.Value)
							}
						}
						part, err := connect.GetAwsPartition(partition)
						if err != nil {
							resp.Code = int(CodeExecError)
							resp.Error = fmt.Sprintf("Cannot access partition: %s", partition)
						} else {
							resp.Content = part.ServiceNames()
						}
					} else {
						resp.Code = int(CodeExecError)
						resp.Error = fmt.Sprintf("Unknown list type: %s, available: %s", subCommand, "partitions, regions, services")
					}
				}
			}
			return resp
		},
		Parser: func(c *Command, args []string, fg *flag.FlagSet) {
			c.data["flagSet"]=fg
			c.data["error"]=nil
			if len(args) < 3 {
				c.data["code"]=CodeInsufficientArgs
			} else {
				c.data["subject"]=args[2]
				c.data["code"]=CodeOk
				if len(args) >= 3 {
					c.data["error"]=fg.Parse(args[3:])
				} else if len(c.Args) > 0 {
					c.data["code"]=CodeInsufficientArgs
				}
			}
		},
		data: make(map[string]interface{}),
	}
}


