<template>
    <v-app class="gradient">
        <v-main>
            <router-view />
        </v-main>
    </v-app>
</template>

<script lang="ts">
import { Vue } from 'vue-class-component';

export default class App extends Vue {

}
</script>

<style lang="scss">
#app {
    height: 100%;
    width: 100%;
    margin: 0;
    background-color: #202020;
}
body,
html {
    width: 100%;
    height: 100%;
}

.gradient {
    display: block;
    position: fixed;
    top: 0;
    background: radial-gradient(#292929cc, #202020ed), url(assets/noise.svg);
    width: 100%;
    height: 100%;
    z-index: 0;
}

.v-main {
    height: 100;
}
.v-main__wrap {
    justify-self: center;
    height: 100;
}

.v-application__wrap {
    height: 100%;
}

.fullHeight {
    // vuetify 3.X beta fill-height seems broken
    height: 100%;
}

button {
    transition-duration: 600ms !important;
    animation-timing-function: cubic-bezier(0.65, 0.05, 0.36, 1) !important;
    &:hover {
        background-color: rgba(255, 255, 255, 0.2);
        box-shadow: 0 3px 8px 4px #161616af !important;
    }
}
</style>
