/**
 * Compiler configuration
 *
 * @see {@link https://roots.io/sage/docs sage documentation}
 * @see {@link https://bud.js.org/learn/config bud.js configuration guide}
 *
 * @type {import('@roots/bud').Config}
 */
export default async (app) => {
    /**
     * Application assets & entrypoint
     *
     * @see {@link https://bud.js.org/reference/bud.entry}
     * @see {@link https://bud.js.org/reference/bud.assets}
     */
    app
        .entry('app', ['@scripts/app', '@styles/app'])
        .assets(['images']);

    /**
     * Set public path
     *
     * @see {@link https://bud.js.org/reference/bud.setPublicPath}
     */
    app.setPublicPath('/static/dist/');
}
