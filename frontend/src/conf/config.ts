export const config = {
    API_URL: process.env?.VUE_APP_API_URL || 'http://localhost:7080',
    SHORTEN_PATH: process.env?.VUE_APP_SHORTEN_PATH || 'u',
    AUTH_PATH: process.env?.VUE_APP_AUTH_PATH || 'auth'
};
