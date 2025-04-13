package main

import (
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
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
	gnarkConstraintSystem, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit, frontend.IgnoreUnconstrainedInputs())
	if err != nil {
		log.Fatalf("Failed to compile circuit: %v", err)
	}

	outFile, err := os.Create(gnarkCircuitOutPath)
	if err != nil {
		log.Fatalf("Failed to create gnark circuit output file: %v", err)
	}
	_, err = gnarkConstraintSystem.WriteTo(outFile)
	if err != nil {
		log.Fatalf("Failed to write gnark circuit to output file: %v", err)
	}
	outFile.Close()
	log.Println("Circuit defined successfully.")
}

func setup_circuit(gnarkCircuitInPath string, pkOutPath string, vkOutPath string) {
	inFile, err := os.Open(gnarkCircuitInPath)
	if err != nil {
		log.Fatalf("Failed to open gnark circuit input file: %v", err)
	}
	var gnarkConstraintSystem constraint.ConstraintSystem = groth16.NewCS(ecc.BN254)
	gnarkConstraintSystem.ReadFrom(inFile)
	inFile.Close()

	prover_key, verifier_key, err := groth16.Setup(gnarkConstraintSystem)
	if err != nil {
		log.Fatalf("Failed to setup circuit: %v", err)
	}
	pkFile, err := os.Create(pkOutPath)
	if err != nil {
		log.Fatalf("Failed to create pk output file: %v", err)
	}
	_, err = prover_key.WriteTo(pkFile)
	if err != nil {
		log.Fatalf("Failed to write pk to output file: %v", err)
	}
	pkFile.Close()
	vkFile, err := os.Create(vkOutPath)
	if err != nil {
		log.Fatalf("Failed to create vk output file: %v", err)
	}
	_, err = verifier_key.WriteTo(vkFile)
	if err != nil {
		log.Fatalf("Failed to write vk to output file: %v", err)
	}
	vkFile.Close()
	log.Println("Setup completed successfully.")
}

func import_witness(circomWitnessInPath string, gnarkCircuitInPath string, wtnsOutPath string, pubInputsOutPath string) {
	// Log function arguments
	log.Println("Importing witness. circomWitnessInPath=" + circomWitnessInPath +
		" gnarkCircuitInPath=" + gnarkCircuitInPath +
		" wtnsOutPath=" + wtnsOutPath +
		" pubInputsOutPath=" + pubInputsOutPath)

	panic("Import witness not implemented yet")
	// log.Println("Witness imported successfully.")
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
		if len(os.Args) != 5 {
			fmt.Println("Incorrect Usage. Expected Usage: " + SetupCircuitUsage)
			os.Exit(1)
		}
		setup_circuit(os.Args[2], os.Args[3], os.Args[4])
	case ImportWitnessCommand:
		if len(os.Args) != 6 {
			fmt.Println("Incorrect Usage. Expected Usage: " + ImportWitnessUsage)
			os.Exit(1)
		}
		import_witness(os.Args[2], os.Args[3], os.Args[4], os.Args[5])
	default:
		fmt.Println("Unknown command:", command)
		fmt.Println("Supported commands:" + SupportedCommands)
		os.Exit(1)
	}
}
