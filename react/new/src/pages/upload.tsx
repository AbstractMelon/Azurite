import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Header from '../components/Header';
import styles from '../stylesheets/Upload.module.css';

interface Game {
  id: string;
  name: string;
}

const UploadMod: React.FC = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [formData, setFormData] = useState<FormData>(() => new FormData());
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    /*
    if (!document.cookie.includes('username')) {
      router.replace('/login');
    }
    */

    fetch('/api/games')
      .then((response) => response.json())
      .then((data) => {
        // Convert the games object to an array
        const gamesArray = Object.keys(data).map((key) => ({
          id: key,
          name: data[key].name,
        }));
        setGames(gamesArray);
      })
      .catch((error) => {
        console.error('Error fetching games:', error);
        setErrorMessage('Error fetching games');
      });
  }, [router]);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, files } = event.target;
    const updatedFormData = new FormData();

    formData.forEach((value, key) => {
      updatedFormData.append(key, value);
    });

    if (files) {
      updatedFormData.set(name, files[0]);
    } else {
      updatedFormData.set(name, value);
    }

    setFormData(updatedFormData);
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    fetch('/api/upload', {
      method: 'POST',
      body: formData,
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          console.error('Error uploading mod:', data.error);
          setErrorMessage(data.error);
        } else {
          console.log('Upload successful:', data);
        }
      })
      .catch((error) => {
        console.error('Error uploading mod:', error);
        setErrorMessage('Error uploading mod');
      });
  };

  return (
    <>
      <Head>
        <title>Azurite - Upload Mod</title>
      </Head>

      <Header />

      <div className={styles.uploadContainer}>
        <h2>Upload Mod</h2>
        {errorMessage && <p className={styles.error}>{errorMessage}</p>}
        <form id="upload-form" onSubmit={handleSubmit} encType="multipart/form-data">
          <input type="text" name="modName" placeholder="Mod Name" required onChange={handleChange} />
          <input type="text" name="modVersion" placeholder="Version" onChange={handleChange} />
          <textarea name="modDescription" placeholder="Mod Description" required onChange={handleChange} />
          <p>Mod icon:</p>
          <input type="file" name="modIcon" accept="image/*" required onChange={handleChange} />
          <p>Mod .dll:</p>
          <input type="file" name="modFile" accept=".dll" required onChange={handleChange} />
          <label htmlFor="gameSelect">Select Game:</label>
          <select id="gameSelect" name="gameId" required defaultValue="" onChange={handleChange}>
            <option value="" disabled>Select a game</option>
            {games.map((game) => (
              <option key={game.id} value={game.id}>{game.name}</option>
            ))}
          </select>
          <button type="submit">Upload</button>
        </form>
      </div>
    </>
  );
};

export default UploadMod;
