'use client';

import React, { useEffect, useState } from 'react';
import { MapPin, ExternalLink } from 'lucide-react';

interface RoomItem {
  id?: string;
  ID?: string;
  name?: string;
  Name?: string;
  description?: string;
  Description?: string;
  capacity?: number;
  Capacity?: number;
  price?: number;
  Price?: number;
  image_url?: string;
  ImageURL?: string;
  features?: string;
  Features?: string;
}

export default function Dashboard() {
  const [rooms, setRooms] = useState<RoomItem[]>([]);

  useEffect(() => {
    const fetchRooms = async () => {
      try {
        const token = localStorage.getItem('token');
        const headers: Record<string, string> = {};
        if (token) {
          headers['Authorization'] = `Bearer ${token}`;
        }

        const res = await fetch('/api/Beranda', { headers });
        if (res.ok) {
          const data = await res.json();
          if (data.data && Array.isArray(data.data)) {
            setRooms(data.data);
          }
        }
      } catch (err) {
        console.error('Error fetching rooms:', err);
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
          {rooms.length > 0 ? rooms.map((room, idx) => {
            const roomId = room.id || room.ID || `room-${idx}`;
            const roomName = room.name || room.Name || 'Ruangan';
            const roomDesc = room.description || room.Description || '';
            const roomPrice = room.price ?? room.Price ?? 0;
            const roomImage = room.image_url || room.ImageURL || "https://images.unsplash.com/photo-1580274455191-1c62238fa333?auto=format&fit=crop&q=80&w=1000";
            const roomFeatures = room.features || room.Features || 'Jangka Pendek';

            return (
              <div key={roomId} className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden hover:shadow-md transition-shadow">
                <div className="h-64 bg-gray-200 relative">
                  <img 
                    src={roomImage} 
                    alt={roomName} 
                    className="w-full h-full object-cover"
                  />
                </div>
                <div className="p-8">
                  <div className="inline-block px-4 py-1.5 bg-[#007bff] text-white text-sm font-medium rounded-full mb-4">
                    {roomFeatures}
                  </div>
                  <h3 className="text-3xl font-bold text-gray-900 mb-1">{roomName}</h3>
                  <p className="text-gray-500 mb-6 text-xl">{roomDesc}</p>
                  
                  <div className="flex items-start justify-between gap-6 mb-6">
                    <div className="flex items-start gap-3">
                      <MapPin className="w-6 h-6 text-[#007bff] shrink-0 mt-0.5" />
                      <p className="text-gray-600 text-base leading-relaxed max-w-sm">
                        Lokasi: {roomDesc}
                      </p>
                    </div>
                    <a href="#" className="flex items-center gap-1 text-[#007bff] text-sm font-medium hover:underline shrink-0">
                      <ExternalLink className="w-4 h-4" />
                      Lihat di Maps
                    </a>
                  </div>
                  
                  <div className="pt-6 border-t border-gray-100">
                    <p className="text-4xl font-bold text-[#007bff] tracking-tight">
                      Rp {roomPrice.toLocaleString('id-ID')}
                    </p>
                    <p className="text-gray-800 text-xl font-bold mt-2">per Jam</p>
                  </div>
                </div>
              </div>
            );
          }) : (
            <p className="text-gray-500">Memuat data sirkuit...</p>
          )}
        </div>
      </div>
    </div>
  );
}
