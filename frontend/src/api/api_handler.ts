import { buildUrl } from '@/util';
import { config } from '@/conf';
import { ref } from 'vue';
import { post } from './http';

export interface CreateLinkBody {
    customShortText?: string
}

export interface HttpError {
    status: number,
    error?: string,
    message?: string
}

export type ShortenResponse = {
    shortUrl: string,
    // TODO: hint: customShortTextAlreadyUsed (Needs backend)
}

interface AuthSuccessResponse {
    status: '200' | 200,
    jwt: string
}

export class ApiHandler {
    // TODO: proper error handling...
    private static authToken: string | undefined = localStorage.getItem('authToken') || undefined;
    public static isAuthenticated = ref(!!this.authToken);

    // handle auth here, we have only 2 endpoints anyway...
    // eslint-disable-next-line
    public static async authenticate(user: string, pass: string): Promise<true> {
        localStorage.removeItem('authToken');
        this.authToken = undefined;
        this.isAuthenticated.value = false;
        const response = await post(
            buildUrl(config.API_URL, config.AUTH_PATH),
            { user, pass });
        if (response.ok) {
            const responseBody = await response.json() as AuthSuccessResponse;
            // deliberate autocast
            // eslint-disable-next-line eqeqeq
            if (responseBody.status == 200 && !!responseBody.jwt) {
                this.authToken = responseBody.jwt;
                localStorage.setItem('authToken', this.authToken);
                this.isAuthenticated.value = true;
                return true;
            }
            throw new Error(`${response?.status}`); // shouldn't happen...
        }
        const responseErr = await response.json() as HttpError;
        throw new Error(`${responseErr.status}: ${responseErr.error} - ${responseErr.message}`);
    }

    public static async shortUrl(url: string, customShortText?: string): Promise<ShortenResponse> {
        if (!this.isAuthenticated) {
            throw new Error('Not Authenticated');
        }
        const response = await post(
            buildUrl(config.API_URL, config.SHORTEN_PATH, url),
            customShortText ? { customShortText } : undefined,
            { Authorization: `Bearer ${this.authToken}` });
        if (response.ok) {
            return {
                shortUrl: await response.text() // TODO: backend should really return json
            };
        }
        const responseErr = await response.json() as HttpError;
        // deliberate autocast
        // eslint-disable-next-line eqeqeq
        if (responseErr.status == 401) {
            this.authToken = undefined;
            localStorage.removeItem('authToken');
            this.isAuthenticated.value = false;
        }
        throw new Error(`${responseErr.status}: ${responseErr.error} - ${responseErr.message}`);
    }
}
