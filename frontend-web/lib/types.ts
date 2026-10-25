/**
 * Shared TypeScript type definitions for the Sirkuit ITS / TrackIQ frontend.
 * Derived strictly from the authoritative PostgreSQL database models and backend JSON contracts.
 */

export type OrderStatus = 'Menunggu Pembayaran' | 'Selesai' | 'Dibatalkan';

/**
 * Room catalogue model matching backend schemas.Room JSON serialization.
 */
export interface Room {
  id: string;
  name: string;
  description: string;
  capacity: number;
  price: number;
  image_url: string;
  features: string;
  created_at?: string;
  updated_at?: string;
}

/**
 * Authoritative Order record returned by the backend (schemas.Order).
 */
export interface Order {
  id: string;
  user_id: string;
  room_id: string;
  date_start: string;
  date_finish: string;
  status: OrderStatus | string;
  total_price: number;
  created_at: string;
  updated_at: string;
}

/**
 * Public/Customer-visible occupied slot for schedule conflict presentation.
 * Customer identities and prices are strictly excluded for privacy.
 */
export interface AvailabilitySlot {
  room_id: string;
  date_start: string;
  date_finish: string;
  status: OrderStatus | string;
}

/**
 * Booking creation request payload for POST /api/Bookings.
 * Only room_id, date_start, and date_finish are client-settable.
 * User identity, status, and price are strictly server-derived.
 */
export interface BookingCreationRequest {
  room_id: string;
  date_start: string; // RFC 3339
  date_finish: string; // RFC 3339
}

/**
 * Standard backend response envelope.
 */
export interface ApiResponse<T> {
  message?: string;
  data: T;
}

/**
 * Standard backend error envelope.
 */
export interface ApiErrorResponse {
  error: string;
}

/**
 * Type guard for Room runtime validation.
 * Ensures the response strictly adheres to backend schemas.Room contract.
 */
export function isRoom(value: unknown): value is Room {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return false;
  }
  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.id === 'string' &&
    candidate.id.trim() !== '' &&
    typeof candidate.name === 'string' &&
    typeof candidate.description === 'string' &&
    typeof candidate.capacity === 'number' &&
    !isNaN(candidate.capacity) &&
    typeof candidate.price === 'number' &&
    !isNaN(candidate.price) &&
    typeof candidate.image_url === 'string' &&
    typeof candidate.features === 'string'
  );
}

/**
 * Validates and safely parses an unknown payload into a valid Room object.
 * Throws a descriptive TypeError if data is malformed.
 */
export function parseRoom(value: unknown): Room {
  if (!isRoom(value)) {
    throw new TypeError('Data ruangan dari server tidak valid atau format tidak sesuai.');
  }
  return value;
}

/**
 * Validates an array of room records from /api/Beranda.
 * If the input is not an array or contains invalid room records, throws TypeError.
 */
export function parseRoomCatalogue(data: unknown): Room[] {
  if (!Array.isArray(data)) {
    throw new TypeError('Katalog sirkuit dari server tidak valid (harus berupa array).');
  }
  return data.map((item, idx) => {
    if (!isRoom(item)) {
      throw new TypeError(`Data fasilitas sirkuit pada indeks ${idx} tidak valid atau rusak.`);
    }
    return item;
  });
}

/**
 * Type guard for AvailabilitySlot runtime validation.
 */
export function isAvailabilitySlot(value: unknown): value is AvailabilitySlot {
  if (!value || typeof value !== 'object' || Array.isArray(value)) {
    return false;
  }
  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.room_id === 'string' &&
    typeof candidate.date_start === 'string' &&
    typeof candidate.date_finish === 'string' &&
    typeof candidate.status === 'string'
  );
}

/**
 * Validates an array of availability slots.
 */
export function parseAvailabilitySlots(data: unknown): AvailabilitySlot[] {
  if (!Array.isArray(data)) {
    throw new TypeError('Data ketersediaan sirkuit tidak valid (harus berupa array).');
  }
  return data.map((slot, idx) => {
    if (!isAvailabilitySlot(slot)) {
      throw new TypeError(`Data slot ketersediaan pada indeks ${idx} tidak valid.`);
    }
    return slot;
  });
}
