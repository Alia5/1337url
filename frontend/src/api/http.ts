export async function post(url: string, body?: Record<string, unknown>, headers?: Record<string, string>): Promise<Response> {
    return fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            ...(headers || {})
        },
        body: body ? JSON.stringify(body) : undefined
    });
}
// TODO: other request wrappers
