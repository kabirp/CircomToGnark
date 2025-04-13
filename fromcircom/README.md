# Circom + gnark Bridge (BN254 Groth16)

> **Disclaimer**  
> This is a prototype and **not intended for production**.  
> The code has **not been audited**.  
> Suggestions and contributions are welcome.

---

## Motivation

- Circom is great for defining R1CS circuits—especially with its robust library ecosystem.
- The gnark prover is significantly faster than Circom's.
- Goal: Combine Circom's ergonomics with gnark's performance.

**Note:** This library currently supports only **Groth16 over BN254**.

---

## Solution Overview

### Setup

- Define your R1CS circuit in Circom and compile it into a `.r1cs` binary.
- Use this tool to:
  - Convert the `.r1cs` file to a gnark-compatible format.
  - Run `groth16.Setup()` to generate the prover and verifier keys.

### Compute and Prove

- Generate your witness using Circom/snarkjs to produce `witness.wtns`.
- Use this tool to:
  - Convert the witness into a format gnark understands.
  - Generate a proof using gnark.

### Verify

- Use this tool to verify the proof with gnark’s verifier.

---

## Example Commands

These should be run from specific working directories.

### Directory: `fromcircom/circom/simple_circom_circuit`

```bash
alias circom_simple="circom --r1cs --wasm --sym simple.circom"
alias circom_simple_witgen="node simple_js/generate_witness.js simple_js/simple.wasm input.json witness.wtns"
```

### Directory: `fromcircom/`

```bash
alias do_define="go run main.go parse_circom_helper.go define_circuit circom/simple_circom_circuit/simple.r1cs gnarkCirc2.txt"
alias do_setup="go run main.go parse_circom_helper.go setup_circuit gnarkCirc2.txt provKey2.txt verKey2.txt"
alias do_witness="go run main.go parse_circom_helper.go import_witness circom/simple_circom_circuit/witness.wtns gnarkCirc2.txt gnarkWit2.txt pubInps2.txt"
alias do_prove="go run main.go parse_circom_helper.go prove gnarkCirc2.txt provKey2.txt gnarkWit2.txt proof2.txt"
alias do_verify="go run main.go parse_circom_helper.go verify verKey2.txt proof2.txt pubInps2.txt"
```

### Example `input.json`

Place this file in: `fromcircom/circom/simple_circom_circuit`

```json
{
  "x": "3",
  "y": "11"
}
```

---

## Running an Example

First, you will have to update line 3 of the top level directory's go.mod file appropriate. 

### Setup

```bash
cd fromcircom/circom/simple_circom_circuit
circom_simple

cd ../../../fromcircom
do_define
do_setup
```

### Compute and Prove

```bash
cd circom/simple_circom_circuit
# Make sure input.json is present
circom_simple_witgen

cd ../../
do_witness
do_prove
```

### Verify

```bash
do_verify
```
