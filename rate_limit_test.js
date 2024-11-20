import http from "k6/http";
import { check } from "k6";
import { Rate } from "k6/metrics";

let rateLimitExceeded = new Rate("rate_limit_exceeded");
let successRate = new Rate("success_rate");

export let options = {
  scenarios: {
    high_rps: {
      executor: "constant-arrival-rate",
      rate: 2000,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 500,
      maxVUs: 4000,
    },
  },
  thresholds: {
    rate_limit_exceeded: ["rate>0.1"],
    success_rate: ["rate<0.9"],
  },
};

export default function () {
  let res = http.get("http://localhost:8080/resolve?domain=youtube.com");

  check(res, {
    "status is 200": (r) => r.status === 200,
    "status is 429": (r) => r.status === 429,
  });

  if (res.status === 429) {
    rateLimitExceeded.add(1);
  } else if (res.status === 200) {
    successRate.add(1);
  }
}
