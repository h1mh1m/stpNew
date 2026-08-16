'use client';

import React, { useState } from 'react';
import { Search, Filter, ChevronDown } from 'lucide-react';
import { RoomCard } from '@/components/dashboard/RoomCard';

const allRooms = [
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
  {
    id: 5,
    title: 'Ruang Rapat Utama',
    department: 'Rektorat ITS',
    location: 'Zona B',
    capacity: 40,
    price: 'Rp 200.000',
    unit: 'per Jam',
    tag: 'Jangka Pendek'
  },
  {
    id: 6,
    title: 'Auditorium Pascasarjana',
    department: 'Sekolah Pascasarjana ITS',
    location: 'Zona B',
    capacity: 200,
    price: 'Rp 15.000.000',
    unit: 'per Hari',
    tag: 'Jangka Pendek'
  }
];

export default function RoomsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [activeFilter, setActiveFilter] = useState('Semua');

  const filteredRooms = allRooms.filter(room => {
    const matchesSearch = room.title.toLowerCase().includes(searchTerm.toLowerCase()) || 
                          room.department.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesFilter = activeFilter === 'Semua' || room.tag === activeFilter;
    
    return matchesSearch && matchesFilter;
  });

  return (
    <div className="max-w-6xl mx-auto pb-12">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Daftar Ruangan</h1>
        <p className="text-gray-500">Temukan ruangan yang sesuai dengan kebutuhan Anda</p>
      </div>

      {/* Search and Filters Bar */}
      <div className="bg-white p-4 rounded-xl border border-gray-200 shadow-sm mb-8 flex flex-col md:flex-row md:items-center justify-between gap-4">
        {/* Search Input */}
        <div className="relative flex-1 max-w-md">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search className="h-5 w-5 text-gray-400" />
          </div>
          <input
            type="text"
            className="block w-full pl-10 pr-3 py-2.5 border border-gray-300 rounded-lg leading-5 bg-gray-50 placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-[#0055b8] focus:border-[#0055b8] sm:text-sm text-gray-900"
            placeholder="Cari nama ruangan atau departemen..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        {/* Filters */}
        <div className="flex items-center space-x-2 overflow-x-auto pb-2 md:pb-0">
          <button 
            onClick={() => setActiveFilter('Semua')}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              activeFilter === 'Semua' 
                ? 'bg-[#0055b8] text-white' 
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            Semua
          </button>
          <button 
            onClick={() => setActiveFilter('Jangka Pendek')}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              activeFilter === 'Jangka Pendek' 
                ? 'bg-[#0055b8] text-white' 
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            Jangka Pendek
          </button>
          <button 
            onClick={() => setActiveFilter('Jangka Panjang')}
            className={`px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
              activeFilter === 'Jangka Panjang' 
                ? 'bg-[#0055b8] text-white' 
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
          >
            Jangka Panjang
          </button>
          
          <div className="w-px h-8 bg-gray-300 mx-2 hidden md:block"></div>
          
          <button className="flex items-center px-4 py-2 bg-white border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors whitespace-nowrap">
            <Filter className="w-4 h-4 mr-2" />
            Filter Lainnya
            <ChevronDown className="w-4 h-4 ml-2" />
          </button>
        </div>
      </div>

      {/* Room Grid */}
      {filteredRooms.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {filteredRooms.map((room) => (
            <RoomCard key={room.id} {...room} />
          ))}
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
          <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-100 mb-4">
            <Search className="w-8 h-8 text-gray-400" />
          </div>
          <h3 className="text-lg font-bold text-gray-900 mb-2">Ruangan tidak ditemukan</h3>
          <p className="text-gray-500">Coba sesuaikan kata kunci pencarian atau filter Anda.</p>
        </div>
      )}
    </div>
  );
}
