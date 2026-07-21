package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/nstranquist/solo-ops/internal/soloops"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: solo-ops inquiry-add|inquiry-list|inquiry-set|policy-pack|serve|cockpit-html ...")
		os.Exit(2)
	}
	switch os.Args[1] {
	case "inquiry-add":
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
		encode(inq)
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
		encode(list)
	case "inquiry-set":
		var ledger, id, status string
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--ledger":
				i++
				ledger = os.Args[i]
			case "--id":
				i++
				id = os.Args[i]
			case "--status":
				i++
				status = os.Args[i]
			}
		}
		if ledger == "" || id == "" || status == "" {
			fmt.Fprintln(os.Stderr, "usage: solo-ops inquiry-set --ledger f --id id --status s")
			os.Exit(2)
		}
		inq, err := soloops.SetStatus(ledger, id, status)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		encode(inq)
	case "policy-pack":
		biz := "Demo Business"
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--business" {
				i++
				biz = os.Args[i]
			}
		}
		encode(soloops.GeneratePolicyPack(biz))
	case "cockpit-html":
		ledger := "testdata/inquiries.jsonl"
		biz := "Solo Ops"
		var out string
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--ledger":
				i++
				ledger = os.Args[i]
			case "--business":
				i++
				biz = os.Args[i]
			case "--out":
				i++
				out = os.Args[i]
			}
		}
		list, err := soloops.ListInquiries(ledger)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		html := soloops.RenderCockpitHTML(biz, list)
		if out != "" {
			if err := os.WriteFile(out, []byte(html), 0o644); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		fmt.Print(html)
	case "serve":
		ledger := "testdata/inquiries.jsonl"
		biz := "Solo Ops"
		addr := "127.0.0.1:8791"
		for i := 2; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--ledger":
				i++
				ledger = os.Args[i]
			case "--business":
				i++
				biz = os.Args[i]
			case "--addr":
				i++
				addr = os.Args[i]
			}
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			list, err := soloops.ListInquiries(ledger)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(soloops.RenderCockpitHTML(biz, list)))
		})
		fmt.Fprintln(os.Stderr, "solo-ops listening on http://"+addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown %q\n", os.Args[1])
		os.Exit(2)
	}
}

func encode(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
