<template>
    <v-container
        fill-height
        fluid
        class="fullHeight"
        @keypress.enter="login"
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
                        md=4
                        xs=auto
                    >
                        <!-- eslint-disable vuejs-accessibility/no-autofocus -->
                        <v-text-field
                            variant="outlined"
                            clearable
                            label="Username"
                            v-model="user"
                            autofocus
                            :error="!!loginError"
                            :error-messages="loginError"
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
                            label="Password"
                            v-model="pass"
                            type="password"
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
                            height="56"
                            width="132"
                            class="ml-4 mb-9"
                            @click="login"
                        >
                            Login
                        </v-btn>
                    </v-col>
                </v-row>
            </v-col>
        </v-row>
        <v-overlay
            :model-value="loading"
            class="align-center justify-center"
        >
            <v-progress-circular
                indeterminate
                size="64"
            ></v-progress-circular>
        </v-overlay>
    </v-container>
</template>
<script lang='ts'>
import { Vue } from 'vue-class-component';
import { ApiHandler } from '@/api';

export default class AuthComp extends Vue {
    public user = '';
    public pass = '';
    public loading = false;

    public loginError = '';
    public login() {
        this.loading = true;
        ApiHandler.authenticate(this.user, this.pass).catch((err) => {
            this.loginError = `${err}`;
        }).finally(() => {
            this.loading = false;
        });
    }
}
</script>
