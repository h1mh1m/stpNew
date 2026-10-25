import React, { useState, useEffect, useRef } from 'react';
import { Calendar as CalendarIcon, Clock, ShieldCheck, X } from 'lucide-react';
import { Room, AvailabilitySlot, BookingCreationRequest, Order } from '@/lib/types';
import { localDatetimeToISO, findConflictingSlot, formatRupiah, formatIndonesianDateTime } from '@/lib/datetime';
import { authenticatedFetch, ApiError } from '@/lib/api-client';
import { Alert } from '@/components/ui/Alert';

interface BookingFormProps {
  room: Room;
  availabilitySlots: AvailabilitySlot[];
  onBookingSuccess: (order: Order) => void;
  onRefreshAvailability?: (range?: { start: string; finish: string }) => void;
  onEnsureRange?: (range: { start: string; finish: string }) => void;
}

export const BookingForm: React.FC<BookingFormProps> = ({
  room,
  availabilitySlots,
  onBookingSuccess,
  onRefreshAvailability,
  onEnsureRange,
}) => {
  const [dateStart, setDateStart] = useState<string>('');
  const [dateFinish, setDateFinish] = useState<string>('');
  const [isTosModalOpen, setIsTosModalOpen] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [validationError, setValidationError] = useState<string | null>(null);
  const [apiError, setApiError] = useState<string | null>(null);
  const [conflictError, setConflictError] = useState<string | null>(null);

  // Focus trap and submission refs
  const isSubmittingRef = useRef<boolean>(false);
  const openTermsButtonRef = useRef<HTMLButtonElement | null>(null);
  const closeButtonRef = useRef<HTMLButtonElement | null>(null);
  const batalButtonRef = useRef<HTMLButtonElement | null>(null);
  const modalRef = useRef<HTMLDivElement | null>(null);

  // Range alignment: if selected dates extend beyond default range, ensure availability window covers them
  useEffect(() => {
    if (!dateStart || !dateFinish) return;
    try {
      const startIso = localDatetimeToISO(dateStart);
      const finishIso = localDatetimeToISO(dateFinish);
      if (new Date(finishIso).getTime() > new Date(startIso).getTime()) {
        onEnsureRange?.({ start: startIso, finish: finishIso });
      }
    } catch {
      // Ignored while user is actively typing
    }
  }, [dateStart, dateFinish, onEnsureRange]);

  // Modal accessibility: focus management, Tab trap, and ESC key listener
  useEffect(() => {
    if (!isTosModalOpen) return;

    // Focus initial safe action (close button) when modal opens
    const focusTimer = setTimeout(() => {
      closeButtonRef.current?.focus();
    }, 40);

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        e.preventDefault();
        setIsTosModalOpen(false);
        openTermsButtonRef.current?.focus();
        return;
      }

      if (e.key === 'Tab' && modalRef.current) {
        const focusable = modalRef.current.querySelectorAll<HTMLElement>(
          'button:not([disabled]), [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
        );
        if (focusable.length === 0) return;

        const first = focusable[0];
        const last = focusable[focusable.length - 1];

        if (e.shiftKey) {
          if (document.activeElement === first) {
            e.preventDefault();
            last.focus();
          }
        } else {
          if (document.activeElement === last) {
            e.preventDefault();
            first.focus();
          }
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      clearTimeout(focusTimer);
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [isTosModalOpen]);

  // Client-side pre-validation
  const validateForm = (): { startIso: string; finishIso: string } | null => {
    setValidationError(null);
    setConflictError(null);
    setApiError(null);

    if (!dateStart) {
      setValidationError('Silakan pilih tanggal dan waktu mulai reservasi.');
      return null;
    }
    if (!dateFinish) {
      setValidationError('Silakan pilih tanggal dan waktu selesai reservasi.');
      return null;
    }

    let startIso: string;
    let finishIso: string;
    try {
      startIso = localDatetimeToISO(dateStart);
      finishIso = localDatetimeToISO(dateFinish);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Format tanggal tidak valid';
      setValidationError(msg);
      return null;
    }

    const startDate = new Date(startIso);
    const finishDate = new Date(finishIso);
    const now = new Date();

    if (finishDate.getTime() <= startDate.getTime()) {
      setValidationError('Waktu selesai harus lebih lambat dari waktu mulai.');
      return null;
    }

    // 5-minute tolerance for clock differences between client and server
    if (startDate.getTime() < now.getTime() - 5 * 60 * 1000) {
      setValidationError('Waktu mulai reservasi tidak boleh berada di masa lampau.');
      return null;
    }

    // Visible pre-conflict check against loaded availability slots
    const conflict = findConflictingSlot(startIso, finishIso, availabilitySlots);
    if (conflict) {
      setConflictError(
        `Jadwal yang Anda pilih bentrok dengan reservasi yang sudah ada (${formatIndonesianDateTime(
          conflict.date_start
        )} s/d ${formatIndonesianDateTime(conflict.date_finish)}). Silakan pilih waktu lain.`
      );
      return null;
    }

    return { startIso, finishIso };
  };

  const handleOpenTermsModal = (e: React.FormEvent) => {
    e.preventDefault();
    if (isSubmittingRef.current) return;
    const validated = validateForm();
    if (validated) {
      setIsTosModalOpen(true);
    }
  };

  const handleCloseModal = () => {
    if (isSubmittingRef.current) return;
    setIsTosModalOpen(false);
    openTermsButtonRef.current?.focus();
  };

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget && !isSubmittingRef.current) {
      handleCloseModal();
    }
  };

  const handleConfirmAndSubmit = async () => {
    // In-flight guard: prevent duplicate submission in the same render cycle
    if (isSubmittingRef.current) {
      return;
    }

    const validated = validateForm();
    if (!validated) {
      handleCloseModal();
      return;
    }

    isSubmittingRef.current = true;
    setIsSubmitting(true);
    setIsTosModalOpen(false);
    setApiError(null);
    setConflictError(null);

    const payload: BookingCreationRequest = {
      room_id: room.id,
      date_start: validated.startIso,
      date_finish: validated.finishIso,
    };

    try {
      const res = await authenticatedFetch<{ message: string; data: Order }>(
        '/api/Bookings',
        {
          method: 'POST',
          body: JSON.stringify(payload),
        }
      );

      // Successfully created order
      onBookingSuccess(res.data);
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        if (err.status === 409) {
          setConflictError(
            'Jadwal bentrok: Sirkuit baru saja dipesan oleh pemesan lain untuk jadwal tersebut. Jadwal ketersediaan telah diperbarui.'
          );
          onRefreshAvailability?.({ start: validated.startIso, finish: validated.finishIso });
        } else if (err.status === 404) {
          setApiError('Fasilitas sirkuit yang dipilih tidak ditemukan dalam katalog.');
        } else if (err.status === 400) {
          setApiError(err.message || 'Permintaan tidak valid. Periksa kembali jadwal yang dimasukkan.');
        } else if (err.status >= 500) {
          setApiError('Terjadi kesalahan pada server. Silakan coba kembali dalam beberapa saat.');
        } else {
          setApiError(err.message);
        }
      } else {
        setApiError('Gagal mengirimkan pemesanan. Periksa koneksi internet Anda.');
      }
    } finally {
      isSubmittingRef.current = false;
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Validation & Conflict Feedback Alerts */}
      {validationError && (
        <Alert type="error" title="Validasi Jadwal">
          {validationError}
        </Alert>
      )}

      {conflictError && (
        <Alert type="warning" title="Jadwal Tidak Tersedia">
          {conflictError}
        </Alert>
      )}

      {apiError && (
        <Alert type="error" title="Kesalahan Pemesanan">
          {apiError}
        </Alert>
      )}

      {/* Selected Facility Read-Only Summary */}
      <div className="bg-purple-50/60 border border-purple-200/80 rounded-2xl p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <span className="text-xs font-bold uppercase tracking-wider text-[#512BD4]">
            Fasilitas Terpilih
          </span>
          <h4 className="text-lg font-bold text-gray-900 mt-0.5">{room.name}</h4>
          <p className="text-xs text-gray-600 mt-1 max-w-xl">{room.description}</p>
        </div>
        <div className="text-right sm:text-right shrink-0">
          <span className="text-xs text-gray-500 block">Biaya Sewa</span>
          <span className="text-2xl font-bold text-[#512BD4]">
            {formatRupiah(room.price)}
          </span>
        </div>
      </div>

      <form onSubmit={handleOpenTermsModal} noValidate className="space-y-6">
        <div>
          <h4 className="text-sm font-bold text-gray-900 mb-1">
            Tentukan Waktu Penggunaan Sirkuit <span className="text-red-500">*</span>
          </h4>
          <p className="text-xs text-gray-500 mb-4">
            Waktu yang dipilih disesuaikan dengan waktu lokal perangkat Anda.
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label htmlFor="date-start" className="block text-xs font-semibold text-gray-700 mb-2">
                Tanggal dan Waktu Mulai
              </label>
              <div className="relative">
                <CalendarIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" aria-hidden="true" />
                <input
                  id="date-start"
                  type="datetime-local"
                  required
                  value={dateStart}
                  onChange={(e) => {
                    setDateStart(e.target.value);
                    setValidationError(null);
                    setConflictError(null);
                  }}
                  className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white transition-colors"
                />
              </div>
            </div>

            <div>
              <label htmlFor="date-finish" className="block text-xs font-semibold text-gray-700 mb-2">
                Tanggal dan Waktu Berakhir
              </label>
              <div className="relative">
                <Clock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" aria-hidden="true" />
                <input
                  id="date-finish"
                  type="datetime-local"
                  required
                  value={dateFinish}
                  onChange={(e) => {
                    setDateFinish(e.target.value);
                    setValidationError(null);
                    setConflictError(null);
                  }}
                  className="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg text-sm text-gray-900 focus:outline-none focus:ring-2 focus:ring-[#512BD4]/20 focus:border-[#512BD4] bg-white transition-colors"
                />
              </div>
            </div>
          </div>
        </div>

        <div className="pt-4 flex flex-col items-center gap-3">
          <button
            ref={openTermsButtonRef}
            type="submit"
            disabled={isSubmitting}
            className="w-full sm:w-auto bg-[#512BD4] hover:bg-[#4323af] text-white px-10 py-3 rounded-lg font-medium transition-colors shadow-sm disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            {isSubmitting ? (
              <>
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                <span>Memproses Reservasi...</span>
              </>
            ) : (
              <span>Lanjutkan ke Syarat & Ketentuan</span>
            )}
          </button>
          <p className="text-xs text-gray-400 text-center">
            Pemesanan akan dibuat dengan status &quot;Menunggu Pembayaran&quot; setelah Anda menyetujui ketentuan.
          </p>
        </div>
      </form>

      {/* Terms of Service Modal */}
      {isTosModalOpen && (
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby="tos-modal-title"
          aria-describedby="tos-modal-description"
          onClick={handleBackdropClick}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-xs"
        >
          <div
            ref={modalRef}
            onClick={(e) => e.stopPropagation()}
            className="bg-white rounded-3xl p-6 md:p-8 max-w-lg w-full shadow-2xl relative animate-in fade-in zoom-in-95 duration-150"
          >
            <button
              ref={closeButtonRef}
              type="button"
              onClick={handleCloseModal}
              className="absolute top-6 right-6 text-gray-400 hover:text-gray-600 p-1 rounded-lg hover:bg-gray-100 transition-colors"
              aria-label="Tutup jendela syarat dan ketentuan"
            >
              <X className="w-6 h-6" />
            </button>

            <div className="flex items-center gap-2.5 mb-4">
              <div className="p-2 bg-[#512BD4]/10 rounded-xl">
                <ShieldCheck className="w-6 h-6 text-[#512BD4]" aria-hidden="true" />
              </div>
              <h3 id="tos-modal-title" className="text-xl font-bold text-gray-900">
                Ketentuan Layanan Sirkuit ITS
              </h3>
            </div>

            <div
              id="tos-modal-description"
              className="text-gray-600 text-xs sm:text-sm space-y-3 max-h-[50vh] overflow-y-auto mb-6 pr-2 border-y border-gray-100 py-4"
            >
              <p>1. Pemesan wajib menjaga kebersihan dan fasilitas sirkuit selama masa sewa.</p>
              <p>2. Segala bentuk kerusakan yang diakibatkan kelalaian pengguna menjadi tanggung jawab pemesan sepenuhnya.</p>
              <p>3. Dilarang keras membawa atau mengonsumsi minuman keras dan obat-obatan terlarang di seluruh area sirkuit.</p>
              <p>4. Pemesan harus mematuhi jadwal yang telah disepakati dan tidak melewati batas waktu reservasi.</p>
              <p>5. Pembatalan pemesanan harus dilakukan minimal 24 jam sebelum jadwal penggunaan.</p>
              <p>6. Pengelola Sirkuit ITS tidak bertanggung jawab atas kehilangan barang pribadi.</p>
              <p>7. Reservasi yang dibuat berstatus Menunggu Pembayaran sampai verifikasi pembayaran selesai.</p>
            </div>

            <div className="flex flex-col-reverse sm:flex-row items-center justify-end gap-3">
              <button
                ref={batalButtonRef}
                type="button"
                onClick={handleCloseModal}
                className="w-full sm:w-auto px-5 py-2.5 rounded-lg border border-gray-300 text-gray-700 text-sm font-medium hover:bg-gray-50 transition-colors"
              >
                Batal
              </button>
              <button
                type="button"
                onClick={handleConfirmAndSubmit}
                disabled={isSubmitting}
                className="w-full sm:w-auto bg-[#512BD4] hover:bg-[#4323af] text-white px-6 py-2.5 rounded-lg text-sm font-medium transition-colors shadow-sm disabled:opacity-50"
              >
                Saya Setuju &amp; Kirim Pemesanan
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
