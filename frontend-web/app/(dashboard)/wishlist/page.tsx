'use client';

import React from 'react';
import { Heart } from 'lucide-react';
import { RoomCard } from '@/components/dashboard/RoomCard';
import Link from 'next/link';

const wishlistRooms = [
  {
    id: 1,
    title: 'Lab Pemrograman 1',
    department: 'Departemen Teknik Informatika ITS',
    location: 'Zona A',
    capacity: 25,
    price: 'Rp 90.000',
    unit: 'per Jam',
    tag: 'Jangka Pendek'
  },
  {
    id: 3,
    title: 'Gedung NASDEC',
    department: 'Institut Teknologi Sepuluh Nopember',
    location: 'Zona C',
    capacity: 100,
    price: 'Rp 5.000.000',
    unit: 'per Bulan',
    tag: 'Jangka Panjang'
  }
];

export default function WishlistPage() {
  return (
    <div className="max-w-6xl mx-auto pb-12">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Wishlist</h1>
        <p className="text-gray-500">Ruangan yang Anda simpan untuk dipertimbangkan</p>
      </div>

      {wishlistRooms.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {wishlistRooms.map((room) => (
            <div key={room.id} className="relative group">
              <div className="absolute top-4 right-4 z-10 bg-white p-2 rounded-full shadow-md text-red-500 hover:text-gray-400 cursor-pointer transition-colors">
                <Heart className="w-5 h-5 fill-current" />
              </div>
              <Link href={`/rooms/${room.id}`} className="block">
                <RoomCard {...room} />
              </Link>
            </div>
          ))}
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
            <Heart className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-bold text-gray-900 mb-2">Wishlist Anda masih kosong</h3>
          <p className="text-gray-500 mb-6">Anda belum menyimpan ruangan apapun. Temukan ruangan yang cocok untuk Anda.</p>
          <Link href="/rooms" className="inline-flex bg-[#0055b8] hover:bg-[#004494] text-white px-6 py-2.5 rounded-lg font-semibold transition-colors">
            Cari Ruangan
          </Link>
        </div>
      )}
    </div>
  );
}
