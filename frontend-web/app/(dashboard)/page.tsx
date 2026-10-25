'use client';

import React, { useEffect, useState } from 'react';
import { RoomCard } from '@/components/dashboard/RoomCard';
import { Room, ApiResponse, parseRoomCatalogue } from '@/lib/types';
import { authenticatedFetch, ApiError } from '@/lib/api-client';
import { Alert } from '@/components/ui/Alert';
import { RefreshCw } from 'lucide-react';

export default function Dashboard() {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let ignore = false;
    authenticatedFetch<ApiResponse<unknown>>('/api/Beranda')
      .then((res) => {
        if (!ignore) {
          const validatedRooms = parseRoomCatalogue(res.data);
          setRooms(validatedRooms);
          setIsLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!ignore) {
          if (err instanceof ApiError) {
            setError(err.message);
          } else if (err instanceof TypeError) {
            setError(err.message);
          } else {
            setError('Gagal memuat katalog sirkuit. Pastikan server backend sedang berjalan.');
          }
          setIsLoading(false);
        }
      });
    return () => {
      ignore = true;
    };
  }, []);

  const handleRetry = () => {
    setIsLoading(true);
    setError(null);
    authenticatedFetch<ApiResponse<unknown>>('/api/Beranda')
      .then((res) => {
        const validatedRooms = parseRoomCatalogue(res.data);
        setRooms(validatedRooms);
      })
      .catch((err: unknown) => {
        if (err instanceof ApiError) {
          setError(err.message);
        } else if (err instanceof TypeError) {
          setError(err.message);
        } else {
          setError('Gagal memuat katalog sirkuit. Pastikan server backend sedang berjalan.');
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  };

  return (
    <div className="max-w-7xl mx-auto pb-12 space-y-8">
      {/* Hero Section */}
      <div className="w-full h-[400px] rounded-3xl overflow-hidden relative bg-gray-200">
        <img
          src="https://images.unsplash.com/photo-1541339907198-e08756dedf3f?auto=format&fit=crop&q=80&w=2070"
          alt="Campus Building"
          className="w-full h-full object-cover"
        />
      </div>

      {/* Main Content */}
      <div className="bg-white p-8 md:p-12 rounded-3xl shadow-sm border border-gray-100">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-10">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2 flex items-center gap-3">
              <span>🏁</span> Temukan Sirkuit Terbaik untuk Kebutuhan Anda
            </h2>
            <p className="text-gray-600 text-base md:text-lg">
              Pilih sirkuit yang sesuai dengan kebutuhan Anda, lengkap dengan informasi fasilitas, lokasi, dan harga sewa.
            </p>
          </div>
          {error && (
            <button
              type="button"
              onClick={handleRetry}
              className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-[#512BD4] bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors shrink-0"
            >
              <RefreshCw className="w-4 h-4" />
              <span>Coba Lagi</span>
            </button>
          )}
        </div>

        {/* Loading State */}
        {isLoading && (
          <div className="py-20 text-center text-gray-500">
            <div className="inline-block animate-spin w-8 h-8 border-4 border-[#512BD4] border-t-transparent rounded-full mb-4" />
            <p className="text-base font-medium">Memuat katalog sirkuit dari server...</p>
          </div>
        )}

        {/* Error State */}
        {error && !isLoading && (
          <Alert type="error" title="Gagal Memuat Katalog Sirkuit">
            {error}
          </Alert>
        )}

        {/* Empty State */}
        {!isLoading && !error && rooms.length === 0 && (
          <div className="py-16 text-center text-gray-500 bg-gray-50 rounded-2xl border border-dashed border-gray-200">
            <p className="text-lg font-bold text-gray-800">Tidak ada fasilitas sirkuit yang tersedia</p>
            <p className="text-sm text-gray-500 mt-1">Saat ini belum ada data sirkuit yang terdaftar di basis data.</p>
          </div>
        )}

        {/* Success State: Circuit Cards Grid */}
        {!isLoading && !error && rooms.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {rooms.map((room) => (
              <RoomCard key={room.id} room={room} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
