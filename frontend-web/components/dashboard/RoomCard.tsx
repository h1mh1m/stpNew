import React from 'react';
import Link from 'next/link';
import { MapPin, Users, ArrowRight } from 'lucide-react';
import { Room } from '@/lib/types';
import { formatRupiah } from '@/lib/datetime';

interface RoomCardProps {
  room: Room;
}

export const RoomCard: React.FC<RoomCardProps> = ({ room }) => {
  const bookingUrl = `/booking?room_id=${encodeURIComponent(room.id)}`;

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden hover:shadow-md transition-shadow flex flex-col justify-between">
      <div>
        <div className="h-64 bg-gray-200 relative overflow-hidden">
          <img
            src={room.image_url || 'https://images.unsplash.com/photo-1580274455191-1c62238fa333?auto=format&fit=crop&q=80&w=1000'}
            alt={room.name}
            className="w-full h-full object-cover"
          />
        </div>
        <div className="p-8 pb-4">
          <div className="inline-block px-4 py-1.5 bg-[#007bff] text-white text-xs font-semibold rounded-full mb-4">
            {room.features || 'Sirkuit Balap & Pengujian'}
          </div>
          <h3 className="text-2xl font-bold text-gray-900 mb-2">{room.name}</h3>
          <p className="text-gray-500 mb-6 text-sm leading-relaxed">{room.description}</p>

          <div className="flex flex-wrap items-center gap-6 text-xs text-gray-700 mb-6">
            <div className="flex items-center gap-2">
              <MapPin className="w-4 h-4 text-[#007bff] shrink-0" aria-hidden="true" />
              <span>Kampus ITS Sukolilo, Surabaya</span>
            </div>
            {room.capacity > 0 && (
              <div className="flex items-center gap-2">
                <Users className="w-4 h-4 text-[#007bff] shrink-0" aria-hidden="true" />
                <span>Kapasitas: {room.capacity} Orang</span>
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="p-8 pt-0">
        <div className="pt-6 border-t border-gray-100 flex items-center justify-between gap-4">
          <div>
            <span className="text-xs text-gray-500 block">Biaya Reservasi</span>
            <p className="text-3xl font-bold text-[#007bff] tracking-tight">
              {formatRupiah(room.price)}
            </p>
          </div>
          <Link
            href={bookingUrl}
            className="inline-flex items-center gap-2 bg-[#512BD4] hover:bg-[#4323af] text-white px-6 py-3 rounded-xl text-sm font-semibold transition-colors shadow-xs"
          >
            <span>Pesan Sirkuit</span>
            <ArrowRight className="w-4 h-4" aria-hidden="true" />
          </Link>
        </div>
      </div>
    </div>
  );
};
