/* eslint-disable no-undef */
const selfsigned = require('selfsigned');
const fs = require('fs');
const path = require('path');

// Generate a self-signed certificate
const pems = selfsigned.generate([{ name: 'commonName', value: 'localhost' }], {
  keySize: 2048,
  days: 365,
  algorithm: 'sha256',
  extensions: [{ name: 'basicConstraints', cA: true }],
});

const sslDir = path.join(__dirname, '../src/config/ssl');

// Make sure ssl exsists
if (!fs.existsSync(sslDir)) {
  fs.mkdirSync(sslDir);
}

// Write the key and certificate to files
fs.writeFileSync(path.join(sslDir, 'privateKey.key'), pems.private);
fs.writeFileSync(path.join(sslDir, 'certificate.crt'), pems.cert);

console.log('SSL certificates generated successfully!');
