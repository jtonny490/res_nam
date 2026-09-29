import test from "node:test";
import assert from "node:assert/strict";

import { CATEGORIES, STATUSES, formatDate, isWithinPilot, severityLabel, titleCase } from "../src/utils.js";

test("titleCase normalizes separators", () => {
  assert.equal(titleCase("water"), "Water");
  assert.equal(titleCase("authority-requests"), "Authority Requests");
  assert.equal(titleCase(""), "");
});

test("severity labels are bounded", () => {
  assert.equal(severityLabel(1), "Minor");
  assert.equal(severityLabel(5), "Severe");
  assert.equal(severityLabel(9), "Unknown");
});

test("pilot coverage bounds", () => {
  assert.equal(isWithinPilot(-1.6, 31.5), true);
  assert.equal(isWithinPilot(0.7, 35.0), true);
  assert.equal(isWithinPilot(-3, 34), false);
  assert.equal(isWithinPilot(0, 40), false);
  assert.equal(isWithinPilot("abc", 34), false);
});

test("formatDate handles valid and invalid input", () => {
  assert.equal(formatDate(""), "");
  assert.equal(formatDate("not-a-date"), "");
  assert.notEqual(formatDate("2026-08-25T10:00:00Z"), "");
});

test("enum lists match the backend contract", () => {
  assert.deepEqual(CATEGORIES, ["water", "air", "waste", "deforestation", "other"]);
  assert.deepEqual(STATUSES, ["open", "investigating", "stale", "resolved"]);
});
