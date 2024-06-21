import { useEffect } from 'react';
import Link from 'next/link';

const Header = () => {
    const toggleMobileNav = () => {
        const mobileNav = document.getElementById('mobile-nav');
        const hamburger = document.querySelector('.hamburger');
        if (mobileNav && hamburger) {
            mobileNav.classList.toggle('show');
            hamburger.classList.toggle('active');
        }
    };

    useEffect(() => {
        const getCookie = (name: string) => {
            const cookieString = document.cookie;
            const cookies = cookieString.split('; ');

            for (let cookie of cookies) {
                const [cookieName, cookieValue] = cookie.split('=');
                if (cookieName === name) {
                    return decodeURIComponent(cookieValue);
                }
            }

            return '';
        };

        const username = getCookie('username');
        const accountLink = document.getElementById('accountLink');
        const mobileAccountLink = document.getElementById('mobileAccountLink');

        if (username) {
            if (accountLink) {
                accountLink.href = `/profile/${username}`;
                accountLink.textContent = username;
            }

            if (mobileAccountLink) {
                mobileAccountLink.href = `/profile/${username}`;
                mobileAccountLink.textContent = username;
            }
        }
    }, []);

    return (
        <>
            <header>
                <Link href="/">
                    <img src="/assets/images/azuritelogo.png" alt="Azurite Logo" />
                </Link>
                <nav>
                    <Link href="/">Home</Link>
                    <Link href="/games">Games</Link>
                    <Link href="/mod-manager">Mod Manager</Link>
                    <Link href="/upload">Upload</Link>
                    <Link id="accountLink" href="/login">Login/Signup</Link>
                </nav>
                <div className="hamburger" onClick={toggleMobileNav}>
                    <div></div>
                    <div></div>
                    <div></div>
                </div>
            </header>
            <div className="mobile-nav" id="mobile-nav">
                <Link href="/">Home</Link>
                <Link href="/games">Games</Link>
                <Link href="/mod-manager">Mod Manager</Link>
                <Link href="/upload">Upload</Link>
                <Link id="mobileAccountLink" href="/login">Login/Signup</Link>
            </div>
        </>
    );
};

export default Header;
