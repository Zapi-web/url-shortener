import { check } from 'k6';
import http from 'k6/http'

const BASE_URL = 'http://localhost:8081'

export const options = {
    vus: 1000,
    iterations: 10000
};

export default function() {
    const payload = JSON.stringify({
        url: "https://google.com",
        user_id: 1 
    });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    }
    const postRes = http.post(`${BASE_URL}/api/v1/`, payload, params)

    check(postRes, {
        'POST link created (200)': (r) => r.status === 200
    })

    const shortURL = postRes.json('short_url')

    if (shortURL) {
        const fullShortURL = `${BASE_URL}/${shortURL}`;
        const getParams = { redirects: 0 };
        const getRes = http.get(fullShortURL, getParams);

        check(getRes, {
            'GET short link received (302)': (r) => r.status === 302,
            'Redirect location is Google': (r) => r.headers['Location'] === 'https://google.com' || r.headers['location'] === 'https://google.com',
        });
    }
}