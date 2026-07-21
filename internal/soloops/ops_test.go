package soloops

import (
	"path/filepath"
	"testing"
)

func TestInquiryLedgerAndPolicy(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inquiries.jsonl")
	inq := NewInquiry("inq-1", "Alex", "birthday", 12, "outdoor")
	if err := AppendInquiry(path, inq); err != nil {
		t.Fatal(err)
	}
	list, err := ListInquiries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "Alex" {
		t.Fatalf("%+v", list)
	}
	pack := GeneratePolicyPack("Painted by Demo")
	if pack.Markdown == "" || len(pack.Sections) < 3 {
		t.Fatalf("%+v", pack)
	}
	if pack.Business != "Painted by Demo" {
		t.Fatal(pack.Business)
	}
}
