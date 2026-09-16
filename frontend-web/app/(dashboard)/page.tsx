'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { MapPin, ExternalLink } from 'lucide-react';

export default function Dashboard() {
  const [rooms, setRooms] = useState<any[]>([]);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) return;

        const res = await fetch('http://localhost:3000/api/Beranda', {
          headers: { 'Authorization': `Bearer ${token}` }
        });
        const data = await res.json();
        setRooms(data.data || []);
      } catch (err) {
        console.error(err);
      }
    };
    fetchRooms();
  }, []);

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
        <h2 className="text-3xl font-bold text-gray-900 mb-3 flex items-center gap-3">
          <span>🏁</span> Temukan Sirkuit Terbaik untuk Kebutuhan Anda
        </h2>
        <p className="text-gray-600 mb-10 text-lg">
          Pilih sirkuit yang sesuai dengan kebutuhan Anda, lengkap dengan informasi fasilitas, lokasi, dan harga sewa
        </p>
        
        {/* Circuit Cards Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {rooms.length > 0 ? rooms.map((room) => (
            <div key={room.ID} className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden hover:shadow-md transition-shadow">
              <div className="h-64 bg-gray-200 relative">
                <img 
                  src={room.ImageURL || "https://images.unsplash.com/photo-1580274455191-1c62238fa333?auto=format&fit=crop&q=80&w=1000"} 
                  alt={room.Name} 
                  className="w-full h-full object-cover"
                />
              </div>
              <div className="p-8">
                <div className="inline-block px-4 py-1.5 bg-[#007bff] text-white text-sm font-medium rounded-full mb-4">
                  {room.Features || 'Jangka Pendek'}
                </div>
                <h3 className="text-3xl font-bold text-gray-900 mb-1">{room.Name}</h3>
                <p className="text-gray-500 mb-6 text-xl">{room.Description}</p>
                
                <div className="flex items-start justify-between gap-6 mb-6">
                  <div className="flex items-start gap-3">
                    <MapPin className="w-6 h-6 text-[#007bff] shrink-0 mt-0.5" />
                    <p className="text-gray-600 text-base leading-relaxed max-w-sm">
                      Lokasi: {room.Description}
                    </p>
                  </div>
                  <a href="#" className="flex items-center gap-1 text-[#007bff] text-sm font-medium hover:underline shrink-0">
                    <ExternalLink className="w-4 h-4" />
                    Lihat di Maps
                  </a>
                </div>
                
                <div className="pt-6 border-t border-gray-100">
                  <p className="text-4xl font-bold text-[#007bff] tracking-tight">
                    Rp {room.Price ? room.Price.toLocaleString('id-ID') : '0'}
                  </p>
                  <p className="text-gray-800 text-xl font-bold mt-2">per Jam</p>
                </div>
              </div>
            </div>
          )) : (
            <p className="text-gray-500">Memuat data sirkuit...</p>
          )}
        </div>
      </div>
    </div>
  );
}
