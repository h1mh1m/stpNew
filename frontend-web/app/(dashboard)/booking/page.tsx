'use client';

import React, { useState, useEffect, useCallback, useRef, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft, Building2 } from 'lucide-react';
import { Room, AvailabilitySlot, Order, ApiResponse, parseRoomCatalogue, parseAvailabilitySlots } from '@/lib/types';
import { authenticatedFetch, ApiError } from '@/lib/api-client';
import { AvailabilityList } from '@/components/booking/AvailabilityList';
import { BookingForm } from '@/components/booking/BookingForm';
import { BookingSummary } from '@/components/booking/BookingSummary';
import { BookingHistory } from '@/components/booking/BookingHistory';
import { Alert } from '@/components/ui/Alert';

function BookingContent() {
  const searchParams = useSearchParams();
  const roomId = searchParams.get('room_id');

  const [room, setRoom] = useState<Room | null>(null);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [availabilitySlots, setAvailabilitySlots] = useState<AvailabilitySlot[]>([]);
  const [orders, setOrders] = useState<Order[]>([]);
  const [createdOrder, setCreatedOrder] = useState<Order | null>(null);

  const [isLoadingRoom, setIsLoadingRoom] = useState<boolean>(true);
  const [isLoadingAvailability, setIsLoadingAvailability] = useState<boolean>(() => Boolean(roomId));
  const [isLoadingOrders, setIsLoadingOrders] = useState<boolean>(true);

  const [roomError, setRoomError] = useState<string | null>(null);
  const [availabilityError, setAvailabilityError] = useState<string | null>(null);
  const [ordersError, setOrdersError] = useState<string | null>(null);

  // Track the active availability window to dynamically align ranges
  const loadedWindowRef = useRef<{ start: number; finish: number } | null>(null);

  const getLoadedWindow = useCallback(() => {
    if (!loadedWindowRef.current) {
      const now = Date.now();
      loadedWindowRef.current = {
        start: now - 5 * 60 * 1000,
        finish: now + 90 * 24 * 60 * 60 * 1000,
      };
    }
    return loadedWindowRef.current;
  }, []);

  // Manual refresh callbacks (triggered by user action or booking completion)
  const refreshOrders = useCallback(() => {
    setIsLoadingOrders(true);
    setOrdersError(null);
    authenticatedFetch<ApiResponse<Order[]>>('/api/Bookings')
      .then((res) => {
        setOrders(res.data || []);
      })
      .catch((err: unknown) => {
        setOrdersError(err instanceof ApiError ? err.message : 'Gagal memuat riwayat pemesanan Anda.');
      })
      .finally(() => {
        setIsLoadingOrders(false);
      });
  }, []);

  const refreshAvailability = useCallback(
    (range?: { start: string; finish: string }) => {
      if (!roomId) return;
      setIsLoadingAvailability(true);
      setAvailabilityError(null);

      let path = `/api/rooms/${encodeURIComponent(roomId)}/availability`;
      if (range?.start && range?.finish) {
        const params = new URLSearchParams({ start: range.start, finish: range.finish });
        path += `?${params.toString()}`;
      }

      authenticatedFetch<ApiResponse<unknown>>(path)
        .then((res) => {
          const validated = parseAvailabilitySlots(res.data);
          setAvailabilitySlots(validated);
          if (range?.start && range?.finish) {
            const current = getLoadedWindow();
            loadedWindowRef.current = {
              start: Math.min(current.start, new Date(range.start).getTime()),
              finish: Math.max(current.finish, new Date(range.finish).getTime()),
            };
          }
        })
        .catch((err: unknown) => {
          if (err instanceof ApiError) {
            setAvailabilityError(err.message);
          } else if (err instanceof TypeError) {
            setAvailabilityError(err.message);
          } else {
            setAvailabilityError('Gagal memuat ketersediaan jadwal sirkuit.');
          }
        })
        .finally(() => {
          setIsLoadingAvailability(false);
        });
    },
    [roomId, getLoadedWindow]
  );

  // Dynamic range alignment: if customer selects dates outside currently loaded window, refetch bounded range
  const handleEnsureRange = useCallback(
    (range: { start: string; finish: string }) => {
      if (!roomId) return;
      const startMs = new Date(range.start).getTime();
      const finishMs = new Date(range.finish).getTime();
      if (isNaN(startMs) || isNaN(finishMs)) return;

      const current = getLoadedWindow();
      if (startMs < current.start || finishMs > current.finish) {
        // Expand window bounded by user selection, with 24h buffer
        const queryStart = new Date(Math.min(current.start, startMs - 24 * 60 * 60 * 1000)).toISOString();
        const queryFinish = new Date(Math.max(current.finish, finishMs + 24 * 60 * 60 * 1000)).toISOString();
        refreshAvailability({ start: queryStart, finish: queryFinish });
      }
    },
    [roomId, getLoadedWindow, refreshAvailability]
  );

  // Initial load effects
  useEffect(() => {
    let ignore = false;
    authenticatedFetch<ApiResponse<Order[]>>('/api/Bookings')
      .then((res) => {
        if (!ignore) {
          setOrders(res.data || []);
          setIsLoadingOrders(false);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          setOrdersError(err instanceof ApiError ? err.message : 'Gagal memuat riwayat pemesanan Anda.');
          setIsLoadingOrders(false);
        }
      });

    authenticatedFetch<ApiResponse<unknown>>('/api/Beranda')
      .then((res) => {
        if (!ignore) {
          const allRooms = parseRoomCatalogue(res.data);
          setRooms(allRooms);
          if (roomId) {
            const found = allRooms.find((r) => r.id === roomId);
            if (found) {
              setRoom(found);
            } else {
              setRoomError(`Fasilitas sirkuit dengan ID "${roomId}" tidak ditemukan.`);
            }
          }
          setIsLoadingRoom(false);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          if (err instanceof ApiError) {
            setRoomError(err.message);
          } else if (err instanceof TypeError) {
            setRoomError(err.message);
          } else {
            setRoomError('Gagal mengambil data katalog fasilitas.');
          }
          setIsLoadingRoom(false);
        }
      });

    return () => {
      ignore = true;
    };
  }, [roomId]);

  useEffect(() => {
    if (!roomId) {
      return;
    }
    let ignore = false;
    authenticatedFetch<ApiResponse<unknown>>(
      `/api/rooms/${encodeURIComponent(roomId)}/availability`
    )
      .then((res) => {
        if (!ignore) {
          const validated = parseAvailabilitySlots(res.data);
          setAvailabilitySlots(validated);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          if (err instanceof ApiError) {
            setAvailabilityError(err.message);
          } else if (err instanceof TypeError) {
            setAvailabilityError(err.message);
          } else {
            setAvailabilityError('Gagal memuat ketersediaan jadwal sirkuit.');
          }
        }
      })
      .finally(() => {
        if (!ignore) {
          setIsLoadingAvailability(false);
        }
      });

    return () => {
      ignore = true;
    };
  }, [roomId]);

  const handleBookingSuccess = (order: Order) => {
    setCreatedOrder(order);
    refreshAvailability();
    refreshOrders();
  };

  return (
    <div className="max-w-7xl mx-auto pb-12 space-y-10">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <Link
              href="/"
              className="inline-flex items-center gap-1.5 text-xs font-semibold text-gray-500 hover:text-gray-900 transition-colors"
            >
              <ArrowLeft className="w-3.5 h-3.5" />
              <span>Kembali ke Beranda</span>
            </Link>
          </div>
          <h1 className="text-3xl font-bold text-gray-900">Pemesanan Sirkuit ITS</h1>
          <p className="text-sm text-gray-500 mt-1">
            Lakukan reservasi fasilitas sirkuit secara langsung dan terintegrasi
          </p>
        </div>
      </div>

      {/* Case 1: No room_id in search params */}
      {!roomId && (
        <div className="bg-white p-8 md:p-12 rounded-3xl shadow-sm border border-gray-100 text-center space-y-4">
          <div className="w-16 h-16 bg-purple-50 rounded-2xl flex items-center justify-center mx-auto text-[#512BD4]">
            <Building2 className="w-8 h-8" aria-hidden="true" />
          </div>
          <h3 className="text-xl font-bold text-gray-900">Pilih Fasilitas Terlebih Dahulu</h3>
          <p className="text-sm text-gray-500 max-w-md mx-auto">
            Untuk melakukan pemesanan jadwal sirkuit, silakan tentukan fasilitas yang ingin Anda sewa melalui katalog Beranda.
          </p>
          <div className="pt-2">
            <Link
              href="/"
              className="inline-flex items-center gap-2 bg-[#512BD4] hover:bg-[#4323af] text-white px-6 py-3 rounded-xl text-sm font-semibold transition-colors"
            >
              <span>Buka Katalog Fasilitas</span>
            </Link>
          </div>
        </div>
      )}

      {/* Case 2: room_id present, loading room data */}
      {roomId && isLoadingRoom && (
        <div className="bg-white p-12 rounded-3xl shadow-sm border border-gray-100 text-center text-gray-500">
          <div className="inline-block animate-spin w-8 h-8 border-4 border-[#512BD4] border-t-transparent rounded-full mb-3" />
          <p className="text-sm font-medium">Memuat data fasilitas sirkuit...</p>
        </div>
      )}

      {/* Case 3: room_id present, room not found or error */}
      {roomId && !isLoadingRoom && roomError && (
        <div className="bg-white p-8 rounded-3xl shadow-sm border border-gray-100 space-y-4">
          <Alert type="error" title="Fasilitas Tidak Ditemukan">
            {roomError}
          </Alert>
          <div className="pt-2">
            <Link
              href="/"
              className="inline-flex items-center gap-2 text-sm font-semibold text-[#512BD4] hover:underline"
            >
              <ArrowLeft className="w-4 h-4" />
              <span>Kembali ke Katalog Sirkuit</span>
            </Link>
          </div>
        </div>
      )}

      {/* Case 4: room_id present, room loaded successfully */}
      {roomId && !isLoadingRoom && room && (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* Main Booking Form / Summary Section */}
          <div className="lg:col-span-7 bg-white p-6 md:p-8 rounded-3xl shadow-sm border border-gray-100">
            {createdOrder ? (
              <BookingSummary
                order={createdOrder}
                room={room}
                onReset={() => setCreatedOrder(null)}
              />
            ) : (
              <div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">Formulir Reservasi</h2>
                <p className="text-sm text-gray-500 mb-6">
                  Pilih waktu mulai dan selesai penggunaan. Sistem akan memverifikasi ketersediaan jadwal secara real-time.
                </p>
                <BookingForm
                  room={room}
                  availabilitySlots={availabilitySlots}
                  onBookingSuccess={handleBookingSuccess}
                  onRefreshAvailability={refreshAvailability}
                  onEnsureRange={handleEnsureRange}
                />
              </div>
            )}
          </div>

          {/* Availability Side Panel */}
          <div className="lg:col-span-5">
            <AvailabilityList
              slots={availabilitySlots}
              isLoading={isLoadingAvailability}
              error={availabilityError}
              onRefresh={refreshAvailability}
            />
          </div>
        </div>
      )}

      {/* Booking History Section: scoped strictly to authenticated customer */}
      <BookingHistory
        orders={orders}
        rooms={rooms}
        isLoading={isLoadingOrders}
        error={ordersError}
        onRefresh={refreshOrders}
      />
    </div>
  );
}

export default function BookingPage() {
  return (
    <Suspense
      fallback={
        <div className="max-w-7xl mx-auto py-20 text-center text-gray-500">
          <div className="inline-block animate-spin w-8 h-8 border-4 border-[#512BD4] border-t-transparent rounded-full mb-3" />
          <p className="text-sm">Memuat halaman pemesanan...</p>
        </div>
      }
    >
      <BookingContent />
    </Suspense>
  );
}
