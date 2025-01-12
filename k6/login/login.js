import http from 'k6/http';

const url = 'http://localhost:8080/user/login_email_password';

export default function () {
    let data = {
        email: "17858644295@qq.com",
        password: "125aAss#"
    };
    const success = http.expectedStatuses(200);
    // 传入 json
    http.post(url, JSON.stringify(data), {
        headers: { 'Content-Type': 'application/json' },
        responseCallback: success,
    });
}
