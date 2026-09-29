import test from "node:test";
import assert from "node:assert/strict";

class MemoryStorage {
  constructor() {
    this.map = new Map();
  }
  getItem(key) {
    return this.map.has(key) ? this.map.get(key) : null;
  }
  setItem(key, value) {
    this.map.set(key, String(value));
  }
  removeItem(key) {
    this.map.delete(key);
  }
}

globalThis.localStorage = new MemoryStorage();

const calls = [];
globalThis.fetch = async (url, options = {}) => {
  calls.push({ url, options });
  return { status: 200, ok: true, json: async () => ({ ok: true }) };
};

const { api } = await import("../src/api.js");
const { setSession } = await import("../src/state.js");

test("login posts JSON without an auth header", async () => {
  calls.length = 0;
  await api.login({ email: "a@b.c", password: "secret" });
  assert.equal(calls[0].url, "/api/auth/login");
  assert.equal(calls[0].options.method, "POST");
  assert.equal(calls[0].options.headers.Authorization, undefined);
  assert.equal(JSON.parse(calls[0].options.body).email, "a@b.c");
});

test("report list serializes query params and adds the bearer token", async () => {
  setSession({ token: "abc", user: { id: 1, role: "public" } });
  calls.length = 0;
  await api.reports({ category: "water", page: 2, status: "" });
  assert.equal(calls[0].url, "/api/reports?category=water&page=2");
  assert.equal(calls[0].options.headers.Authorization, "Bearer abc");
});

test("createReport with a photo uses multipart form data", async () => {
  setSession({ token: "abc", user: { id: 1, role: "public" } });
  calls.length = 0;
  const file = new File(["x"], "a.png", { type: "image/png" });
  await api.createReport({ title: "Spill", category: "water" }, file);
  assert.equal(calls[0].options.headers["Content-Type"], undefined);
  assert.ok(calls[0].options.body instanceof FormData);
  assert.equal(calls[0].options.body.get("title"), "Spill");
  assert.equal(calls[0].options.body.get("photo").name, "a.png");
});

test("like and unlike use POST and DELETE on the same route", async () => {
  calls.length = 0;
  await api.like(5);
  await api.unlike(5);
  assert.equal(calls[0].options.method, "POST");
  assert.equal(calls[1].options.method, "DELETE");
  assert.equal(calls[0].url, "/api/reports/5/like");
});

test("error responses throw the server message", async () => {
  globalThis.fetch = async () => ({
    status: 403,
    ok: false,
    json: async () => ({ error: "forbidden" }),
  });
  await assert.rejects(() => api.deleteReport(1), /forbidden/);
});
