import test from 'node:test';
import assert from 'node:assert/strict';

import {
  localDatetimeToISO,
  isoToLocalDatetimeInput,
  hasIntervalOverlap,
  findConflictingSlot,
  formatRupiah,
  formatIndonesianDateTime,
} from './datetime.ts';

test('localDatetimeToISO converts local datetime to RFC 3339 UTC ISO', () => {
  const localInput = '2026-11-01T14:30';
  const iso = localDatetimeToISO(localInput);
  
  assert.ok(iso.endsWith('Z'), 'ISO string should end with Z');
  const d = new Date(iso);
  assert.equal(d.getMinutes(), 30);
});

test('localDatetimeToISO rejects empty and invalid inputs', () => {
  assert.throws(() => localDatetimeToISO(''), /tidak boleh kosong/);
  assert.throws(() => localDatetimeToISO('invalid-string'), /harus YYYY-MM-DDTHH:mm/);
  assert.throws(() => localDatetimeToISO('2026-13-45T99:99'), /tidak valid/);
});

test('isoToLocalDatetimeInput converts ISO/Date to YYYY-MM-DDTHH:mm', () => {
  const d = new Date(2026, 10, 1, 14, 30);
  const formatted = isoToLocalDatetimeInput(d);
  assert.equal(formatted, '2026-11-01T14:30');
});

test('hasIntervalOverlap detects overlapping intervals correctly', () => {
  // [10:00, 12:00) vs [11:00, 13:00) -> Overlap
  assert.equal(
    hasIntervalOverlap('2026-11-01T10:00:00Z', '2026-11-01T12:00:00Z', '2026-11-01T11:00:00Z', '2026-11-01T13:00:00Z'),
    true
  );

  // [10:00, 12:00) vs [12:00, 14:00) -> Adjacent, NO overlap (half-open)
  assert.equal(
    hasIntervalOverlap('2026-11-01T10:00:00Z', '2026-11-01T12:00:00Z', '2026-11-01T12:00:00Z', '2026-11-01T14:00:00Z'),
    false
  );

  // [10:00, 12:00) vs [08:00, 10:00) -> Adjacent, NO overlap
  assert.equal(
    hasIntervalOverlap('2026-11-01T10:00:00Z', '2026-11-01T12:00:00Z', '2026-11-01T08:00:00Z', '2026-11-01T10:00:00Z'),
    false
  );

  // [10:00, 12:00) vs [14:00, 16:00) -> Disjoint, NO overlap
  assert.equal(
    hasIntervalOverlap('2026-11-01T10:00:00Z', '2026-11-01T12:00:00Z', '2026-11-01T14:00:00Z', '2026-11-01T16:00:00Z'),
    false
  );
});

test('findConflictingSlot returns the conflicting slot when overlapping', () => {
  const slots = [
    {
      room_id: 'room-1',
      date_start: '2026-11-01T10:00:00Z',
      date_finish: '2026-11-01T12:00:00Z',
      status: 'Menunggu Pembayaran',
    },
    {
      room_id: 'room-1',
      date_start: '2026-11-01T14:00:00Z',
      date_finish: '2026-11-01T16:00:00Z',
      status: 'Selesai',
    },
  ];

  // Request [11:30, 13:00) overlaps first slot
  const conflict = findConflictingSlot('2026-11-01T11:30:00Z', '2026-11-01T13:00:00Z', slots);
  assert.ok(conflict !== null);
  assert.equal(conflict?.status, 'Menunggu Pembayaran');

  // Request [12:00, 14:00) is adjacent to both, no conflict
  const noConflict = findConflictingSlot('2026-11-01T12:00:00Z', '2026-11-01T14:00:00Z', slots);
  assert.equal(noConflict, null);
});

test('formatRupiah formats IDR currency correctly', () => {
  const formatted = formatRupiah(150000);
  assert.ok(formatted.includes('150.000'), `Expected 150.000 in ${formatted}`);
  assert.ok(formatted.includes('Rp'), `Expected Rp in ${formatted}`);
});

test('formatIndonesianDateTime formats date with timezone', () => {
  const formatted = formatIndonesianDateTime('2026-11-01T14:00:00Z');
  assert.ok(formatted.length > 0);
});

test('localDatetimeToISO strictly rejects non-existent calendar dates and invalid leap days', () => {
  // February 30 does not exist
  assert.throws(() => localDatetimeToISO('2026-02-30T10:00'), /tidak valid/);
  // April 31 does not exist (April has 30 days)
  assert.throws(() => localDatetimeToISO('2026-04-31T10:00'), /tidak valid/);
  // February 29 in non-leap year 2026 does not exist
  assert.throws(() => localDatetimeToISO('2026-02-29T10:00'), /tidak valid/);
  // February 29 in leap year 2024 DOES exist and must succeed
  const leapIso = localDatetimeToISO('2024-02-29T10:00');
  assert.ok(leapIso.endsWith('Z'), 'Leap year ISO should end in Z');

  // Missing time
  assert.throws(() => localDatetimeToISO('2026-02-15'), /harus YYYY-MM-DDTHH:mm/);
  // Partial inputs
  assert.throws(() => localDatetimeToISO('2026-02'), /harus YYYY-MM-DDTHH:mm/);
  assert.throws(() => localDatetimeToISO('2026-02-15T'), /harus YYYY-MM-DDTHH:mm/);
  assert.throws(() => localDatetimeToISO('2026-02-15T10'), /harus YYYY-MM-DDTHH:mm/);
  // Non-numeric tokens
  assert.throws(() => localDatetimeToISO('2026-02-15Taa:bb'), /harus YYYY-MM-DDTHH:mm/);
});

import { isRoom, parseRoom, parseRoomCatalogue } from './types.ts';

test('isRoom and parseRoom validate Go Room model JSON contracts', () => {
  const validRoom = {
    id: 'sirkuit-utama',
    name: 'Sirkuit Utama TrackIQ',
    description: 'Sirkuit aspal 1.2km',
    capacity: 50,
    price: 150000,
    image_url: 'https://example.com/circuit.jpg',
    features: 'Pit stop, timing gate',
    created_at: '2026-10-25T00:00:00Z',
    updated_at: '2026-10-25T00:00:00Z',
  };

  assert.equal(isRoom(validRoom), true);
  assert.deepEqual(parseRoom(validRoom), validRoom);

  // Missing id
  assert.equal(isRoom({ ...validRoom, id: '' }), false);
  // price is string instead of number
  assert.equal(isRoom({ ...validRoom, price: '150000' }), false);
  // capacity is string instead of number
  assert.equal(isRoom({ ...validRoom, capacity: '50' }), false);
  // missing required field
  const { name: _unusedName, ...missingName } = validRoom;
  void _unusedName;
  assert.equal(isRoom(missingName), false);
  assert.throws(() => parseRoom(missingName), /Data ruangan dari server tidak valid/);
});

test('parseRoomCatalogue validates array of rooms and rejects malformed collections', () => {
  const validCatalogue = [
    {
      id: 'room-1',
      name: 'Room 1',
      description: 'Desc 1',
      capacity: 20,
      price: 150000,
      image_url: 'https://example.com/1.jpg',
      features: 'Feature 1',
    },
  ];

  assert.deepEqual(parseRoomCatalogue(validCatalogue), validCatalogue);

  // Non-array
  assert.throws(() => parseRoomCatalogue(null), /harus berupa array/);
  assert.throws(() => parseRoomCatalogue({}), /harus berupa array/);

  // Array containing corrupted object
  const corruptedCatalogue = [
    validCatalogue[0],
    { id: 'corrupted', name: 123 }, // invalid
  ];
  assert.throws(() => parseRoomCatalogue(corruptedCatalogue), /Data fasilitas sirkuit pada indeks 1 tidak valid/);
});
