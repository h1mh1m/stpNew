import React, { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary';
  leftIcon?: ReactNode;
  isLoading?: boolean;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ children, variant = 'primary', leftIcon, isLoading = false, disabled, className = '', ...props }, ref) => {
    const baseStyles = 'inline-flex items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-60 disabled:cursor-not-allowed';
    
    const variants = {
      primary: 'bg-[#512BD4] text-white hover:bg-[#4323B0] focus:ring-[#512BD4]',
      secondary: 'bg-[#F0F4F8] text-[#171C1F] hover:bg-[#E2E8F0] focus:ring-[#F0F4F8]',
    };

    return (
      <button
        ref={ref}
        disabled={disabled || isLoading}
        className={`${baseStyles} ${variants[variant]} ${className}`}
        {...props}
      >
        {isLoading ? (
          <span className="inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
        ) : (
          leftIcon && <span className="flex items-center">{leftIcon}</span>
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';

