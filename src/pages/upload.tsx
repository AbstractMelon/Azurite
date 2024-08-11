import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import Head from 'next/head';
import Header from '../components/Header';
import styles from '../stylesheets/upload.module.css';

interface Game {
  id: string;
  name: string;
}

const UploadMod: React.FC = () => {
  const [games, setGames] = useState<Game[]>([]);
  const [formData, setFormData] = useState<FormData>(() => new FormData());
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [uploading, setUploading] = useState<boolean>(false);
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
  }, []);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, files } = event.target as HTMLInputElement;
    const updatedFormData = new FormData();
    formData.forEach((value, key) => {
      updatedFormData.append(key, value);
    });

    if (files && files.length > 0) {
      updatedFormData.set(name, files[0]);
    } else {
      updatedFormData.set(name, value);
    }

    setFormData(updatedFormData);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setUploading(true);
    setErrorMessage(null);
    setSuccessMessage(null);

    try {
      const response = await fetch('/api/upload', {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();
      setUploading(false);

      if (response.ok) {
        setSuccessMessage('Mod uploaded successfully');
        setTimeout(() => {
          router.push(`/mods/${data.manifestUrl}`);
        }, 2000);
      } else {
        setErrorMessage(data.error || 'Error uploading mod');
      }
    } catch (error) {
      console.error('Error uploading mod:', error);
      setErrorMessage('Error uploading mod');
      setUploading(false);
    }
  };

  return (
    <>
      <Head>
        <title>Azurite - Upload Mod</title>
        <link rel="icon" type="image/x-icon" href="/assets/images/icon.ico" />
      </Head>

      <Header />

      <div className={styles.uploadContainer}>
        <h2>Upload Mod</h2>
        {errorMessage && <p className={styles.error}>{errorMessage}</p>}
        {successMessage && <p className={styles.success}>{successMessage}</p>}
        <form id="upload-form" onSubmit={handleSubmit} encType="multipart/form-data" className={styles.uploadForm}>
          <input type="text" name="modName" placeholder="Mod Name" required onChange={handleChange} />
          <input type="text" name="modVersion" placeholder="Version" onChange={handleChange} />
          <textarea name="modDescription" placeholder="Mod Description" required onChange={handleChange} />
          <input type="text" name="author" placeholder="Author" required onChange={handleChange} />
          <textarea name="teamMembers" placeholder="Team Members (comma separated)" onChange={handleChange} />
          <p>Mod icon:</p>
          <input type="file" name="modIcon" accept="image/*" required onChange={handleChange} />
          <p>Mod .dll:</p>
          <input type="file" name="modFile" accept=".dll" required onChange={handleChange} />
          <p>Optional Screenshots:</p>
          <input type="file" name="screenshots" accept="image/*" multiple onChange={handleChange} />
          <label htmlFor="gameSelect">Select Game:</label>
          <select id="gameSelect" name="gameId" required defaultValue="" onChange={handleChange}>
            <option value="" disabled>Select a game</option>
            {games.map((game) => (
              <option key={game.id} value={game.id}>{game.name}</option>
            ))}
          </select>
          <button type="submit" disabled={uploading}>
            {uploading ? 'Uploading...' : 'Upload'}
          </button>
        </form>
      </div>
    </>
  );
};

export default UploadMod;
