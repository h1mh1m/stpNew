import React, { InputHTMLAttributes, ReactNode } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  error?: string;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, leftIcon, rightIcon, error, className = '', ...props }, ref) => {
    return (
      <div className="flex flex-col gap-1 w-full">
        {label && (
          <label className="text-sm font-semibold text-[#171C1F]">
            {label}
          </label>
        )}
        <div className="relative flex items-center">
          {leftIcon && (
            <div className="absolute left-3 text-gray-500">
              {leftIcon}
            </div>
          )}
          <input
            ref={ref}
            className={`w-full rounded-md border ${
              error ? 'border-red-500' : 'border-gray-300'
            } bg-white text-[#171C1F] px-3 py-2 text-sm placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-[#512BD4] focus:border-[#512BD4] transition-colors ${
              leftIcon ? 'pl-10' : ''
            } ${rightIcon ? 'pr-10' : ''} ${className}`}
            {...props}
          />
          {rightIcon && (
            <div className="absolute right-3 text-gray-500">
              {rightIcon}
            </div>
          )}
        </div>
        {error && <span className="text-xs text-red-500">{error}</span>}
      </div>
    );
  }
);

Input.displayName = 'Input';
