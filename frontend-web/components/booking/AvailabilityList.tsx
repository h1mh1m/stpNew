import React from 'react';
import { Calendar, RefreshCw } from 'lucide-react';
import { AvailabilitySlot } from '@/lib/types';
import { formatIndonesianDateTime } from '@/lib/datetime';
import { BookingStatusBadge } from './BookingStatusBadge';
import { Alert } from '@/components/ui/Alert';

interface AvailabilityListProps {
  slots: AvailabilitySlot[];
  isLoading: boolean;
  error: string | null;
  onRefresh?: () => void;
}

export const AvailabilityList: React.FC<AvailabilityListProps> = ({
  slots,
  isLoading,
  error,
  onRefresh,
}) => {
  return (
    <div className="bg-gray-50/70 border border-gray-200 rounded-2xl p-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Calendar className="w-5 h-5 text-[#512BD4]" aria-hidden="true" />
          <h4 className="text-base font-bold text-gray-900">
            Jadwal Terisi Saat Ini
          </h4>
        </div>
        {onRefresh && (
          <button
            type="button"
            onClick={onRefresh}
            disabled={isLoading}
            className="flex items-center gap-1.5 text-xs text-[#512BD4] hover:text-[#4323af] font-medium p-1.5 rounded-lg hover:bg-gray-100 disabled:opacity-50 transition-colors"
            title="Muat ulang ketersediaan jadwal"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${isLoading ? 'animate-spin' : ''}`} aria-hidden="true" />
            <span>Segarkan</span>
          </button>
        )}
      </div>

      <p className="text-xs text-gray-500 mb-4">
        Daftar reservasi aktif di bawah ini merupakan jadwal yang memblokir ketersediaan sirkuit. Demi privasi, identitas pemesan lain tidak ditampilkan.
      </p>

      {isLoading && (
        <div className="py-6 text-center text-sm text-gray-500">
          <div className="inline-block animate-spin w-5 h-5 border-2 border-[#512BD4] border-t-transparent rounded-full mb-2" />
          <p>Memeriksa ketersediaan jadwal sirkuit...</p>
        </div>
      )}

      {error && !isLoading && (
        <Alert type="error" title="Gagal Memuat Ketersediaan">
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

      {!isLoading && !error && slots.length === 0 && (
        <div className="py-6 text-center text-sm text-gray-500 bg-white rounded-xl border border-dashed border-gray-200">
          <p className="font-medium text-gray-700">Belum ada pemesanan aktif</p>
          <p className="text-xs text-gray-400 mt-1">Sirkuit sepenuhnya tersedia untuk tanggal dan jam yang Anda inginkan.</p>
        </div>
      )}

      {!isLoading && !error && slots.length > 0 && (
        <ul className="space-y-3" aria-label="Daftar jadwal terisi">
          {slots.map((slot, index) => (
            <li
              key={`${slot.room_id}-${slot.date_start}-${index}`}
              className="bg-white p-3.5 rounded-xl border border-gray-200 flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-xs"
            >
              <div className="space-y-1">
                <div className="flex items-center gap-2 text-xs font-semibold text-gray-900">
                  <span className="w-2 h-2 rounded-full bg-amber-500 inline-block" />
                  <span>Mulai: {formatIndonesianDateTime(slot.date_start)}</span>
                </div>
                <div className="text-xs text-gray-600 pl-4">
                  <span>Selesai: {formatIndonesianDateTime(slot.date_finish)}</span>
                </div>
              </div>
              <div className="self-start sm:self-center">
                <BookingStatusBadge status={slot.status} />
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
