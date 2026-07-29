import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// 自定义指标
const errorRate = new Rate('errors');
const latency = new Trend('api_latency');

// 测试配置
export const options = {
  stages: [
    { duration: '30s', target: 20 },   // 预热：30秒内爬升到20个虚拟用户
    { duration: '1m', target: 50 },    // 负载：保持50个用户1分钟
    { duration: '30s', target: 100 },  // 压力：爬升到100个用户
    { duration: '1m', target: 100 },   // 峰值：保持100个用户1分钟
    { duration: '30s', target: 0 },    // 释放：降到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],  // 95% 请求 < 500ms
    errors: ['rate<0.01'],             // 错误率 < 1%
  },
};

const BASE_URL = __ENV.TARGET_URL || 'http://localhost:8080';

export default function () {
  // 测试 1: 健康检查
  const healthRes = http.get(`${BASE_URL}/healthz`);
  check(healthRes, {
    'health check status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.5);

  // 测试 2: 获取公司列表
  const companiesRes = http.get(`${BASE_URL}/api/v1/companies`);
  latency.add(companiesRes.timings.duration);
  check(companiesRes, {
    'companies status 200': (r) => r.status === 200,
    'companies has data': (r) => JSON.parse(r.body).data !== undefined,
  }) || errorRate.add(1);

  sleep(1);

  // 测试 3: 获取 SME 公司列表
  const smeRes = http.get(`${BASE_URL}/api/v1/sme-companies`);
  check(smeRes, {
    'sme companies status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(1);

  // 测试 4: SRE 指标端点
  const metricsRes = http.get(`${BASE_URL}/api/v1/sre/metrics`);
  check(metricsRes, {
    'metrics status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  sleep(0.5);
}

// 测试完成后输出摘要
export function handleSummary(data) {
  return {
    'stdout': JSON.stringify(data, null, 2),
    'load-test-results.json': JSON.stringify(data, null, 2),
  };
}