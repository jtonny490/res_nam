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
const state = await import("../src/state.js");

test("session roundtrip and role helpers", () => {
  assert.equal(state.getSession(), null);
  assert.equal(state.currentRole(), "guest");
  assert.equal(state.isAuthenticated(), false);

  state.setSession({ token: "t", user: { id: 3, role: "admin" } });
  assert.equal(state.isAuthenticated(), true);
  assert.equal(state.currentRole(), "admin");
  assert.equal(state.isAdmin(), true);
  assert.equal(state.canAccessMap(), true);
  assert.equal(state.canModerate(), true);

  state.setSession({ token: "t", user: { id: 4, role: "authority" } });
  assert.equal(state.isAuthority(), true);
  assert.equal(state.canAccessMap(), true);
  assert.equal(state.canModerate(), false);

  state.clearSession();
  assert.equal(state.isAuthenticated(), false);
});

test("ownership helpers honour admin override", () => {
  state.setSession({ token: "t", user: { id: 7, role: "public" } });
  assert.equal(state.userId(), 7);
  assert.equal(state.isOwner({ user_id: 7 }), true);
  assert.equal(state.isOwner({ user_id: 8 }), false);
  assert.equal(state.canEdit({ user_id: 8 }), false);

  state.setSession({ token: "t", user: { id: 9, role: "admin" } });
  assert.equal(state.canEdit({ user_id: 8 }), true);
});
