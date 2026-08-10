'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Footer } from '@/components/layout/Footer';

export default function SignUpPage() {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="min-h-screen flex flex-col bg-[#484555]">
      <main className="flex-grow flex flex-col items-center justify-center p-6 py-12">
        <div className="flex flex-col items-center text-center mb-10 text-white">
          <div className="bg-[#EBE2FF] rounded-full p-4 mb-4 text-[#512BD4]">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 4h4v4H4V4zm6 0h4v4h-4V4zm6 0h4v4h-4V4zM4 10h4v4H4v-4zm6 0h4v4h-4v-4zm6 0h4v4h-4v-4zM4 16h4v4H4v-4zm6 0h4v4h-4v-4zm6 0h4v4h-4v-4z" fill="currentColor" fillOpacity="0.5"/>
              <path d="M4 4h4v4H4V4zm12 12h4v4h-4v-4z" fill="currentColor"/>
            </svg>
          </div>
          <h1 className="text-4xl font-bold mb-3">Welcome to TrackIQ Access Portal</h1>
          <p className="text-gray-300">Please sign up to access circuit features and data.</p>
        </div>

        <div className="w-full max-w-[960px] bg-white rounded-lg shadow-xl overflow-hidden flex flex-col md:flex-row">
          
          {/* Left Column - ITS Track Team */}
          <div className="w-full md:w-1/2 bg-[#F6FAFE] p-10 md:p-14 flex flex-col justify-center border-b md:border-b-0 md:border-r border-gray-100">
            <div className="flex items-center gap-3 mb-6">
              <div className="w-2 h-8 bg-[#512BD4] rounded-full"></div>
              <h2 className="text-2xl font-bold text-[#171C1F]">ITS Track Team</h2>
            </div>
            <p className="text-[#171C1F] mb-10">Use your authorized organization SSO credentials.</p>
            
            <Button className="w-full h-12" leftIcon={
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                <line x1="16" y1="2" x2="16" y2="6"/>
                <line x1="8" y1="2" x2="8" y2="6"/>
                <line x1="3" y1="10" x2="21" y2="10"/>
              </svg>
            }>
              Sign up with Team SSO
            </Button>
          </div>

          {/* Right Column - External Clients */}
          <div className="w-full md:w-1/2 bg-white p-10 md:p-14 flex flex-col justify-center">
            <h2 className="text-xl font-bold text-[#171C1F] mb-8">Register External Clients & Drivers</h2>
            
            <div className="space-y-5 mb-6">
              <Input 
                label="Email Address" 
                type="email" 
                placeholder="driver@example.com" 
                leftIcon={
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                    <polyline points="22,6 12,13 2,6"></polyline>
                  </svg>
                }
              />
              
              <div>
                <Input 
                  label="Password"
                  type={showPassword ? "text" : "password"} 
                  placeholder="••••••••" 
                  leftIcon={
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                      <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                    </svg>
                  }
                  rightIcon={
                    <div onClick={() => setShowPassword(!showPassword)} className="cursor-pointer hover:text-black text-gray-400">
                      {showPassword ? (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                          <circle cx="12" cy="12" r="3"></circle>
                          <line x1="3" y1="3" x2="21" y2="21"></line>
                        </svg>
                      ) : (
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                          <circle cx="12" cy="12" r="3"></circle>
                        </svg>
                      )}
                    </div>
                  }
                />
              </div>
            </div>

            <Link href="#" className="text-sm font-semibold text-[#512BD4] hover:underline mb-8 self-start">
              Forgot Password?
            </Link>

            <Button className="w-full h-12 mb-6">
              Sign Up
            </Button>
            
            <p className="text-center text-sm text-[#171C1F]">
              Already have an account? <Link href="/signin" className="font-semibold text-[#512BD4] hover:underline">Sign in</Link>
            </p>
          </div>
        </div>
      </main>
      <Footer />
    </div>
  );
}
