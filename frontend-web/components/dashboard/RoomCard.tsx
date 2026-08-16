import React from 'react';
import { MapPin, Users } from 'lucide-react';

interface RoomCardProps {
  title: string;
  department: string;
  location: string;
  capacity: number;
  price: string;
  unit: string;
  tag: string;
}

export const RoomCard = ({ title, department, location, capacity, price, unit, tag }: RoomCardProps) => {
  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden shadow-sm hover:shadow-md transition-shadow">
      {/* Blank Purple Block instead of Image */}
      <div className="h-48 w-full bg-[#512BD4]"></div>
      
      <div className="p-5">
        <div className="inline-block bg-[#007bff] text-white text-xs font-semibold px-3 py-1 rounded-md mb-4">
          {tag}
        </div>
        <h3 className="font-bold text-gray-900 text-lg mb-1">{title}</h3>
        <p className="text-sm text-gray-500 mb-5">{department}</p>
        
        <div className="flex items-center space-x-6 text-xs text-gray-700 mb-5">
          <div className="flex items-center">
            <MapPin className="w-3.5 h-3.5 mr-1.5 text-[#007bff]" />
            {location}
          </div>
          <div className="flex items-center">
            <Users className="w-3.5 h-3.5 mr-1.5 text-[#007bff]" />
            {capacity} Orang
          </div>
        </div>
        
        <div className="text-sm mt-auto pt-2 border-t border-gray-100">
          <span className="font-bold text-[#007bff] text-lg">{price}</span>
          <span className="text-gray-600 ml-1">{unit}</span>
        </div>
      </div>
    </div>
  );
};
