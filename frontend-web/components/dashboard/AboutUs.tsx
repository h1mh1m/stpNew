import React from 'react';

const team = [
  { 
    id: 1, 
    name: 'Jane Doe', 
    role: 'Director', 
    image: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?q=80&w=300&auto=format&fit=crop' 
  },
  { 
    id: 2, 
    name: 'John Smith', 
    role: 'Lead Developer', 
    image: 'https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?q=80&w=300&auto=format&fit=crop' 
  },
];

export default function AboutUs() {
  return (
    <div className="py-12 mt-4 border-t border-gray-200">
      <h2 className="text-2xl font-bold text-gray-900 mb-8 px-2 border-l-4 border-[#512BD4]">Who's Us</h2>
      <div className="flex flex-wrap gap-8 px-2">
        {team.map((member) => (
          <div key={member.id} className="flex flex-col items-center flex-shrink-0 w-48 group">
            <div className="w-32 h-32 rounded-full overflow-hidden mb-4 border-2 border-transparent group-hover:border-[#512BD4] transition-all duration-300 shadow-md group-hover:shadow-lg transform group-hover:-translate-y-1">
              <img src={member.image} alt={member.name} className="w-full h-full object-cover" />
            </div>
            <h3 className="text-lg font-bold text-gray-900 text-center">{member.name}</h3>
            <p className="text-sm text-center text-[#512BD4] font-medium transition-colors">
              {member.role}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
