import { useEffect, useState } from "react";
import Link from "next/link";
import { parseCookies } from "nookies";
import { jwtDecode } from 'jwt-decode';

interface DecodedToken {
    id: string;
    username: string;
}

const Header = () => {
    const [username, setUsername] = useState<string | null>(null);

    const toggleMobileNav = () => {
        const mobileNav = document.getElementById("mobile-nav");
        const hamburger = document.querySelector(".hamburger");
        if (mobileNav && hamburger) {
            mobileNav.classList.toggle("show");
            hamburger.classList.toggle("active");
        }
    };

    useEffect(() => {
        const cookies = parseCookies();
        console.log("Cookies:", cookies);

        const token = cookies.authToken;
        console.log("Token:", token);

        if (token) {
            try {
                const decoded: DecodedToken = jwtDecode(token);
                console.log("Decoded Token:", decoded);
                setUsername(decoded.username);
            } catch (error) {
                console.error("Failed to decode token:", error);
            }
        }
    }, []);

    return (
        <>
            <header>
                <Link href="/">
                    <img
                        src="/assets/images/azuritelogo.png"
                        alt="Azurite Logo"
                    />
                </Link>
                <nav>
                    <Link href="/">Home</Link>
                    <Link href="/games">Games</Link>
                    <Link href="/mod-manager">Mod Manager</Link>
                    <Link href="/upload">Upload</Link>
                    {username ? (
                        <Link id="accountLink" href={`/profile/${username}`}>
                            {username}
                        </Link>
                    ) : (
                        <Link id="accountLink" href="/login">
                            Login/Signup
                        </Link>
                    )}
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
                {username ? (
                    <Link id="mobileAccountLink" href={`/profile/${username}`}>
                        {username}
                    </Link>
                ) : (
                    <Link id="mobileAccountLink" href="/login">
                        Login/Signup
                    </Link>
                )}
            </div>
        </>
    );
};

export default Header;
