import React from 'react';
import { CheckCircle2, Calendar, CreditCard, Hash, MapPin } from 'lucide-react';
import { Order, Room } from '@/lib/types';
import { formatIndonesianDateTime, formatRupiah } from '@/lib/datetime';
import { BookingStatusBadge } from './BookingStatusBadge';

interface BookingSummaryProps {
  order: Order;
  room?: Room | null;
  onReset?: () => void;
}

export const BookingSummary: React.FC<BookingSummaryProps> = ({
  order,
  room,
  onReset,
}) => {
  return (
    <div className="bg-white rounded-2xl border-2 border-emerald-500/30 p-6 md:p-8 shadow-sm space-y-6">
      <div className="flex items-start gap-4">
        <div className="w-12 h-12 rounded-full bg-emerald-100 flex items-center justify-center shrink-0">
          <CheckCircle2 className="w-6 h-6 text-emerald-600" aria-hidden="true" />
        </div>
        <div className="flex-1">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <h3 className="text-xl font-bold text-gray-900">
              Pemesanan Berhasil Dibuat
            </h3>
            <BookingStatusBadge status={order.status} />
          </div>
          <p className="text-sm text-gray-600 mt-1">
            Reservasi Anda telah tercatat di sistem kami dengan status{' '}
            <span className="font-semibold text-amber-700">Menunggu Pembayaran</span>.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 bg-gray-50 p-5 rounded-xl border border-gray-100 text-sm">
        <div className="space-y-1">
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <Hash className="w-3.5 h-3.5" aria-hidden="true" />
            <span>ID Pemesanan (Order ID)</span>
          </div>
          <p className="font-mono font-bold text-gray-900 break-all select-all">
            {order.id}
          </p>
        </div>

        <div className="space-y-1">
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <MapPin className="w-3.5 h-3.5" aria-hidden="true" />
            <span>Fasilitas</span>
          </div>
          <p className="font-semibold text-gray-900">
            {room?.name || `Fasilitas #${order.room_id}`}
          </p>
        </div>

        <div className="space-y-1">
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <Calendar className="w-3.5 h-3.5" aria-hidden="true" />
            <span>Waktu Penggunaan</span>
          </div>
          <p className="font-medium text-gray-900">
            Mulai: {formatIndonesianDateTime(order.date_start)}
          </p>
          <p className="font-medium text-gray-900">
            Selesai: {formatIndonesianDateTime(order.date_finish)}
          </p>
        </div>

        <div className="space-y-1">
          <div className="flex items-center gap-1.5 text-xs text-gray-500">
            <CreditCard className="w-3.5 h-3.5" aria-hidden="true" />
            <span>Total Tagihan (Server-Calculated)</span>
          </div>
          <p className="text-2xl font-bold text-[#512BD4]">
            {formatRupiah(order.total_price)}
          </p>
        </div>
      </div>

      <div className="bg-amber-50/80 border border-amber-200 rounded-xl p-4 text-xs text-amber-900 space-y-1">
        <p className="font-bold">⚠️ Informasi Pembayaran:</p>
        <p>
          Reservasi ini belum selesai dan memerlukan pembayaran. Simpan Order ID di atas. Penyelesaian pembayaran akan diaktifkan secara aman pada tahap rilis berikutnya.
        </p>
      </div>

      {onReset && (
        <div className="pt-2 flex justify-end">
          <button
            type="button"
            onClick={onReset}
            className="px-6 py-2.5 text-sm font-semibold text-[#512BD4] bg-white border border-[#512BD4] rounded-lg hover:bg-gray-50 transition-colors"
          >
            Buat Pemesanan Lain
          </button>
        </div>
      )}
    </div>
  );
};
