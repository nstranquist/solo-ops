package soloops

import (
	"fmt"
	"html"
	"strings"
	"time"
)

// RenderCockpitHTML builds a static, read-only HTML view of inquiries.
func RenderCockpitHTML(business string, list []Inquiry) string {
	if strings.TrimSpace(business) == "" {
		business = "Solo Ops"
	}
	counts := map[string]int{}
	for _, inq := range list {
		counts[statusKind(inq.Status)]++
	}

	var b strings.Builder
	b.WriteString("<!doctype html><html lang=en><head><meta charset=utf-8>")
	b.WriteString("<meta name=viewport content=\"width=device-width, initial-scale=1\">")
	b.WriteString("<title>Solo Ops — " + html.EscapeString(business) + "</title>")
	b.WriteString(cockpitCSS)
	b.WriteString("</head><body><div class=wrap>")
	b.WriteString("<header class=hero>")
	b.WriteString("<div class=kicker>Solo Ops · offline cockpit</div>")
	b.WriteString("<h1>" + html.EscapeString(business) + "</h1>")
	b.WriteString("<p class=lede>Read-only view of the local JSONL ledger. Nothing is uploaded.</p>")
	b.WriteString("</header>")

	b.WriteString("<section class=stats aria-label=Summary>")
	writeStat(&b, "inquiries", fmt.Sprintf("%d", len(list)))
	for _, key := range []string{"new", "review", "booked", "done", "declined"} {
		if n := counts[key]; n > 0 {
			writeStat(&b, key, fmt.Sprintf("%d", n))
		}
	}
	if n := counts["other"]; n > 0 {
		writeStat(&b, "other", fmt.Sprintf("%d", n))
	}
	b.WriteString("</section>")

	b.WriteString("<section class=board>")
	b.WriteString("<div class=board-head><h2>Inquiries</h2><span class=hint>append-only ledger</span></div>")
	if len(list) == 0 {
		b.WriteString("<p class=empty>No inquiries in this ledger.</p>")
	} else {
		b.WriteString("<ul class=cards>")
		for _, inq := range list {
			kind := statusKind(inq.Status)
			b.WriteString("<li class=\"card status-" + kind + "\">")
			b.WriteString("<div class=card-top>")
			b.WriteString("<div class=who>" + html.EscapeString(displayName(inq.Name)) + "</div>")
			b.WriteString("<span class=pill>" + html.EscapeString(displayStatus(inq.Status)) + "</span>")
			b.WriteString("</div>")
			b.WriteString("<div class=event>" + html.EscapeString(displayEvent(inq.Event)) + "</div>")
			b.WriteString("<div class=meta>")
			b.WriteString("<span>" + guestLabel(inq.Guests) + "</span>")
			if d := formatDay(inq.CreatedAt); d != "" {
				b.WriteString("<span>" + html.EscapeString(d) + "</span>")
			}
			if inq.ID != "" {
				b.WriteString("<span class=id>" + html.EscapeString(inq.ID) + "</span>")
			}
			b.WriteString("</div>")
			if strings.TrimSpace(inq.Notes) != "" {
				b.WriteString("<p class=notes>" + html.EscapeString(inq.Notes) + "</p>")
			}
			b.WriteString("</li>")
		}
		b.WriteString("</ul>")
	}
	b.WriteString("</section>")
	b.WriteString("<footer>Local file only · loopback serve · no account</footer>")
	b.WriteString("</div></body></html>")
	return b.String()
}

func writeStat(b *strings.Builder, label, value string) {
	b.WriteString("<div class=stat><div class=stat-n>" + html.EscapeString(value) + "</div>")
	b.WriteString("<div class=stat-l>" + html.EscapeString(label) + "</div></div>")
}

func guestLabel(n int) string {
	if n == 1 {
		return "1 guest"
	}
	return fmt.Sprintf("%d guests", n)
}

func displayName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "Unnamed"
	}
	return s
}

func displayEvent(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "No event listed"
	}
	return s
}

func displayStatus(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	return s
}

func statusKind(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "new":
		return "new"
	case "review":
		return "review"
	case "booked":
		return "booked"
	case "done", "closed", "complete":
		return "done"
	case "declined", "lost":
		return "declined"
	default:
		return "other"
	}
}

func formatDay(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		if len(raw) >= 10 {
			return raw[:10]
		}
		return raw
	}
	return t.UTC().Format("2 Jan 2006")
}

const cockpitCSS = `<style>
:root {
  --bg:#f4efe6; --ink:#1c1915; --muted:#6b645b; --line:#ddd4c6;
  --paper:#fffaf2; --paper-2:#f7f1e6; --accent:#2f5d50;
  --new:#1f6feb; --review:#9a6700; --booked:#1a7f37;
  --done:#6b645b; --declined:#b42318; --other:#8250df;
}
* { box-sizing:border-box; }
html,body { margin:0; background:var(--bg); color:var(--ink);
  font:16px/1.5 -apple-system,BlinkMacSystemFont,"Iowan Old Style","Palatino Linotype",Georgia,serif; }
.wrap { max-width:980px; margin:0 auto; padding:40px 28px 64px; }
.kicker { font:12px/1.3 ui-monospace,Menlo,monospace; letter-spacing:.08em;
  text-transform:uppercase; color:var(--accent); margin-bottom:10px; }
h1 { font-size:clamp(32px,5vw,44px); margin:0 0 8px; letter-spacing:-.03em; font-weight:650; }
.lede { margin:0; color:var(--muted); max-width:46ch; }
.stats { display:flex; flex-wrap:wrap; gap:10px; margin:28px 0 32px; }
.stat { background:var(--paper); border:1px solid var(--line); border-radius:12px;
  padding:12px 16px; min-width:108px; }
.stat-n { font-size:24px; font-weight:650; letter-spacing:-.03em; }
.stat-l { font:12px/1.3 ui-monospace,Menlo,monospace; color:var(--muted); text-transform:uppercase; letter-spacing:.04em; }
.board-head { display:flex; align-items:baseline; justify-content:space-between; gap:12px; margin-bottom:12px; }
h2 { font-size:18px; margin:0; }
.hint { font:12px/1.3 ui-monospace,Menlo,monospace; color:var(--muted); }
.empty { color:var(--muted); padding:28px; background:var(--paper); border:1px dashed var(--line); border-radius:12px; }
.cards { list-style:none; margin:0; padding:0; display:grid; gap:12px; }
.card { background:var(--paper); border:1px solid var(--line); border-radius:14px; padding:16px 18px; }
.card-top { display:flex; justify-content:space-between; gap:12px; align-items:center; }
.who { font-weight:650; font-size:18px; }
.event { margin-top:4px; }
.meta { display:flex; flex-wrap:wrap; gap:10px 14px; margin-top:8px;
  font:12px/1.3 ui-monospace,Menlo,monospace; color:var(--muted); }
.id { opacity:.8; }
.notes { margin:10px 0 0; color:var(--muted); }
.pill { font:11px/1 ui-monospace,Menlo,monospace; text-transform:uppercase; letter-spacing:.06em;
  border-radius:999px; padding:5px 9px; border:1px solid var(--line); background:var(--paper-2); }
.status-new .pill { color:var(--new); border-color:color-mix(in srgb, var(--new) 35%, var(--line)); }
.status-review .pill { color:var(--review); border-color:color-mix(in srgb, var(--review) 35%, var(--line)); }
.status-booked .pill { color:var(--booked); border-color:color-mix(in srgb, var(--booked) 35%, var(--line)); }
.status-done .pill { color:var(--done); }
.status-declined .pill { color:var(--declined); border-color:color-mix(in srgb, var(--declined) 35%, var(--line)); }
footer { margin-top:36px; color:var(--muted); font:12px/1.4 ui-monospace,Menlo,monospace; }
@media (max-width:640px) {
  .wrap { padding:28px 16px 48px; }
  .card-top { align-items:flex-start; }
}
</style>`
