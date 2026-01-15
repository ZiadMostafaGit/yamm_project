




import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 100 }, // 100 users
        { duration: '1m', target: 100 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<100'], // Now we aim for <100ms
    },
};

const BASE_URL = 'http://localhost:8080/api/v1';

// This runs ONCE per user to get the token
export function setup() {
    const res = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
        email: "admin@yamm.com",
        password: "Zizoshata2003@"
    }), { headers: { 'Content-Type': 'application/json' } });
    
    return res.json().data.token;
}

export default function (token) {
    const params = {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    };

    const res = http.get(`${BASE_URL}/faq/get_for_customer?lang=AR`, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        // Check if data is an array (even if empty)
        'is array': (r) => Array.isArray(r.json().data),
    });

    sleep(0.05); // Faster requests
}