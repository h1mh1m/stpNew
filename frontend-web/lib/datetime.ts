/**
 * Pure date, time, and currency utilities for TrackIQ / Sirkuit ITS frontend.
 * - Converts HTML datetime-local inputs strictly via browser local time to RFC 3339 UTC ISO strings.
 * - Formats dates and currency according to Indonesian standards (id-ID).
 * - Pre-validates scheduling conflicts using epoch millisecond comparisons.
 */

import type { AvailabilitySlot } from './types.ts';

/**
 * Converts an HTML datetime-local string (e.g. "2026-11-01T14:30") in the runtime's
 * local timezone to an RFC 3339 UTC ISO string (e.g. "2026-11-01T07:30:00.000Z").
 *
 * Uses explicit date parts with local constructor to prevent date-shifting quirks.
 */
export function localDatetimeToISO(localStr: string): string {
  if (!localStr || typeof localStr !== 'string') {
    throw new Error('Nilai tanggal dan waktu tidak boleh kosong');
  }

  // Strict check for HTML datetime-local format: YYYY-MM-DDTHH:mm (or with seconds)
  const formatRegex = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}(:\d{2})?$/;
  if (!formatRegex.test(localStr)) {
    throw new Error('Format tanggal dan waktu tidak valid (harus YYYY-MM-DDTHH:mm)');
  }

  const [datePart, timePart] = localStr.split('T');
  const [yearStr, monthStr, dayStr] = datePart.split('-');
  const timePieces = timePart.split(':');
  const [hourStr, minuteStr] = timePieces;

  const year = Number(yearStr);
  const month = Number(monthStr);
  const day = Number(dayStr);
  const hours = Number(hourStr);
  const minutes = Number(minuteStr);

  if (
    month < 1 ||
    month > 12 ||
    day < 1 ||
    day > 31 ||
    hours < 0 ||
    hours > 23 ||
    minutes < 0 ||
    minutes > 59
  ) {
    throw new Error('Format tanggal dan waktu tidak valid: nilai di luar rentang');
  }

  const date = new Date(year, month - 1, day, hours, minutes, 0, 0);

  if (
    isNaN(date.getTime()) ||
    date.getFullYear() !== year ||
    date.getMonth() !== month - 1 ||
    date.getDate() !== day ||
    date.getHours() !== hours ||
    date.getMinutes() !== minutes
  ) {
    throw new Error('Tanggal dan waktu tidak valid (tanggal kalender tidak ada)');
  }

  return date.toISOString();
}

/**
 * Converts a Date object or ISO string to the HTML datetime-local input format (YYYY-MM-DDTHH:mm)
 * in the local browser timezone.
 */
export function isoToLocalDatetimeInput(isoOrDate: string | Date): string {
  const d = typeof isoOrDate === 'string' ? new Date(isoOrDate) : isoOrDate;
  if (isNaN(d.getTime())) {
    return '';
  }

  const pad = (n: number) => n.toString().padStart(2, '0');
  const year = d.getFullYear();
  const month = pad(d.getMonth() + 1);
  const day = pad(d.getDate());
  const hours = pad(d.getHours());
  const mins = pad(d.getMinutes());

  return `${year}-${month}-${day}T${hours}:${mins}`;
}

/**
 * Formats an ISO datetime string into readable Indonesian localization with local timezone indicator.
 * e.g. "Minggu, 1 November 2026, 14:00 (WIB/GMT+7)"
 */
export function formatIndonesianDateTime(isoString: string): string {
  const d = new Date(isoString);
  if (isNaN(d.getTime())) {
    return isoString;
  }

  const formatter = new Intl.DateTimeFormat('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
    timeZoneName: 'short',
  });

  return formatter.format(d);
}

/**
 * Formats numeric monetary amounts as Indonesian Rupiah (e.g. "Rp 150.000").
 */
export function formatRupiah(amount: number): string {
  if (isNaN(amount) || amount === null || amount === undefined) {
    return 'Rp 0';
  }
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount);
}

/**
 * Checks if interval [startA, finishA) overlaps with interval [startB, finishB).
 * Half-open interval semantics: start is inclusive, finish is exclusive.
 */
export function hasIntervalOverlap(
  startA: Date | string | number,
  finishA: Date | string | number,
  startB: Date | string | number,
  finishB: Date | string | number
): boolean {
  const sA = new Date(startA).getTime();
  const fA = new Date(finishA).getTime();
  const sB = new Date(startB).getTime();
  const fB = new Date(finishB).getTime();

  if (isNaN(sA) || isNaN(fA) || isNaN(sB) || isNaN(fB)) {
    return false;
  }

  return sA < fB && fA > sB;
}

/**
 * Finds the first conflicting occupied slot from a list of availability slots.
 */
export function findConflictingSlot(
  startIso: string,
  finishIso: string,
  slots: AvailabilitySlot[]
): AvailabilitySlot | null {
  if (!slots || slots.length === 0) {
    return null;
  }

  for (const slot of slots) {
    if (hasIntervalOverlap(startIso, finishIso, slot.date_start, slot.date_finish)) {
      return slot;
    }
  }

  return null;
}
