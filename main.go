package main

import (
	"flag"
	"fmt"
	"github.com/hellgate75/aws-plus/commands"
	"github.com/hellgate75/aws-plus/connect"
	"github.com/hellgate75/aws-plus/io"
	"os"
	"strings"
)

var partition string
var region string
var output string
var silent bool

func init() {
	connect.InitPartitions()
	commands.InitCommands()
	preConditions()
}


func printHelpUsage(cmd ...string) {
	printUsage()
}

func printUsage() {
	names := strings.Join(commands.CommandNames(), ", ")
	fmt.Println("aws-plus [command] -opt0=value0 -opt2=value1 ....  -optN=valueN")
	fmt.Println("use aws-plus help <command> for details")
	fmt.Println("Commands:", names)
	fmt.Println("Options:")
	initParameters().PrintDefaults()
	fmt.Println()
	fmt.Println()
}

func preConditions() {
	if len(os.Args) < 2 {
		fmt.Println("Insufficient parameters")
		printUsage()
		os.Exit(1)
	}
}
var cFlag *flag.FlagSet

func initParameters() *flag.FlagSet {
	if cFlag != nil {
		return cFlag
	}
	cFlag = flag.NewFlagSet("aws-plus", flag.ContinueOnError)
	cFlag.StringVar(&partition, "partition", "aws", fmt.Sprintf("Working partitions (available: %s)", partitionDescriptors()))
	cFlag.StringVar(&region, "region", "eu-west-1", fmt.Sprintf("Partition selected region (available: %s)", partitionRegionsDescriptors("aws")))
	cFlag.StringVar(&output, "output", "text", "Output format (text, json, yaml, xml)")
	cFlag.BoolVar(&silent, "silent", false, "Produce silent output")
	return cFlag
}

func partitionDescriptors() string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", connect.GetPartitionIds()), " ", ", "), "]", ""), "[", "")
}


func partitionRegionsDescriptors(region string) string {
	p, _ := connect.GetAwsPartition(region)
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", p.GetRegionNames()), " ", ", "), "]", ""), "[", "")
}

func main() {
	cmd := strings.ToLower(os.Args[1])
	if ! commands.ValidateCommand(cmd) {
		fmt.Printf("Unknown command: %s", cmd)
		printUsage()
		os.Exit(2)
	}
	command, err := commands.CommandByName(cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		command.Parser(command, os.Args, initParameters())
		var args = make([]commands.Arg, 0)
		args = append(args, commands.Arg {
			Name: "partition",
			Value: partition,
		})
		args = append(args, commands.Arg {
			Name: "region",
			Value: region,
		})
		args = append(args, commands.Arg {
			Name: "output",
			Value: output,
		})
		args = append(args, commands.Arg {
			Name: "silent",
			Value: silent,
		})
		if errC := command.Accepts(args...); errC != commands.CodeOk {
			fmt.Printf("Invalid parameters code: %v\n", errC)
			printHelpUsage(cmd)
		}
		response := command.Action(command, args...)
		enc := io.ParseEncoding(output)
		//fmt.Printf("Output: %+v\n", response)
		//fmt.Printf("Using encoding format: %v\n", enc)
		response.Code += 200
		data, err := io.EncodeValue(&response, enc)
		if err != nil {
			fmt.Printf("Error encoding response: %v\n", err)
			printHelpUsage(cmd)
			os.Exit(2)
		}
		if ! silent {
			fmt.Println("Command:", command.Name)
		}
		if response.Code != 0 {
			if silent {
				fmt.Println(string(data))
			} else {
				fmt.Printf("Error, response:\n%s\n", string(data))
				printHelpUsage(cmd)
			}
		} else {
			if silent {
				fmt.Println(string(data))
			} else {
				fmt.Printf("Success, response:\n%s\n", string(data))
			}
		}
		os.Exit(0)
	}
}
