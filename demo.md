# ID Generator Demo

This guide shows you how to quickly set up and run the `id-generator` demo.

## 1. Setup

Clone the repository and enter the project directory with this single command block:

```bash
git clone https://github.com/vaibhav423/id-generator.git && cd id-generator
```

## 2. Run Examples

Once the setup is complete, you can run the demo directly from your terminal. The program is located in the `demo` directory.

### Example 1: Default Run

Generate 100 IDs with the default machine ID (1) and save them to `generated_ids.txt`.

```bash
go run demo/main.go
```

### Example 2: Generate a Specific Number of IDs

Generate 1,000 IDs and save them to the default file.

```bash
go run demo/main.go -n 1000
```

### Example 3: Specify Machine ID and Output File

Generate 500 IDs using machine ID `25` and save them to a custom file named `machine_25_ids.txt`.

```bash
go run demo/main.go -m 25 -n 500 -o machine_25_ids.txt
```