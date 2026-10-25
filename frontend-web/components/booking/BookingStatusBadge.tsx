import React from 'react';
import { Clock, CheckCircle2, XCircle, HelpCircle } from 'lucide-react';
import { OrderStatus } from '@/lib/types';

interface BookingStatusBadgeProps {
  status: OrderStatus | string;
  className?: string;
}

export const BookingStatusBadge: React.FC<BookingStatusBadgeProps> = ({
  status,
  className = '',
}) => {
  switch (status) {
    case 'Menunggu Pembayaran':
      return (
        <span
          className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-amber-50 text-amber-800 border border-amber-200 ${className}`}
        >
          <Clock className="w-3.5 h-3.5 text-amber-600" aria-hidden="true" />
          Menunggu Pembayaran
        </span>
      );
    case 'Selesai':
      return (
        <span
          className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-800 border border-emerald-200 ${className}`}
        >
          <CheckCircle2 className="w-3.5 h-3.5 text-emerald-600" aria-hidden="true" />
          Selesai
        </span>
      );
    case 'Dibatalkan':
      return (
        <span
          className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-gray-100 text-gray-700 border border-gray-300 ${className}`}
        >
          <XCircle className="w-3.5 h-3.5 text-gray-500" aria-hidden="true" />
          Dibatalkan
        </span>
      );
    default:
      return (
        <span
          className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-gray-100 text-gray-700 border border-gray-200 ${className}`}
        >
          <HelpCircle className="w-3.5 h-3.5 text-gray-400" aria-hidden="true" />
          {status}
        </span>
      );
  }
};
