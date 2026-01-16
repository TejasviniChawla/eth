import './globals.css';
import type { Metadata } from 'next';
import Providers from './providers';

export const metadata: Metadata = {
  title: 'ether.fi Staking Dashboard',
  description: 'Track staking positions, yields, and protocol APY across ether.fi, Lido, and Rocket Pool.'
};

export default function RootLayout({
  children
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
