import { useState } from 'react';

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
          className={`block w-5 h-0.5 bg-white transition-all duration-300 ${
            isOpen ? 'rotate-45 translate-y-1' : ''
          }`}
        />
        <span
          className={`block w-5 h-0.5 bg-white transition-all duration-300 ${
            isOpen ? '-rotate-45 -translate-y-1' : ''
          }`}
        />
      </button>

      {isOpen && (
        <div className="fixed inset-0 z-40 bg-[#050508] flex flex-col items-center justify-center gap-8">
          {navItems.map((item) => (
            <a
              key={item.label}
              href={item.href}
              onClick={() => setIsOpen(false)}
              className="text-white text-2xl font-bold hover:text-[#6366f1] transition-colors"
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
              className="px-6 py-2 border border-[#6366f1]/40 rounded-lg text-[#6366f1] text-sm"
            >
              RESUME
            </a>
            <a
              href="#contact"
              onClick={() => setIsOpen(false)}
              className="px-6 py-2 bg-[#6366f1] rounded-lg text-white text-sm"
            >
              LET'S TALK
            </a>
          </div>
        </div>
      )}
    </div>
  );
}
