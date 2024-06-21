import type { AppProps } from 'next/app'
import '../stylesheets/root.css'
import '../stylesheets/homepage.css'
import '../stylesheets/games.css'
import '../stylesheets/download.css'

export default function MyApp({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}
