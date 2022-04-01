export function trimDoubleSlash(i: string): string {
    return i.replace(/\/\//g, '/');
}

export function trimHttpS(i: string): string {
    return i.replace(/^http(s)*?:\/\//, '');
}

export function buildUrl(...i: (string | unknown)[]): string {
    return `${
        `${i?.[0]}`.startsWith('http:') ? 'http://' : 'https://'
    }${
        trimDoubleSlash(i.concat(trimHttpS(`${i.reverse().pop()}`)).reverse().join('/'))
    }`;
}
