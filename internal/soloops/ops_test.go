package soloops

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInquiryStatusAndCockpit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "inquiries.jsonl")
	inq := NewInquiry("inq-1", "Alex", "birthday", 12, "outdoor")
	if err := AppendInquiry(path, inq); err != nil {
		t.Fatal(err)
	}
	upd, err := SetStatus(path, "inq-1", "booked")
	if err != nil {
		t.Fatal(err)
	}
	if upd.Status != "booked" {
		t.Fatal(upd.Status)
	}
	list, err := ListInquiries(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Status != "booked" {
		t.Fatalf("%+v", list)
	}
	html := RenderCockpitHTML("Demo Co", list)
	if !strings.Contains(html, "Alex") || !strings.Contains(html, "booked") {
		t.Fatal(html)
	}
	pack := GeneratePolicyPack("Demo Co")
	if pack.Markdown == "" {
		t.Fatal("policy")
	}
}
