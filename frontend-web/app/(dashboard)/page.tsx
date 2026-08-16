import React from 'react';
import { ChevronRight } from 'lucide-react';
import { RoomCard } from '@/components/dashboard/RoomCard';
import Link from 'next/link';

const rooms = [
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
    id: 2,
    title: 'Lab Pemrograman 2',
    department: 'Departemen Teknik Informatika ITS',
    location: 'Zona A',
    capacity: 25,
    price: 'Rp 100.000',
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
  },
  {
    id: 4,
    title: 'Menara Sains',
    department: 'Institut Teknologi Sepuluh Nopember',
    location: 'Zona C',
    capacity: 50,
    price: 'Rp 10.000.000',
    unit: 'per Tahun',
    tag: 'Jangka Panjang'
  },
];

export default function Dashboard() {
  return (
    <div className="max-w-6xl mx-auto pb-12">
      {/* Hero Section */}
      <div className="relative w-full h-[350px] overflow-hidden rounded-2xl mb-12">
        <div className="absolute inset-0 bg-[#512BD4]"></div>
        <div className="absolute inset-0 flex flex-col justify-center px-10 md:px-16 text-white">
          <h1 className="text-3xl md:text-4xl font-bold mb-4 max-w-2xl">Selamat Datang di Wiranala!</h1>
          <p className="text-lg text-white/90 max-w-xl mb-8">
            Temukan ruangan dengan kualitas dan harga terbaik untuk menunjang segala kebutuhan Anda
          </p>
          <Link href="/rooms" className="inline-flex items-center justify-center bg-[#007bff] hover:bg-[#0056b3] text-white px-6 py-3 rounded-lg font-semibold w-max transition-colors">
            Cari Ruangan
            <ChevronRight className="w-5 h-5 ml-2" />
          </Link>
        </div>
      </div>

      {/* Tempat Kerja yang Nyaman Section */}
      <div className="bg-white p-8 rounded-2xl border border-gray-200">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Tempat Kerja yang Nyaman</h2>
        <p className="text-gray-500 mb-8">Temukan berbagai pilihan ruangan kantor yang nyaman untuk menunjang kebutuhan kerja Anda</p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
          {rooms.map((room) => (
            <RoomCard key={room.id} {...room} />
          ))}
        </div>

        <div className="flex justify-center">
          <Link href="/rooms" className="bg-[#007bff] hover:bg-[#0056b3] text-white px-8 py-3 rounded-lg font-semibold transition-colors">
            Lihat Semua Ruangan
          </Link>
        </div>
      </div>
    </div>
  );
}
