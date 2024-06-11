const Greenlock = require('greenlock');
const GreenlockStore = require('greenlock-store-fs');
const AcmeHttp01Standalone = require('acme-http-01-standalone');

const greenlock = Greenlock.create({
    packageRoot: __dirname,
    configDir: './greenlock.d',
    maintainerEmail: 'your-email@example.com',
    cluster: false,
    store: GreenlockStore.create({ configDir: './greenlock.d' })
});

greenlock.manager.defaults({
    agreeToTerms: true,
    subscriberEmail: 'your-email@example.com',
    challenges: {
        'http-01': AcmeHttp01Standalone
    }
}).then(() => {
    return greenlock.add({
        subject: 'your-domain.com',
        altnames: ['your-domain.com']
    });
}).then(() => {
    console.log('Success! Certificates obtained and stored.');
}).catch((err) => {
    console.error('Error:', err);
});
