/* eslint-disable no-underscore-dangle */
/* eslint-disable no-use-before-define */
const { defineConfig } = require('@vue/cli-service');
const fs = require('fs');
const { SourceMapConsumer, SourceMapGenerator } = require('source-map');

module.exports = defineConfig({
    // handled by vscode. Always, not just on save.
    lintOnSave: false,

    pluginOptions: {
        vuetify: {
            // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vuetify-loader
        }
    },
    configureWebpack: {
        devtool: 'source-map',

        // stolen from: https://github.com/fearnycompknowhow/vue-typescript-debugger/blob/master/vue.config.js
        // holy crap that code needs to be fixed!
        // but it works for now... SO what gives...
        plugins: [{
            apply(compiler) {
                compiler.hooks.thisCompilation.tap('Initializing Compilation', (compilation) => {
                    compilation.hooks.finishModules.tapPromise('All Modules Built', async (modules) => {
                        // eslint-disable-next-line no-restricted-syntax
                        for (const module of modules) {
                            // eslint-disable-next-line no-continue
                            if (shouldSkipModule(module)) continue;

                            const pathWithoutQuery = module.resource.replace(/\?.*$/, '');
                            const sourceFile = fs.readFileSync(pathWithoutQuery).toString('utf-8');
                            const sourceMap = extractSourceMap(module);

                            sourceMap.sources = [pathWithoutQuery];
                            sourceMap.sourcesContent = [sourceFile];
                            // eslint-disable-next-line no-await-in-loop
                            sourceMap.mappings = await shiftMappings(sourceMap, sourceFile, pathWithoutQuery);
                        }
                    });
                });
            }
        }]
    }
});

function shouldSkipModule(module) {
    const { resource = '' } = module;

    if (!resource) return true;
    if (/node_modules/.test(resource)) return true;
    if (!/\.vue/.test(resource)) return true;
    if (!/type=script/.test(resource)) return true;
    if (!/lang=ts/.test(resource)) return true;
    if (isMissingSourceMap(module)) return true;

    return false;
}

function isMissingSourceMap(module) {
    return !extractSourceMap(module);
}

function extractSourceMap(module) {
    if (!module._source) return null;

    return module._source._sourceMap
        || module._source._sourceMapAsObject
        || null;
}

async function shiftMappings(sourceMap, sourceFile, sourcePath) {
    const indexOfScriptTag = getIndexOfScriptTag(sourceFile);

    const shiftedSourceMap = await SourceMapConsumer.with(sourceMap, null, async (consumer) => {
        const generator = new SourceMapGenerator();
        let original;
        consumer.eachMapping((mapping) => {
            const {
                generatedColumn,
                generatedLine,
                originalColumn,
                originalLine
            } = mapping;

            let { name } = mapping;
            let source = sourcePath;

            if (originalLine === null || originalColumn === null) {
                name = null;
                source = null;
            } else {
                original = {
                    column: originalColumn,
                    line: originalLine + indexOfScriptTag
                };
            }

            generator.addMapping({
                generated: {
                    column: generatedColumn,
                    line: generatedLine
                },
                original,
                source,
                name
            });
        });

        return generator.toJSON();
    });

    return shiftedSourceMap.mappings;
}

function getIndexOfScriptTag(sourceFile) {
    const lines = sourceFile.match(/.+/g);
    let indexOfScriptTag = 0;

    // eslint-disable-next-line no-restricted-syntax
    for (const line of lines) {
        // eslint-disable-next-line no-plusplus
        ++indexOfScriptTag;
        if (/<script/.test(line)) break;
    }

    return indexOfScriptTag;
}
