// load-test-single-wallet.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const WALLET_ID = '550e8400-e29b-41d4-a716-446655440000'; // фиксированный ID — вся нагрузка на ОДИН кошелёк
const BASE_URL = 'http://server:8080';

export const options = {
    scenarios: {
        single_wallet_load: {
            executor: 'constant-arrival-rate',
            rate: 1000,              // 1000 итераций в секунду
            timeUnit: '1s',
            duration: '30s',
            preAllocatedVUs: 200,    // сколько виртуальных пользователей держать наготове
            maxVUs: 1000,
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.01'],   // допускаем <1% ошибок (ожидаемые 422 не считаем ошибкой ниже)
        http_req_duration: ['p(95)<500'], // 95% запросов быстрее 500ms
    },
};

export default function () {
    const payload = JSON.stringify({
        amount: 1,
        valletId: WALLET_ID,
        operationType: 'WITHDRAW',
    });

    const params = {
        headers: { 'Content-Type': 'application/json' },
    };

    const res = http.post(`${BASE_URL}/api/v1/wallet`, payload, params);

    check(res, {
        'status is 200 or 422 (insufficient funds)': (r) => r.status === 200 || r.status === 422,
        'no 500 errors': (r) => r.status !== 500,
    });
}