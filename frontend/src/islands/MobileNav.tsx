import { useState } from 'react';
import ThemeToggle from './ThemeToggle';

interface NavItem {
  label: string;
  href: string;
}

interface Props {
  navItems: NavItem[];
  resumePath: string;
}

export default function MobileNav({ navItems, resumePath }: Props) {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="md:hidden">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="relative z-50 w-8 h-8 flex flex-col justify-center items-center gap-1.5"
        aria-label={isOpen ? 'Close menu' : 'Open menu'}
      >
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? 'rotate-45 translate-y-1' : ''
          }`}
        />
        <span
          className={`block w-5 h-0.5 bg-text-primary transition-all duration-300 ${
            isOpen ? '-rotate-45 -translate-y-1' : ''
          }`}
        />
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-40 bg-base flex flex-col items-center justify-center gap-8">
          <div className="absolute top-5 right-14">
            <ThemeToggle />
          </div>
          {navItems.map((item) => (
            <a
              key={item.label}
              href={item.href}
              onClick={() => setIsOpen(false)}
              className="text-text-primary text-2xl font-bold hover:text-accent-indigo transition-colors"
              style={{ letterSpacing: '-0.06em' }}
            >
              {item.label}
            </a>
          ))}
          <div className="flex gap-4 mt-4">
            <a
              href={resumePath}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 border border-accent-indigo/40 rounded-lg text-accent-indigo text-sm"
            >
              RESUME
            </a>
            <a
              href="/#contact"
              onClick={() => setIsOpen(false)}
              className="px-6 py-2 bg-accent-indigo rounded-lg text-white text-sm"
            >
              LET'S TALK
            </a>
          </div>
        </div>
      )}
    </div>
  );
}
