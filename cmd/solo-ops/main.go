package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/nstranquist/solo-ops/internal/soloops"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: solo-ops inquiry-add|inquiry-list|policy-pack ...")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "inquiry-add":
		// --ledger --id --name --event --guests --notes
		var ledger, id, name, event, notes string
		guests := 0
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--ledger":
				i++
				ledger = os.Args[i]
			case "--id":
				i++
				id = os.Args[i]
			case "--name":
				i++
				name = os.Args[i]
			case "--event":
				i++
				event = os.Args[i]
			case "--guests":
				i++
				guests, _ = strconv.Atoi(os.Args[i])
			case "--notes":
				i++
				notes = os.Args[i]
			}
		}
		if ledger == "" || name == "" {
			fmt.Fprintln(os.Stderr, "need --ledger and --name")
			os.Exit(2)
		}
		inq := soloops.NewInquiry(id, name, event, guests, notes)
		if err := soloops.AppendInquiry(ledger, inq); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(inq)
	case "inquiry-list":
		ledger := "testdata/inquiries.jsonl"
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--ledger" {
				i++
				ledger = os.Args[i]
			}
		}
		list, err := soloops.ListInquiries(ledger)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(list)
	case "policy-pack":
		biz := "Demo Business"
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--business" {
				i++
				biz = os.Args[i]
			}
		}
		pack := soloops.GeneratePolicyPack(biz)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(pack)
	default:
		fmt.Fprintf(os.Stderr, "unknown %q\n", os.Args[1])
		os.Exit(2)
	}
}
