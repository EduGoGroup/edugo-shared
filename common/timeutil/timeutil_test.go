package timeutil_test

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/common/timeutil"
)

func TestNowUTC(t *testing.T) {
	now := timeutil.NowUTC()

	if now.Location() != time.UTC {
		t.Errorf("NowUTC debe estar en UTC, got %v", now.Location())
	}
}

func TestFormatISO(t *testing.T) {
	t.Run("zona_no_utc_produce_sufijo_Z", func(t *testing.T) {
		// 10:00 en una zona +05:00 == 05:00 UTC.
		loc := time.FixedZone("PLUS5", 5*60*60)
		instant := time.Date(2026, 5, 15, 10, 0, 0, 0, loc)

		got := timeutil.FormatISO(instant)

		want := "2026-05-15T05:00:00Z"
		if got != want {
			t.Errorf("FormatISO: want %q, got %q", want, got)
		}
	})

	t.Run("ya_en_utc", func(t *testing.T) {
		instant := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)

		got := timeutil.FormatISO(instant)

		want := "2026-05-15T00:00:00Z"
		if got != want {
			t.Errorf("FormatISO: want %q, got %q", want, got)
		}
	})
}

func TestParseISO(t *testing.T) {
	t.Run("utc_sin_corrimiento", func(t *testing.T) {
		got, err := timeutil.ParseISO("2026-05-15T00:00:00Z")
		if err != nil {
			t.Fatalf("ParseISO falló: %v", err)
		}

		if got.Location() != time.UTC {
			t.Errorf("ParseISO debe devolver UTC, got %v", got.Location())
		}

		want := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("ParseISO: want %v, got %v", want, got)
		}
	})

	t.Run("con_offset_se_normaliza_a_utc", func(t *testing.T) {
		got, err := timeutil.ParseISO("2026-05-15T10:00:00+05:00")
		if err != nil {
			t.Fatalf("ParseISO falló: %v", err)
		}

		if got.Location() != time.UTC {
			t.Errorf("ParseISO debe devolver UTC, got %v", got.Location())
		}

		want := time.Date(2026, 5, 15, 5, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("ParseISO: want %v, got %v", want, got)
		}
	})

	t.Run("invalido_devuelve_error", func(t *testing.T) {
		_, err := timeutil.ParseISO("no-es-una-fecha")
		if err == nil {
			t.Error("ParseISO debe devolver error con entrada inválida")
		}
	})
}

func TestFormatDate(t *testing.T) {
	d := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)

	got := timeutil.FormatDate(d)

	want := "2026-05-15"
	if got != want {
		t.Errorf("FormatDate: want %q, got %q", want, got)
	}
}

func TestParseDate(t *testing.T) {
	t.Run("medianoche_utc_sin_corrimiento", func(t *testing.T) {
		got, err := timeutil.ParseDate("2026-05-15")
		if err != nil {
			t.Fatalf("ParseDate falló: %v", err)
		}

		want := time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC)
		if !got.Equal(want) {
			t.Errorf("ParseDate: want %v, got %v", want, got)
		}
	})

	t.Run("invalido_devuelve_error", func(t *testing.T) {
		_, err := timeutil.ParseDate("2026/05/15")
		if err == nil {
			t.Error("ParseDate debe devolver error con formato inválido")
		}
	})
}

func TestDateRoundTrip(t *testing.T) {
	const in = "2026-05-15"

	parsed, err := timeutil.ParseDate(in)
	if err != nil {
		t.Fatalf("ParseDate falló: %v", err)
	}

	out := timeutil.FormatDate(parsed)
	if out != in {
		t.Errorf("round-trip de fecha: want %q, got %q (corrimiento de día)", in, out)
	}
}

func TestISORoundTrip(t *testing.T) {
	const in = "2026-05-15T08:30:00Z"

	parsed, err := timeutil.ParseISO(in)
	if err != nil {
		t.Fatalf("ParseISO falló: %v", err)
	}

	out := timeutil.FormatISO(parsed)
	if out != in {
		t.Errorf("round-trip ISO: want %q, got %q", in, out)
	}
}
