package store

const schemaSQL = `
CREATE TABLE IF NOT EXISTS chambers (id TEXT PRIMARY KEY, name TEXT NOT NULL, site_code TEXT NOT NULL, depth_m REAL NOT NULL, temperature REAL NOT NULL, created_at TEXT NOT NULL, tags TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS instruments (id TEXT PRIMARY KEY, serial TEXT NOT NULL, firmware TEXT NOT NULL, frequency_hz REAL NOT NULL, drift_ppm REAL NOT NULL, calibrated_at TEXT NOT NULL, enabled INTEGER NOT NULL, tags TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS surveys (id TEXT PRIMARY KEY, chamber_id TEXT NOT NULL, lead TEXT NOT NULL, status TEXT NOT NULL, started_at TEXT NOT NULL, closed_at TEXT NOT NULL, pulse_count INTEGER NOT NULL, echo_count INTEGER NOT NULL, open_anomaly INTEGER NOT NULL, notes TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS pulses (id TEXT PRIMARY KEY, survey_id TEXT NOT NULL, instrument_id TEXT NOT NULL, sequence INTEGER NOT NULL, emitted_at TEXT NOT NULL, duration_ms INTEGER NOT NULL, gain_db REAL NOT NULL, samples TEXT NOT NULL, tags TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS echoes (id TEXT PRIMARY KEY, survey_id TEXT NOT NULL, pulse_id TEXT NOT NULL, peak_db REAL NOT NULL, noise_db REAL NOT NULL, distance_m REAL NOT NULL, confidence REAL NOT NULL, bands TEXT NOT NULL, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS anomalies (id TEXT PRIMARY KEY, survey_id TEXT NOT NULL, echo_id TEXT NOT NULL, kind TEXT NOT NULL, severity INTEGER NOT NULL, state TEXT NOT NULL, evidence TEXT NOT NULL, created_at TEXT NOT NULL, reviewed_at TEXT NOT NULL, reviewer TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS reviews (id TEXT PRIMARY KEY, anomaly_id TEXT NOT NULL, decision TEXT NOT NULL, comment TEXT NOT NULL, reviewer TEXT NOT NULL, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS routes (id TEXT PRIMARY KEY, survey_id TEXT NOT NULL, status TEXT NOT NULL, stops TEXT NOT NULL, distance_m REAL NOT NULL, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS calibrations (id TEXT PRIMARY KEY, instrument_id TEXT NOT NULL, survey_id TEXT NOT NULL, reference_db REAL NOT NULL, measured_db REAL NOT NULL, offset_db REAL NOT NULL, operator TEXT NOT NULL, passed INTEGER NOT NULL, created_at TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS archives (id TEXT PRIMARY KEY, survey_id TEXT NOT NULL, object_key TEXT NOT NULL, digest TEXT NOT NULL, size_bytes INTEGER NOT NULL, completed_at TEXT NOT NULL, verified INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY AUTOINCREMENT, survey_id TEXT NOT NULL, kind TEXT NOT NULL, payload TEXT NOT NULL, created_at TEXT NOT NULL);
`
