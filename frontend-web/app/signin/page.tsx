'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Footer } from '@/components/layout/Footer';

export default function SignInPage() {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="min-h-screen flex flex-col bg-[#F6FAFE]">
      <main className="flex-grow flex items-center justify-center p-6">
        <div className="w-full max-w-[480px] bg-white rounded-lg p-10 shadow-[0_4px_24px_rgba(0,0,0,0.04)] text-center">
          
          <div className="flex justify-center mb-6 text-[#512BD4]">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M4 4h4v4H4V4zm6 0h4v4h-4V4zm6 0h4v4h-4V4zM4 10h4v4H4v-4zm6 0h4v4h-4v-4zm6 0h4v4h-4v-4zM4 16h4v4H4v-4zm6 0h4v4h-4v-4zm6 0h4v4h-4v-4z" fill="currentColor" fillOpacity="0.5"/>
              <path d="M4 4h4v4H4V4zm12 12h4v4h-4v-4z" fill="currentColor"/>
            </svg>
            <span className="text-2xl font-bold ml-2 text-black">TrackIQ</span>
          </div>

          <h1 className="text-3xl font-bold text-[#171C1F] mb-2">Welcome Back</h1>
          <p className="text-gray-500 mb-8">Sign in to your TrackIQ account</p>

          <Button 
            variant="secondary" 
            className="w-full mb-6 font-medium text-black bg-[#F0F4F8] hover:bg-[#E2E8F0]"
            leftIcon={
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="text-[#512BD4]">
                <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2"/>
                <path d="M12 16a4 4 0 1 0 0-8 4 4 0 0 0 0 8z" fill="currentColor"/>
              </svg>
            }
          >
            Sign In with myITS (SSO)
          </Button>

          <div className="relative mb-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-white text-gray-500">or sign in with email</span>
            </div>
          </div>

          <div className="text-left space-y-4 mb-8">
            <Input 
              label="Email Address" 
              type="email" 
              placeholder="Enter your email" 
            />
            
            <div>
              <div className="flex justify-between items-center mb-1">
                <label className="text-sm font-semibold text-[#171C1F]">Password</label>
                <Link href="#" className="text-sm font-semibold text-[#512BD4] hover:underline">
                  Forgot Password?
                </Link>
              </div>
              <Input 
                type={showPassword ? "text" : "password"} 
                placeholder="Enter your password" 
                rightIcon={
                  <div onClick={() => setShowPassword(!showPassword)} className="cursor-pointer hover:text-black">
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

          <Button className="w-full h-12 text-[15px]">
            Sign In
          </Button>
        </div>
      </main>
      <Footer />
    </div>
  );
}
