import Head from "next/head";
import Link from "next/link";
import Header from "../components/Header";

const Test = () => {
    return (
        <>
            <Head>
                <title>Azurite - Pagename</title>
                <link
                    rel="icon"
                    type="image/x-icon"
                    href="/assets/images/icon.ico"
                />
            </Head>
            <Header />
        </>
    );
};

export default Test;
