package schemas

import (
	"testing"
)

func TestOrderStatuses(t *testing.T) {
	if StatusMenungguPembayaran != "Menunggu Pembayaran" {
		t.Fatalf("expected 'Menunggu Pembayaran', got '%s'", StatusMenungguPembayaran)
	}
	if StatusSelesai != "Selesai" {
		t.Fatalf("expected 'Selesai', got '%s'", StatusSelesai)
	}
	if StatusDibatalkan != "Dibatalkan" {
		t.Fatalf("expected 'Dibatalkan', got '%s'", StatusDibatalkan)
	}
}

func TestBlockingStatuses(t *testing.T) {
	if !IsBlockingStatus(StatusMenungguPembayaran) {
		t.Error("expected StatusMenungguPembayaran to be blocking")
	}
	if !IsBlockingStatus(StatusSelesai) {
		t.Error("expected StatusSelesai to be blocking")
	}
	if IsBlockingStatus(StatusDibatalkan) {
		t.Error("expected StatusDibatalkan to be non-blocking")
	}
	if IsBlockingStatus("RandomStatus") {
		t.Error("expected unknown status to be non-blocking")
	}
	if IsBlockingStatus("") {
		t.Error("expected empty status to be non-blocking")
	}
}
