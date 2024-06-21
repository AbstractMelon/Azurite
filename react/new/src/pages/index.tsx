import Head from 'next/head';
import Link from 'next/link';
import Header from '../components/Header'

const Home = () => {
    return (
        <>
            <Head>
                <meta charSet="UTF-8" />
                <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                <title>Azurite - Homepage</title>
                <link rel="icon" type="image/x-icon" href="/assets/images/icon.ico" />
                <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.4/css/all.min.css" />
            </Head>

            <Header />
            <div style={{ paddingTop: '5px' }}>
                <section className="welcome-section">
                    <div className="container">
                        <h2>Welcome to Azurite</h2>
                        <p>Your ultimate destination for discovering and sharing amazing mods!</p>
                    </div>
                </section>

                <section className="features-section">
                    <div className="container">
                        <div className="feature">
                            <i className="fas fa-gamepad"></i>
                            <h3>Explore Mods</h3>
                            <p>Discover a vast collection of mods for your favorite games.</p>
                        </div>
                        <div className="feature">
                            <i className="fas fa-users"></i>
                            <h3>Engage with Community</h3>
                            <p>Connect with fellow modders, share tips, and showcase your creations.</p>
                        </div>
                        <div className="feature">
                            <i className="fas fa-shield-alt"></i>
                            <h3>Secure Downloads</h3>
                            <p>Download mods safely with our secure platform and verified sources.</p>
                        </div>
                    </div>
                </section>

                <section className="get-started-section">
                    <div className="container">
                        <h2>Ready to get started?</h2>
                        <Link href="/register">
                            Create an Account
                        </Link>
                    </div>
                </section>

                <footer>
                    <div className="footer-section">
                        <h3>About Us</h3>
                        <ul>
                            <li><Link href="/about">Our Story</Link></li>
                            <li><Link href="/team">Team</Link></li>
                            <li><Link href="/careers">Careers</Link></li>
                        </ul>
                    </div>
                    <div className="footer-section">
                        <h3>Documentation</h3>
                        <ul>
                            <li><Link href="/docs/uploading-mods">Uploading Mods</Link></li>
                            <li><Link href="/docs/api">API Documentation</Link></li>
                            <li><Link href="/docs/mod-guidelines">Mod Guidelines</Link></li>
                        </ul>
                    </div>
                    <div className="footer-section">
                        <h3>Support</h3>
                        <ul>
                            <li><Link href="/support/discord">Join Discord</Link></li>
                            <li><Link href="/support/contact">Contact Support</Link></li>
                            <li><Link href="/support/faq">FAQ</Link></li>
                        </ul>
                    </div>
                    <div className="footer-section">
                        <h3>Helpful Links</h3>
                        <ul>
                            <li><Link href="/support/discord">Join Discord</Link></li>
                            <li><Link href="/support/contact">Contact Support</Link></li>
                            <li><Link href="/support/faq">FAQ</Link></li>
                        </ul>
                    </div>
                    <div className="footer-section">
                        <h3>Newsletter</h3>
                        <div className="newsletter">
                            <input type="email" placeholder="Enter your email" />
                            <button type="button">Subscribe</button>
                        </div>
                    </div>
                    <div className="footer-bottom">
                        &copy; 2024 Azurite. All rights reserved. | <Link href="/privacy">Privacy Policy</Link> | <Link href="/terms">Terms of Service</Link>
                    </div>
                </footer>
            </div>
        </>
    );
}

export default Home;
