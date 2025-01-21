import http from 'k6/http';

import { FormData } from 'https://jslib.k6.io/formdata/0.0.2/index.js';

const url = 'http://localhost:8080/files/upload';
const jwtToken = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MzMsIk5pY2tOYW1lIjoiMjRiMDYzZTAtNjU5OS00MDI4LTk3YzctZTgwY2M5NDI4ZGU4IiwiZXhwIjoxNzM4MDU1NzQ3LCJTZXNzaW9uSWQiOiI5ODEwMjMwYy0zNzdjLTQ3ZTAtYjA0My00MjY3ODRkODYyOTAiLCJBdXRob3JpdHkiOjJ9.ll-ZnMneWJjWNgzYG-KBo9x7H7k36c8eg21yddVg5ARGyhcbw5sy3AQYcfzwNLGZ-6zLSp80CkX6VAKZxqYUuA'
// 使用 open 函数加载文件内容
const filePath = 'D:/GoLangProjects/minio/开题报告1.pdf'; // 文件路径相对于脚本的位置
const fileContent = open(filePath, 'b'); // 'b' 表示以二进制模式读取文件
const fileName = '开题报告.pdf';
const mimeType = 'application/pdf';

export default function () {
    const formData = new FormData();
    formData.append('name', fileName);
    formData.append('content', fileContent);
    formData.append('file', http.file(fileName, fileContent, mimeType));

    const headers = {
        'Authorization': `Bearer ${jwtToken}`,
        ...formData.headers,
    };

    const response = http.post(url, formData.body(), { headers: headers });

    if (response.status !== 200) {
        console.error(`Upload failed with status code: ${response.status}`);
    } else {
        console.log('Upload successful');
    }
}