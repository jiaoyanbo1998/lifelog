import http from 'k6/http';
import { FormData } from 'https://jslib.k6.io/formdata/0.0.2/index.js';

const url = 'http://localhost:8080/files/upload';
const jwtToken = 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJZCI6MzMsIk5pY2tOYW1lIjoiMjRiMDYzZTAtNjU5OS00MDI4LTk3YzctZTgwY2M5NDI4ZGU4IiwiZXhwIjoxNzM4MDU1NzQ3LCJTZXNzaW9uSWQiOiI5ODEwMjMwYy0zNzdjLTQ3ZTAtYjA0My00MjY3ODRkODYyOTAiLCJBdXRob3JpdHkiOjJ9.ll-ZnMneWJjWNgzYG-KBo9x7H7k36c8eg21yddVg5ARGyhcbw5sy3AQYcfzwNLGZ-6zLSp80CkX6VAKZxqYUuA';

// 将文件放在脚本同目录下使用相对路径
const filePath = './开题报告.pdf';
const fileName = '开题报告.pdf';
const mimeType = 'application/pdf';

// 验证文件存在性
let fileContent;
try {
    fileContent = open(filePath, 'b');
} catch (e) {
    console.error(`无法读取文件: ${e.message}`);
    throw e;
}

export default function () {
    const formData = new FormData();
    formData.append('file', http.file(fileName, fileContent, mimeType));

    const headers = {
        'Authorization': `Bearer ${jwtToken}`,
        ...formData.headers,
    };

    const response = http.post(url, formData.body(), { headers });

    if (response.status !== 200) {
        console.error(`上传失败 (${response.status}): ${response.body}`);
    } else {
        console.log('上传成功');
    }
}