import React from 'react';
import { History, RefreshCw, Calendar, Clock, MapPin } from 'lucide-react';
import { Order, Room } from '@/lib/types';
import { formatIndonesianDateTime, formatRupiah } from '@/lib/datetime';
import { BookingStatusBadge } from './BookingStatusBadge';
import { Alert } from '@/components/ui/Alert';

interface BookingHistoryProps {
  orders: Order[];
  rooms?: Room[];
  isLoading: boolean;
  error: string | null;
  onRefresh?: () => void;
}

export const BookingHistory: React.FC<BookingHistoryProps> = ({
  orders,
  rooms = [],
  isLoading,
  error,
  onRefresh,
}) => {
  const getRoomName = (roomId: string): string => {
    const found = rooms.find((r) => r.id === roomId);
    return found?.name || `Fasilitas #${roomId}`;
  };

  return (
    <div className="bg-white rounded-3xl p-6 md:p-8 border border-gray-100 shadow-sm space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2.5 bg-[#512BD4]/10 rounded-xl">
            <History className="w-5 h-5 text-[#512BD4]" aria-hidden="true" />
          </div>
          <div>
            <h3 className="text-xl font-bold text-gray-900">
              Riwayat Pemesanan Saya
            </h3>
            <p className="text-xs text-gray-500">
              Daftar reservasi fasilitas yang terdaftar pada akun Anda
            </p>
          </div>
        </div>
        {onRefresh && (
          <button
            type="button"
            onClick={onRefresh}
            disabled={isLoading}
            className="flex items-center gap-1.5 text-xs text-[#512BD4] hover:text-[#4323af] font-medium px-3 py-1.5 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-50 transition-colors"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${isLoading ? 'animate-spin' : ''}`} aria-hidden="true" />
            <span>Segarkan</span>
          </button>
        )}
      </div>

      {isLoading && (
        <div className="py-12 text-center text-sm text-gray-500">
          <div className="inline-block animate-spin w-6 h-6 border-2 border-[#512BD4] border-t-transparent rounded-full mb-3" />
          <p>Memuat riwayat pemesanan Anda...</p>
        </div>
      )}

      {error && !isLoading && (
        <Alert type="error" title="Gagal Memuat Riwayat">
          {error}
          {onRefresh && (
            <button
              type="button"
              onClick={onRefresh}
              className="mt-2 text-xs font-semibold underline block hover:text-red-900"
            >
              Coba lagi
            </button>
          )}
        </Alert>
      )}

      {!isLoading && !error && orders.length === 0 && (
        <div className="py-12 text-center text-sm text-gray-500 bg-gray-50/50 rounded-2xl border border-dashed border-gray-200">
          <Calendar className="w-8 h-8 text-gray-400 mx-auto mb-2" aria-hidden="true" />
          <p className="font-semibold text-gray-700">Belum Ada Riwayat Pemesanan</p>
          <p className="text-xs text-gray-400 mt-1">
            Anda belum pernah membuat reservasi sirkuit atau fasilitas.
          </p>
        </div>
      )}

      {!isLoading && !error && orders.length > 0 && (
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse" aria-label="Tabel riwayat pemesanan">
            <thead>
              <tr className="border-b border-gray-100 text-xs font-bold text-gray-500 uppercase tracking-wider bg-gray-50/50">
                <th scope="col" className="py-3 px-4 rounded-l-lg">ID Pemesanan</th>
                <th scope="col" className="py-3 px-4">Fasilitas</th>
                <th scope="col" className="py-3 px-4">Jadwal Penggunaan</th>
                <th scope="col" className="py-3 px-4">Total Tagihan</th>
                <th scope="col" className="py-3 px-4">Status</th>
                <th scope="col" className="py-3 px-4 rounded-r-lg">Dibuat Pada</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 text-sm">
              {orders.map((order) => (
                <tr key={order.id} className="hover:bg-gray-50/60 transition-colors">
                  <td className="py-4 px-4 font-mono text-xs font-medium text-gray-700 select-all max-w-[140px] truncate" title={order.id}>
                    {order.id}
                  </td>
                  <td className="py-4 px-4 font-medium text-gray-900">
                    <div className="flex items-center gap-1.5">
                      <MapPin className="w-3.5 h-3.5 text-[#007bff] shrink-0" aria-hidden="true" />
                      <span>{getRoomName(order.room_id)}</span>
                    </div>
                  </td>
                  <td className="py-4 px-4 text-xs text-gray-700 space-y-0.5">
                    <div className="font-medium text-gray-900">
                      Mulai: {formatIndonesianDateTime(order.date_start)}
                    </div>
                    <div className="text-gray-500">
                      Selesai: {formatIndonesianDateTime(order.date_finish)}
                    </div>
                  </td>
                  <td className="py-4 px-4 font-bold text-[#512BD4]">
                    {formatRupiah(order.total_price)}
                  </td>
                  <td className="py-4 px-4">
                    <BookingStatusBadge status={order.status} />
                  </td>
                  <td className="py-4 px-4 text-xs text-gray-500">
                    <div className="flex items-center gap-1">
                      <Clock className="w-3 h-3 text-gray-400" aria-hidden="true" />
                      <span>{formatIndonesianDateTime(order.created_at)}</span>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};
