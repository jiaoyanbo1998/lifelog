import http from 'k6/http';

const url = 'http://localhost:8080/collect/detail';
const jwtToken = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MzMsIk5pY2tOYW1lIjoiMjRiMDYzZTAtNjU5OS00MDI4LTk3YzctZTgwY2M5NDI4ZGU4IiwiZXhwIjoxNzM3NjE4NDUzLCJTZXNzaW9uSWQiOiIxNTk1MGI0YS1hZDVmLTQ2ZDgtYTI4Yi0xODA3ZTExNTAzNmYiLCJBdXRob3JpdHkiOjJ9.E-6eTQsqCO_SiB0B8nWa0Cb7Rcy3v08HYe9W_7jVQ-190Z69b0NGnv1Pqw2CLezH4191SFkr_dY5wCDdNqJLzg'; // 替换为你的JWT Token

export default function () {
    let data = {
        collect_id: 8,
        limit: 10,
        offset: 0
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