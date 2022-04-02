<template>
    <v-container
        fill-height
        fluid
        class="fullHeight"
        @keypress.enter="shorten"
    >
        <v-row
            align="center"
            justify="center"
            fill-height
            class="fullHeight"
        >
            <v-col
                cols="12"
                md="6"
            >
                <v-row
                    dense
                    align="center"
                    justify="center"
                >
                    <v-col
                        cols=12
                        md=6
                        xs=auto
                    >
                        <!-- eslint-disable vuejs-accessibility/no-autofocus -->
                        <v-text-field
                            variant="outlined"
                            clearable
                            label="URL"
                            v-model="fullUrl"
                            autofocus
                            :error="!!validators.fullUrl(fullUrl)"
                            :error-messages="validators.fullUrl(fullUrl)"
                        >
                        </v-text-field>
                        <!-- eslint-enable vuejs-accessibility/no-autofocus -->
                    </v-col>
                    <v-col
                        cols=12
                        md=4
                        xs=auto
                    >
                        <v-text-field
                            variant="outlined"
                            clearable
                            label="Custom Text"
                            v-model="customText"
                            :error="!!validators.customText(customText)"
                            :error-messages="validators.customText(customText)"
                        >
                        </v-text-field>
                    </v-col>
                    <v-col
                        cols=12
                        md="auto"
                    >
                        <v-btn
                            elevation="2"
                            variant="outlined"
                            class="ml-4 mb-9"
                            height="56"
                            width="132"
                            @click="shorten"
                        >
                            Short
                        </v-btn>
                    </v-col>
                </v-row>
                <v-row
                    dense
                    justify="center"
                    class="mt-4"
                >
                    <v-col
                        cols=12
                        md=11
                        xs=auto
                        align="end"
                    >
                        <v-text-field
                            v-show="!!shortUrl"
                            readonly
                            variant="underlined"
                            label=""
                            hide-details
                            v-model="shortUrl"
                            ref="shortUrlInput"
                        >
                        </v-text-field>
                        <p class="mt-n7 mr-1">{{copyState}}</p>
                    </v-col>
                </v-row>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang='ts'>
import { Vue } from 'vue-class-component';
import { ref } from 'vue';
import { ApiHandler } from '@/api';
import { trimHttpS } from '@/util';

export default class URLShortner extends Vue {
    protected fullUrl = '';
    protected customText = '';
    protected shortUrl = '';
    protected copyState = '';
    public shortUrlInput = ref<HTMLInputElement | undefined>(undefined);

    // rules also seem broken in vuetify 3.X beta; run my own.
    protected validators = {
        fullUrl: (value: string) => {
            // eslint-disable-next-line eqeqeq
            if (value?.length == 0) {
                return 'URL is required';
            }
            const minUrlLen = 6;
            if (value.length < minUrlLen) {
                return `At least ${minUrlLen} characters`;
            }
            return '';
        },
        customText: (value: string) => (value.match(/^[a-zA-Z0-9]*$/gm)?.length ? '' : 'Custom Text must be alphanumeric')
    };

    protected shorten() {
        if (!this.fullUrl) {
            return;
        }
        ApiHandler.shortUrl(trimHttpS(this.fullUrl), this.customText)
            .then((res) => {
                this.shortUrl = res.shortUrl;
                this.customText = '';
                this.fullUrl = '';
                this.$nextTick(() => {
                    // and this is why i fucking despise frameworks
                    // not exclusively designed for **and** written in typescript
                    // why don't I use angular, again?!
                    // no! allowing any is NEVER an option
                    (this.shortUrlInput as unknown as HTMLInputElement | undefined)?.select();
                    this.copyTextToClipboard(this.shortUrl);
                });
            })
            .catch((err) => {
                this.shortUrl = '';
                this.copyState = `${err}`;
                this.resetCopyStateDelayed();
            });
    }

    private resetCopyStateDelayed() {
        setTimeout(() => {
            this.copyState = '';
        }, 1000);
    }
    // stolen from
    private fallbackCopyTextToClipboard(text: string) {
        const textArea = document.createElement('textarea');
        textArea.value = text;

        // Avoid scrolling to bottom
        textArea.style.top = '0';
        textArea.style.left = '0';
        textArea.style.position = 'fixed';

        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        try {
            document.execCommand('copy');
            this.copyState = 'Copied!';
            this.resetCopyStateDelayed();
        } catch (err) {
            this.copyState = 'Could not copy URL!';
            this.resetCopyStateDelayed();
        }

        document.body.removeChild(textArea);
    }
    private copyTextToClipboard(text: string) {
        if (!navigator.clipboard) {
            this.fallbackCopyTextToClipboard(text);
            return;
        }
        navigator.clipboard.writeText(text).then(() => {
            this.copyState = 'Copied!';
            this.resetCopyStateDelayed();
        }, () => {
            this.copyState = 'Could not copy URL!';
            this.resetCopyStateDelayed();
        });
    }
}
</script>
