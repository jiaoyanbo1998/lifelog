import http from 'k6/http';

const url = 'http://localhost:8080/lifeLog/detail';
const jwtToken = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MzMsIk5pY2tOYW1lIjoiMjRiMDYzZTAtNjU5OS00MDI4LTk3YzctZTgwY2M5NDI4ZGU4IiwiZXhwIjoxNzM4NjU0OTE1LCJTZXNzaW9uSWQiOiI1OGY4NGQ0MC1hZTU0LTRlMWMtOTlkNi03OWQyNzgzYjBiMDAiLCJBdXRob3JpdHkiOjJ9.sFFmBqNXBpl7JN6WDJgpF_YhcWp6lX1aFqbMOLqJ0j12cZIzVYljRh-Xt9rpXG6k69l-SrorGBbPpghw4AyZoQ'

export default function () {
    let data = {
        id: 28,
        public: false
    };
    const success = http.expectedStatuses(200);

    // 传入 json
    http.post(url, JSON.stringify(data), {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwtToken}` // 添加JWT Token到请求头
        },
        responseCallback: success,
    });
}