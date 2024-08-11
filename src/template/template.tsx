import Head from 'next/head';
import Link from 'next/link';
import Header from '../components/Header'

const PageName = () => {
    return (
        <>
            <Head>
                <title>Azurite - Pagename</title>
            </Head>
            <Header />
        </>
    );
}

export default PageName;
