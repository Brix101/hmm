const hostname = window.location.hostname;

export const BACKEND = `http://${hostname}:5000`;
export const BASE_URL = BACKEND + "/api";
export const STATIC_URL = BACKEND + "/data/files";
