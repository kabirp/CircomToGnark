package main

import (
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// Supported commands
const (
	DefineCircuitCommand = "define_circuit"
	SetupCircuitCommand  = "setup_circuit"
	ImportWitnessCommand = "import_witness"
	ProveCommand         = "prove"
	VerifyCommand        = "verify"
)
const SupportedCommands = "define_circuit,setup_circuit,import_witness,prove,verify"

// Supported usage
const (
	DefineCircuitUsage = "define_circuit <circomCircuitInPath> <gnarkCircuitOutPath>"
	SetupCircuitUsage  = "setup_circuit <gnarkCircuitInPath> <pkOutPath> <vkOutPath>"
	ImportWitnessUsage = "import_witness <circomWitnessInPath> <gnarkCircuitInPath> <wtnsOutPath> <pubInputsOutPath>"
	ProveUsage         = "prove <pkInPath> <wtnsInPath> <proofOutPath>"
	VerifyUsage        = "verify <vkInPath> <proofInPath> <pubInputsInPath>"
)

func define_circuit(circomCircuitInPath string, gnarkCircuitOutPath string) {
	// Log function arguments
	log.Println("Defining circuit. circomCircuitInPath=" + circomCircuitInPath +
		" gnarkCircuitOutPath=" + gnarkCircuitOutPath)

	// Open the input file
	sectionTypeToMetadata := isFileValid(circomCircuitInPath, R1CSCircuitBinary)

	// Parse the R1CS header section
	// FIXME: (no need to open file multiple times)
	circuit := parseR1CSHeaderSection(circomCircuitInPath, sectionTypeToMetadata)

	// Compile the circuit
	constraintSystem, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit, frontend.IgnoreUnconstrainedInputs())
	if err != nil {
		log.Fatalf("Failed to compile circuit: %v", err)
	}

	_ = constraintSystem
	panic("define_circuit Not fully implemented yet")
}

func main() {
	command := os.Args[1]

	switch command {
	case DefineCircuitCommand:
		if len(os.Args) != 4 {
			fmt.Println("Incorrect Usage. Expected Usage: " + DefineCircuitUsage)
			os.Exit(1)
		}
		define_circuit(os.Args[2], os.Args[3])
	case SetupCircuitCommand:
		// TODO: Implement setup_circuit
		if len(os.Args) != 5 {
			fmt.Println("Incorrect Usage. Expected Usage: " + SetupCircuitUsage)
			os.Exit(1)
		}
		// setup_circuit(os.Args[2], os.Args[3], os.Args[4])
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Supported commands:" + SupportedCommands)
		os.Exit(1)
	}
}
