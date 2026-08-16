'use client';

import React from 'react';
import Link from 'next/link';
import { ArrowLeft, MapPin, Users, Wifi, Monitor, Coffee, Wind, Calendar, Clock } from 'lucide-react';
import { Button } from '@/components/ui/Button';

// Dummy data for preview purposes
const roomDetail = {
  id: 1,
  title: 'Lab Pemrograman 1',
  department: 'Departemen Teknik Informatika ITS',
  location: 'Zona A - Gedung Teknik Informatika Lt. 2',
  capacity: 25,
  price: 'Rp 90.000',
  unit: 'per Jam',
  tag: 'Jangka Pendek',
  description: 'Laboratorium Pemrograman 1 dilengkapi dengan fasilitas komputer berspesifikasi tinggi, cocok untuk pelatihan, praktikum, atau kegiatan workshop yang membutuhkan perangkat memadai. Ruangan ini ber-AC dan memiliki proyektor.',
  amenities: [
    { icon: Wifi, name: 'Free WiFi' },
    { icon: Monitor, name: '25 PC Desktop' },
    { icon: Wind, name: 'AC' },
    { icon: Coffee, name: 'Pantry Area' },
  ]
};

export default function RoomDetailPage({ params }: { params: { id: string } }) {
  // Normally we would fetch the room details using params.id
  
  return (
    <div className="max-w-6xl mx-auto pb-12">
      {/* Back button */}
      <div className="mb-6">
        <Link href="/rooms" className="inline-flex items-center text-sm font-medium text-gray-500 hover:text-[#0055b8] transition-colors">
          <ArrowLeft className="w-4 h-4 mr-2" />
          Kembali ke Daftar Ruangan
        </Link>
      </div>

      {/* Hero Image (Purple Block) */}
      <div className="w-full h-[400px] bg-[#512BD4] rounded-2xl mb-8 flex items-center justify-center">
        {/* Placeholder for actual room image */}
        <span className="text-white/50 text-xl font-medium tracking-wider uppercase">Gambar Ruangan</span>
      </div>

      <div className="flex flex-col lg:flex-row gap-8">
        {/* Left Column: Details */}
        <div className="flex-1">
          <div className="inline-block bg-[#0055b8] text-white text-xs font-semibold px-3 py-1 rounded-md mb-4">
            {roomDetail.tag}
          </div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">{roomDetail.title}</h1>
          <p className="text-lg text-gray-600 mb-6">{roomDetail.department}</p>
          
          <div className="flex flex-wrap items-center gap-6 mb-8 pb-8 border-b border-gray-200">
            <div className="flex items-center text-gray-700">
              <MapPin className="w-5 h-5 mr-2 text-[#0055b8]" />
              {roomDetail.location}
            </div>
            <div className="flex items-center text-gray-700">
              <Users className="w-5 h-5 mr-2 text-[#0055b8]" />
              Kapasitas {roomDetail.capacity} Orang
            </div>
          </div>

          <div className="mb-8">
            <h2 className="text-xl font-bold text-gray-900 mb-4">Deskripsi Ruangan</h2>
            <p className="text-gray-600 leading-relaxed">
              {roomDetail.description}
            </p>
          </div>

          <div>
            <h2 className="text-xl font-bold text-gray-900 mb-4">Fasilitas Utama</h2>
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-4">
              {roomDetail.amenities.map((item, index) => {
                const Icon = item.icon;
                return (
                  <div key={index} className="flex items-center p-3 rounded-lg border border-gray-200 bg-gray-50">
                    <Icon className="w-5 h-5 mr-3 text-[#0055b8]" />
                    <span className="text-sm font-medium text-gray-700">{item.name}</span>
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Right Column: Booking Card */}
        <div className="w-full lg:w-[400px]">
          <div className="bg-white rounded-2xl border border-gray-200 p-6 shadow-sm sticky top-24">
            <div className="mb-6">
              <span className="text-3xl font-bold text-[#0055b8]">{roomDetail.price}</span>
              <span className="text-gray-500 ml-2">{roomDetail.unit}</span>
            </div>
            
            <div className="space-y-4 mb-6">
              <div>
                <label className="block text-sm font-semibold text-gray-700 mb-1">Tanggal</label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <Calendar className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="date"
                    className="block w-full pl-10 pr-3 py-2.5 border border-gray-300 rounded-lg leading-5 bg-gray-50 text-gray-900 focus:outline-none focus:ring-1 focus:ring-[#0055b8] focus:border-[#0055b8] sm:text-sm"
                  />
                </div>
              </div>
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-1">Waktu Mulai</label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Clock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="time"
                      className="block w-full pl-10 pr-3 py-2.5 border border-gray-300 rounded-lg leading-5 bg-gray-50 text-gray-900 focus:outline-none focus:ring-1 focus:ring-[#0055b8] focus:border-[#0055b8] sm:text-sm"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-semibold text-gray-700 mb-1">Waktu Selesai</label>
                  <div className="relative">
                    <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                      <Clock className="h-5 w-5 text-gray-400" />
                    </div>
                    <input
                      type="time"
                      className="block w-full pl-10 pr-3 py-2.5 border border-gray-300 rounded-lg leading-5 bg-gray-50 text-gray-900 focus:outline-none focus:ring-1 focus:ring-[#0055b8] focus:border-[#0055b8] sm:text-sm"
                    />
                  </div>
                </div>
              </div>
            </div>

            <Button className="w-full bg-[#0055b8] hover:bg-[#004494] text-white py-3 rounded-lg font-semibold text-lg">
              Pesan Ruangan
            </Button>
            
            <p className="text-center text-xs text-gray-500 mt-4">
              Anda belum dikenakan biaya.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
